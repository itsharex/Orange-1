package password

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 算法对明文密码进行加密
// 返回加密后的哈希字符串。如果加密失败，返回 error。
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost 默认工作因子，平衡了安全性与性能
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 校验明文密码与哈希值是否匹配
// 返回 true 表示密码正确，false 表示不匹配。
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
