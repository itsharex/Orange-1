package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一 API 响应结构体
// 所有 API 接口返回的数据都应符合此格式。
type Response struct {
	Code    int         `json:"code"`    // 业务状态码 (0: 成功, 非0: 错误)
	Message string      `json:"message"` // 提示消息
	Data    interface{} `json:"data"`    // 业务数据
}

// PageData 分页数据封装
// 用于列表查询接口的标准返回格式。
type PageData struct {
	List     interface{} `json:"list"`      // 当前页的数据列表
	Total    int64       `json:"total"`     // 总记录数
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页条数
}

// 业务错误码定义
// 0 表示成功，其他值表示各种错误类型。
const (
	CodeSuccess       = 0    // 请求成功
	CodeParamError    = 1001 // 参数错误 (如必填项缺失、格式不对)
	CodeNotFound      = 1002 // 资源未找到
	CodeUnauthorized  = 2001 // 未授权 (未登录或 Token 无效)
	CodeTokenExpired  = 2002 // Token 已过期
	CodeForbidden     = 2003 // 禁止访问 (无权限)
	CodeInternalError = 5000 // 服务器内部错误
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

// Success 发送成功响应 (Code=0)
// data: 返回给前端的业务数据对象
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

// SuccessPage 发送分页数据的成功响应
// list: 当前页的数据切片
// total: 总记录数
// page: 当前页码
// pageSize: 每页大小
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

// Error 发送错误响应
// code: 业务错误码 (非0)
// message: 可选的自定义错误消息。如果不传，则使用 code 对应的默认消息。
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
