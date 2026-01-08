package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// PageData 分页数据结构
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// 错误码定义
const (
	CodeSuccess       = 0
	CodeParamError    = 1001
	CodeNotFound      = 1002
	CodeUnauthorized  = 2001
	CodeTokenExpired  = 2002
	CodeForbidden     = 2003
	CodeInternalError = 5000
)

// 错误消息映射
var codeMessages = map[int]string{
	CodeSuccess:       "success",
	CodeParamError:    "参数错误",
	CodeNotFound:      "资源不存在",
	CodeUnauthorized:  "未授权",
	CodeTokenExpired:  "Token已过期",
	CodeForbidden:     "禁止访问",
	CodeInternalError: "服务器内部错误",
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应带自定义消息
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// SuccessPage 分页成功响应
func SuccessPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data: PageData{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message ...string) {
	msg := codeMessages[code]
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: msg,
		Data:    nil,
	})
}

// ParamError 参数错误
func ParamError(c *gin.Context, message ...string) {
	Error(c, CodeParamError, message...)
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message ...string) {
	Error(c, CodeNotFound, message...)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message ...string) {
	Error(c, CodeUnauthorized, message...)
	c.Abort()
}

// InternalError 服务器错误
func InternalError(c *gin.Context, message ...string) {
	Error(c, CodeInternalError, message...)
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, message ...string) {
	Error(c, CodeForbidden, message...)
	c.Abort()
}
