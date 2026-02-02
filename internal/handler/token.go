package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/FruitsAI/Orange/internal/middleware"
	"github.com/FruitsAI/Orange/internal/models"
	"github.com/FruitsAI/Orange/internal/repository"
	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	repo *repository.TokenRepository
}

func NewTokenHandler() *TokenHandler {
	return &TokenHandler{repo: repository.NewTokenRepository()}
}

// CreateRequest 创建令牌请求参数
type CreateTokenRequest struct {
	Name      string `json:"name" binding:"required"`
	ExpiresIn int    `json:"expires_in"` // 过期时间 (天)，0 表示永不过期
}

// CreateResponse 创建令牌响应 (包含原始 Token)
type CreateTokenResponse struct {
	Token string                      `json:"token"`
	Data  *models.PersonalAccessToken `json:"data"`
}

// generateToken 生成随机 Token 字符串 (32字节 hex)
func generateToken() (string, error) {
	bytes := make([]byte, 16) // 16 bytes = 32 hex chars
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// 添加前缀以便识别
	return "pat_" + hex.EncodeToString(bytes), nil
}

// hashToken 计算 Token Hash
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Create 创建新的访问令牌
func (h *TokenHandler) Create(c *gin.Context) {
	var req CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid parameters"})
		return
	}

	userID := middleware.GetUserID(c)

	// 1. 生成原始 Token
	rawToken, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to generate token"})
		return
	}

	// 2. 计算 Hash
	tokenHash := hashToken(rawToken)

	// 3. 计算过期时间
	var expiresAt *time.Time
	if req.ExpiresIn > 0 {
		t := time.Now().AddDate(0, 0, req.ExpiresIn)
		expiresAt = &t
	}

	// 4. 构建模型
	token := &models.PersonalAccessToken{
		UserID:    userID,
		Name:      req.Name,
		TokenHash: tokenHash,
		Status:    1,
		ExpiresAt: expiresAt,
	}

	// 5. 保存到数据库
	if err := h.repo.Create(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to save token"})
		return
	}

	// 6. 返回结果 (包含原始 Token，仅此一次)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": CreateTokenResponse{
			Token: rawToken,
			Data:  token,
		},
	})
}

// List 获取令牌列表
func (h *TokenHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	tokens, err := h.repo.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to list tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tokens,
	})
}

// Revoke 撤销令牌
func (h *TokenHandler) Revoke(c *gin.Context) {
	idStr := c.Param("id")
	var id int64
	fmt.Sscanf(idStr, "%d", &id)

	userID := middleware.GetUserID(c)

	if err := h.repo.Revoke(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to revoke token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// Delete 删除令牌
func (h *TokenHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	var id int64
	fmt.Sscanf(idStr, "%d", &id)

	userID := middleware.GetUserID(c)

	if err := h.repo.Delete(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to delete token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
