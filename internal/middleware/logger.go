package middleware

import (
	"bytes"
	"log/slog"
	"time"

	"github.com/FruitsAI/Orange/internal/config"
	"github.com/FruitsAI/Orange/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// bodyLogWriter 自定义 ResponseWriter
// 用于在将响应写入 HTTP 连接的同时，将其复制到内存缓冲区中，
// 以便后续在日志中记录响应体内容（主要用于 Debug 模式）。
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写 Write 方法
// 实现了双写机制：既写入原始 ResponseWriter，也写入 buffer。
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger 全局请求日志中间件
// 记录每个 HTTP 请求的详细信息，包括：
// - 状态码、请求方法、请求路径
// - 处理耗时 (Latency)
// - 客户端 IP
// - (Debug模式) 响应体内容
//
// 该中间件依赖 config.AppConfig.LogEnable 配置项决定是否启用。
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

		// 包装 ResponseWriter 以捕获响应体
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
		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
		}

		// 调试模式下记录响应体内容
		if config.AppConfig.LogLevel == "debug" {
			fields = append(fields, zap.String("response", blw.body.String()))
		}

		// 输出结构化日志
		if logger.Log != nil {
			logger.Log.Info("Request", fields...)
		} else {
			// Fallback: 如果 logger 未初始化（理论上不应发生），降级使用 slog
			slog.Info("Request", "status", statusCode, "path", path)
		}
	}
}
