package logger

import (
	"log/slog"
	"sync"

	"github.com/FruitsAI/Orange/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Log 是全局的 Zap Logger 实例，供需要高性能日志的场景使用
	Log *zap.Logger
	// once 保证日志初始化只执行一次
	once sync.Once
)

// Setup 初始化全局日志系统
// 集成了 zap (高吞吐结构化日志)、lumberjack (日志轮转) 和 go 标准库 slog。
func Setup() {
	once.Do(func() {
		// 1. 获取日志文件路径配置
		logPath := config.AppConfig.LogPath
		if logPath == "" {
			logPath = "orange.log" // 默认文件名
		}

		// 2. 配置 Lumberjack 进行日志轮转 (Log Rotation)
		// 防止单个日志文件过大占满磁盘
		rotator := &lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    config.AppConfig.LogMaxSize,    // 单个文件最大尺寸 (MB)
			MaxBackups: config.AppConfig.LogMaxBackups, // 保留旧文件最大个数
			MaxAge:     config.AppConfig.LogMaxAge,     // 保留旧文件最大天数
			Compress:   config.AppConfig.LogCompress,   // 是否压缩旧日志
		}

		// 3. 配置 Zap 编码器 (JSON 格式)
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 时间格式: 2023-01-01T12:00:00.000Z
		encoderConfig.TimeKey = "time"

		// 4. 创建 Zap Core
		// 将日志输出指向 rotator (文件写入器)
		var core zapcore.Core
		fileSyncer := zapcore.AddSync(rotator)

		// 根据配置设置日志级别
		level := zap.InfoLevel
		if config.AppConfig.LogLevel == "debug" {
			level = zap.DebugLevel
		}

		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			fileSyncer,
			level,
		)

		// 5. 初始化 Zap Logger
		// AddCaller: 记录调用日志的文件名和行号
		Log = zap.New(core, zap.AddCaller())

		// 6. 重定向标准库 slog 到 zap 的 rotator
		// 这样应用中所有使用 slog.Info/Error 的地方也会写入同一个日志文件。
		slogHandler := slog.NewJSONHandler(rotator, &slog.HandlerOptions{
			Level: slogLevel(config.AppConfig.LogLevel),
		})
		slog.SetDefault(slog.New(slogHandler))
	})
}

// slogLevel 将配置字符串转换为 slog.Level 枚举
func slogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Sync 刷新缓冲区，确保所有日志写入磁盘
// 通常在 main 函数退出前调用: defer logger.Sync()
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
