package database

import (
	"log/slog"

	"github.com/FruitsAI/Orange/internal/models"
	"gorm.io/gorm"
)

// Seed populates the database with initial data
func Seed(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	slog.Info("Seeding database...")

	return db.Transaction(func(tx *gorm.DB) error {
		// Users
		// Note: mirroring db/seed_users.sql
		usersSQL := []string{
			`INSERT INTO users (id, username, password, name, email, phone, role, department, position, status, create_time) VALUES 
(1, 'admin', '$2a$10$sp/NLWYRQjt9zVCq6HeOieFaFNBl79RoBXdePqxg9UhwQyT1/C7vu', '管理员', 'admin@orange.com', '13800000000', 'admin', '技术部', '系统管理员', 1, CURRENT_TIMESTAMP);`,
			`INSERT INTO users (id, username, password, name, email, phone, role, department, position, status, create_time) VALUES 
(2, 'xu', '$2a$10$sp/NLWYRQjt9zVCq6HeOieFaFNBl79RoBXdePqxg9UhwQyT1/C7vu', '郑旭', 'xu@company.com', '13800000001', 'user', '技术部', '项目经理', 1, CURRENT_TIMESTAMP);`,
		}

		for _, sql := range usersSQL {
			if err := tx.Exec(sql).Error; err != nil {
				return err
			}
		}

		// Dictionaries
		// Note: mirroring db/seed_dictionaries.sql
		dictSQL := []string{
			// payment_stage
			`INSERT INTO dictionaries (id, code, name, status, create_time) VALUES (1, 'payment_stage', '款项阶段', 1, CURRENT_TIMESTAMP);`,
			`INSERT INTO dictionary_item (dictionary_id, label, value, sort, status, create_time) VALUES 
(1, '首付款', 'deposit', 1, 1, CURRENT_TIMESTAMP),
(1, '进度款', 'progress', 2, 1, CURRENT_TIMESTAMP),
(1, '尾款', 'final', 3, 1, CURRENT_TIMESTAMP),
(1, '全款', 'all', 4, 1, CURRENT_TIMESTAMP);`,
			// payment_method
			`INSERT INTO dictionaries (id, code, name, status, create_time) VALUES (2, 'payment_method', '支付方式', 1, CURRENT_TIMESTAMP);`,
			`INSERT INTO dictionary_item (dictionary_id, label, value, sort, status, create_time) VALUES 
(2, '银行转账', 'bank_transfer', 1, 1, CURRENT_TIMESTAMP),
(2, '支付宝', 'alipay', 2, 1, CURRENT_TIMESTAMP),
(2, '微信支付', 'wechat', 3, 1, CURRENT_TIMESTAMP),
(2, '现金', 'cash', 4, 1, CURRENT_TIMESTAMP);`,
			// project_status
			`INSERT INTO dictionaries (id, code, name, status, create_time) VALUES (3, 'project_status', '项目状态', 1, CURRENT_TIMESTAMP);`,
			`INSERT INTO dictionary_item (dictionary_id, label, value, sort, status, create_time) VALUES 
(3, '未开始', 'notstarted', 1, 1, CURRENT_TIMESTAMP),
(3, '进行中', 'active', 2, 1, CURRENT_TIMESTAMP),
(3, '已完成', 'completed', 3, 1, CURRENT_TIMESTAMP),
(3, '已逾期', 'overdue', 4, 1, CURRENT_TIMESTAMP),
(3, '已归档', 'archived', 5, 1, CURRENT_TIMESTAMP);`,
			// project_type
			`INSERT INTO dictionaries (id, code, name, status, create_time) VALUES (4, 'project_type', '项目类型', 1, CURRENT_TIMESTAMP);`,
			`INSERT INTO dictionary_item (dictionary_id, label, value, sort, status, create_time) VALUES 
(4, 'Web开发', 'web', 1, 1, CURRENT_TIMESTAMP),
(4, '移动应用', 'mobile', 2, 1, CURRENT_TIMESTAMP),
(4, 'UI设计', 'design', 3, 1, CURRENT_TIMESTAMP),
(4, 'SaaS系统', 'saas', 4, 1, CURRENT_TIMESTAMP),
(4, '其他', 'other', 5, 1, CURRENT_TIMESTAMP);`,
		}

		for _, sql := range dictSQL {
			if err := tx.Exec(sql).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
