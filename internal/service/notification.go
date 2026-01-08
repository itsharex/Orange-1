package service

import (
	"errors"

	"github.com/FruitsAI/Orange/internal/models"
	"github.com/FruitsAI/Orange/internal/repository"
)

// NotificationService 通知服务
type NotificationService struct {
	notificationRepo *repository.NotificationRepository
}

// NewNotificationService 创建通知服务
func NewNotificationService() *NotificationService {
	return &NotificationService{
		notificationRepo: repository.NewNotificationRepository(),
	}
}

// Create 创建通知（仅管理员）
func (s *NotificationService) Create(senderID int64, title, content, notificationType string, targetUserID int64) (*models.Notification, error) {
	if title == "" {
		return nil, errors.New("标题不能为空")
	}
	if content == "" {
		return nil, errors.New("内容不能为空")
	}

	// 类型转换 string -> int
	// 1:system, 2:activity, 3:private
	var typeInt int
	switch notificationType {
	case "activity":
		typeInt = 2
	case "private":
		typeInt = 3
	default:
		typeInt = 1 // system
	}

	// 判断是否全员
	var isGlobal int
	if targetUserID == 0 {
		isGlobal = 1
	} else {
		isGlobal = 0
	}

	notification := &models.Notification{
		Title:    title,
		Content:  content,
		Type:     typeInt,
		SenderID: senderID,
		IsGlobal: isGlobal,
	}

	if err := s.notificationRepo.Create(notification, targetUserID); err != nil {
		return nil, errors.New("创建通知失败")
	}

	return notification, nil
}

// Update 更新通知（仅管理员）
func (s *NotificationService) Update(id int64, title, content, notificationType string, targetUserID int64) (*models.Notification, error) {
	notification, err := s.notificationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("通知不存在")
	}

	if title == "" {
		return nil, errors.New("标题不能为空")
	}
	if content == "" {
		return nil, errors.New("内容不能为空")
	}

	var typeInt int
	switch notificationType {
	case "activity":
		typeInt = 2
	case "private":
		typeInt = 3
	default:
		typeInt = 1 // system
	}

	var isGlobal int
	if targetUserID == 0 {
		isGlobal = 1
	} else {
		isGlobal = 0
	}

	notification.Title = title
	notification.Content = content
	notification.Type = typeInt
	notification.IsGlobal = isGlobal

	// 注意：Update 目前只更新通知表基本信息，未涉及 user_notifications 关系的调整
	// 如果需要支持从 私信改全员 或 全员改私信 的复杂逻辑，需要repo层配合修改

	if err := s.notificationRepo.Update(notification); err != nil {
		return nil, errors.New("更新通知失败")
	}

	return notification, nil
}

// ListByUser 获取用户通知列表
func (s *NotificationService) ListByUser(userID int64, page, pageSize int) ([]models.Notification, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return s.notificationRepo.ListByUser(userID, offset, pageSize)
}

// ListAll 获取所有通知列表（管理员）
func (s *NotificationService) ListAll(page, pageSize int) ([]models.Notification, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return s.notificationRepo.ListAll(offset, pageSize)
}

// MarkAsRead 标记为已读
func (s *NotificationService) MarkAsRead(id int64, userID int64) error {
	return s.notificationRepo.MarkAsRead(id, userID)
}

// Delete 删除通知（仅管理员）
func (s *NotificationService) Delete(id int64) error {
	return s.notificationRepo.Delete(id)
}

// GetUnreadCount 获取未读通知数量
func (s *NotificationService) GetUnreadCount(userID int64) (int64, error) {
	return s.notificationRepo.GetUnreadCount(userID)
}

// ListUsers 获取用户列表（用于选择发送目标）
func (s *NotificationService) ListUsers() ([]models.User, error) {
	return s.notificationRepo.ListUsers()
}

// Get 获取单个通知
func (s *NotificationService) Get(id int64) (*models.Notification, error) {
	return s.notificationRepo.FindByID(id)
}
