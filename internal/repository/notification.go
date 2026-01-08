package repository

import (
	"github.com/FruitsAI/Orange/internal/database"
	"github.com/FruitsAI/Orange/internal/models"
	"gorm.io/gorm"
)

// NotificationRepository 通知数据仓库
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository 创建通知仓库
func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{db: database.GetDB()}
}

// Create 创建通知
func (r *NotificationRepository) Create(notification *models.Notification, targetUserID int64) error {
	// 如果是全员通知 (IsGlobal=1)，直接插入 notifications 表
	if notification.IsGlobal == 1 {
		return r.db.Create(notification).Error
	}

	// 如果是私信 (IsGlobal=0)，开启事务
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 插入 notifications 表
		if err := tx.Create(notification).Error; err != nil {
			return err
		}

		// 2. 插入 user_notifications 表 (Inbox模式: 私信同时写入接收者记录，状态为未读)
		userNotification := models.UserNotification{
			UserID:         targetUserID,
			NotificationID: notification.ID,
			IsRead:         0,
		}
		if err := tx.Create(&userNotification).Error; err != nil {
			return err
		}

		return nil
	})
}

// Update 更新通知
func (r *NotificationRepository) Update(notification *models.Notification) error {
	return r.db.Model(notification).Omit("Sender").Save(notification).Error
}

// FindByID 根据ID查找通知
func (r *NotificationRepository) FindByID(id int64) (*models.Notification, error) {
	var notification models.Notification
	// 手动关联查询 Sender 信息
	err := r.db.Table("notifications").
		Select("notifications.*, users.id as sender_id, users.username as sender_username, users.name as sender_name").
		Joins("LEFT JOIN users ON users.id = notifications.sender_id").
		Where("notifications.id = ?", id).
		Scan(&notification).Error

	if err != nil {
		return nil, err
	}

	// 简单填充 Sender 指针，保持兼容
	sender := models.User{
		ID: notification.SenderID,
	}
	// 如果有关联查询出来的字段，可以赋值给 sender (需要 struct scan 支持，这里简单再查一次或优化)
	// 为保稳妥，使用 First 查询完整 Sender
	if err := r.db.First(&sender, notification.SenderID).Error; err == nil {
		notification.Sender = &sender
	}

	return &notification, nil
}

// ListByUser 获取用户通知（包含全员通知和私信）
func (r *NotificationRepository) ListByUser(userID int64, offset, limit int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	// 查询逻辑：
	// 1. 全员通知 (IsGlobal=1): 无论 user_notifications 是否有记录都显示。如果 user_notifications 有记录则为已读，否则未读。
	// 2. 私信 (IsGlobal=0): 必须在 user_notifications 有记录才显示。is_read 状态取自 user_notifications。

	db := r.db.Table("notifications").
		Joins("LEFT JOIN user_notifications un ON notifications.id = un.notification_id AND un.user_id = ?", userID).
		Where("(notifications.is_global = 1) OR (un.id IS NOT NULL)")

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := db.Select("notifications.*, COALESCE(un.is_read, 0) as is_read").
		Order("notifications.create_time DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifications).Error

	if err != nil {
		return nil, 0, err
	}

	r.fillSenders(notifications)
	return notifications, total, nil
}

// ListAll 获取所有通知（管理员用）
func (r *NotificationRepository) ListAll(offset, limit int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	if err := r.db.Model(&models.Notification{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Order("create_time DESC").Offset(offset).Limit(limit).Find(&notifications).Error; err != nil {
		return nil, 0, err
	}
	r.fillSenders(notifications)
	return notifications, total, nil
}

// MarkAsRead 标记为已读
func (r *NotificationRepository) MarkAsRead(id int64, userID int64) error {
	// 检查 UserNotification 是否已存在
	var un models.UserNotification
	err := r.db.Where("user_id = ? AND notification_id = ?", userID, id).First(&un).Error

	if err == nil {
		// 存在记录（私信 OR 已读过的全员通知），更新状态为已读
		if un.IsRead == 0 {
			now := database.GetDB().NowFunc() // 使用 gorm 获取当前时间，或者 time.Now()
			un.IsRead = 1
			un.ReadTime = &now
			return r.db.Save(&un).Error
		}
		return nil
	} else if err == gorm.ErrRecordNotFound {
		// 不存在记录（说明是未读的全员通知），插入已读记录
		// 先确认该通知是否真的存在且是 IsGlobal=1
		var n models.Notification
		if err := r.db.First(&n, id).Error; err != nil {
			return err
		}

		if n.IsGlobal == 1 {
			now := database.GetDB().NowFunc()
			newUN := models.UserNotification{
				UserID:         userID,
				NotificationID: id,
				IsRead:         1,
				ReadTime:       &now,
			}
			return r.db.Create(&newUN).Error
		}
		return nil // 如果不是全员通知且没记录，理论上用户不应看到，不做操作
	} else {
		return err
	}
}

// Delete 删除通知
func (r *NotificationRepository) Delete(id int64) error {
	// 事务删除：删除通知和关联的 user_notifications
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("notification_id = ?", id).Delete(&models.UserNotification{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Notification{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetUnreadCount 获取未读通知数量
func (r *NotificationRepository) GetUnreadCount(userID int64) (int64, error) {
	var count int64
	// 统计逻辑：
	// 1. 全员通知 (IsGlobal=1) 且 user_notifications 中不存在记录 (un.id IS NULL)
	// 2. 私信 (IsGlobal=0) 且 user_notifications 中存在记录但 is_read=0

	err := r.db.Table("notifications").
		Joins("LEFT JOIN user_notifications un ON notifications.id = un.notification_id AND un.user_id = ?", userID).
		Where("((notifications.is_global = 1 AND (un.id IS NULL OR un.is_read = 0)) OR (notifications.is_global = 0 AND un.is_read = 0 AND un.id IS NOT NULL))").
		Count(&count).Error

	return count, err
}

// fillSenders 辅助方法：批量填充 Sender 信息
func (r *NotificationRepository) fillSenders(notifications []models.Notification) {
	senderIDs := make([]int64, 0)
	for _, n := range notifications {
		senderIDs = append(senderIDs, n.SenderID)
	}

	if len(senderIDs) == 0 {
		return
	}

	var senders []models.User
	r.db.Where("id IN ?", senderIDs).Find(&senders)

	senderMap := make(map[int64]*models.User)
	for i := range senders {
		senderMap[senders[i].ID] = &senders[i]
	}

	for i := range notifications {
		if sender, ok := senderMap[notifications[i].SenderID]; ok {
			notifications[i].Sender = sender
		}
	}
}

// ListUsers 获取用户列表（用于选择发送目标）
func (r *NotificationRepository) ListUsers() ([]models.User, error) {
	var users []models.User
	if err := r.db.Select("id", "name", "username").
		Where("status = 1").
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
