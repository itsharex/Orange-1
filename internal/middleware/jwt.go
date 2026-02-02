package middleware

import (
	"strings"

	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/FruitsAI/Orange/internal/database"
	"github.com/FruitsAI/Orange/internal/models"
	"github.com/FruitsAI/Orange/internal/pkg/jwt"
	"github.com/FruitsAI/Orange/internal/pkg/response"
	"github.com/FruitsAI/Orange/internal/repository"
	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 鉴权中间件
// 拦截 HTTP 请求，验证 Request Header 中的 Authorization 字段。
// 支持:
// 1. 标准 JWT (Bearer <token>)
// 2. 个人访问令牌 (Bearer pat_<token>)
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "请先登录")
			return
		}

		// 2. 解析 Bearer Token 格式 (Bearer <token>)
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Token格式错误")
			return
		}

		tokenString := parts[1]

		// 3.1 检查是否为 Personal Access Token (以 pat_ 开头)
		if strings.HasPrefix(tokenString, "pat_") {
			// 计算 Hash
			hash := sha256.Sum256([]byte(tokenString))
			tokenHash := hex.EncodeToString(hash[:])

			// 查找 Token
			repo := repository.NewTokenRepository()
			token, err := repo.FindByHash(tokenHash)
			if err != nil {
				response.Unauthorized(c, "无效的访问令牌")
				return
			}

			// 检查过期时间
			if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now()) {
				response.Unauthorized(c, "访问令牌已过期")
				return
			}

			// 更新最后使用时间 (异步，不阻塞请求)
			go func(id int64) {
				// 注意: 这里简单更新，实际可能需要考虑并发或直接使用 SQL Update
				// 这里为了方便直接复用 DB
				database.GetDB().Model(&models.PersonalAccessToken{}).Where("id = ?", id).Update("last_used_at", time.Now())
			}(token.ID)

			// 验证关联用户
			if token.User == nil {
				response.Unauthorized(c, "关联用户无效")
				return
			}

			// 4. 将用户信息注入上下文 (Context)
			c.Set("user_id", token.User.ID)
			c.Set("username", token.User.Username)
			c.Set("role", token.User.Role)
			c.Set("access_token_id", token.ID) // 标记来源

			c.Next()
			return
		}

		// 3.2 校验并解析标准 JWT
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			response.Error(c, response.CodeTokenExpired, "Token已过期或无效")
			c.Abort()
			return
		}

		// 4. 将用户信息注入上下文 (Context)
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) int64 {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(int64)
	}
	return 0
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}

// GetRole 从上下文获取角色
func GetRole(c *gin.Context) string {
	if role, exists := c.Get("role"); exists {
		return role.(string)
	}
	return ""
}
