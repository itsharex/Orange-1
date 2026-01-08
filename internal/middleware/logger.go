package middleware

import (
	"bytes"
	"log/slog"
	"time"

	"github.com/FruitsAI/Orange/internal/config"
	"github.com/gin-gonic/gin"
)

// bodyLogWriter 封装 gin.ResponseWriter 以捕获响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写 Write 方法以同时写入缓冲区
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果日志未启用，直接跳过
		if !config.AppConfig.LogEnable {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method

		// 包装 ResponseWriter
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// 获取查询参数
		if raw != "" {
			path = path + "?" + raw
		}

		// 构建日志字段
		logFields := []interface{}{
			"status", statusCode,
			"method", method,
			"path", path,
			"latency", latency,
			"ip", c.ClientIP(),
		}

		// 调试模式下打印响应体
		if config.AppConfig.LogLevel == "debug" {
			logFields = append(logFields, "response", blw.body.String())
		}

		// 打印日志
		slog.Info("Request", logFields...)
	}
}
