package repository

import (
	"github.com/FruitsAI/Orange/internal/database"
	"github.com/FruitsAI/Orange/internal/models"
	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{db: database.GetDB()}
}

// Create 创建新的令牌
func (r *TokenRepository) Create(token *models.PersonalAccessToken) error {
	return r.db.Create(token).Error
}

// List 获取用户的令牌列表 (隐藏 Hash)
func (r *TokenRepository) List(userID int64) ([]models.PersonalAccessToken, error) {
	var tokens []models.PersonalAccessToken
	// 只查询未删除/未过期的？目前需求是管理所有 token，让用户自己决定撤销
	// 按照创建时间倒序
	if err := r.db.Where("user_id = ?", userID).Order("id DESC").Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

// Revoke 撤销令牌 (软删除或状态变更)
func (r *TokenRepository) Revoke(id int64, userID int64) error {
	return r.db.Model(&models.PersonalAccessToken{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("status", 0).Error
}

// Delete 删除令牌 (硬删除)
func (r *TokenRepository) Delete(id int64, userID int64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.PersonalAccessToken{}).Error
}

// FindByHash 根据 Token Hash 查找有效令牌
func (r *TokenRepository) FindByHash(hash string) (*models.PersonalAccessToken, error) {
	var token models.PersonalAccessToken
	if err := r.db.Preload("User").Where("token_hash = ? AND status = 1", hash).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}
