package database

import (
	"log/slog"
	"sync"

	"github.com/FruitsAI/Orange/internal/config"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// db 是全局唯一的数据库连接实例 (单例模式)
	db *gorm.DB
	// once 用于确保数据库初始化只执行一次
	once sync.Once
)

// GetDB 获取数据库连接实例 (单例)
// 该方法是并发安全的，首次调用时会自动初始化数据库连接。
// 这里的初始化包括打开 SQLite 文件连接和配置 GORM 引擎。
func GetDB() *gorm.DB {
	once.Do(func() {
		var err error
		db, err = initDB()
		if err != nil {
			slog.Error("Failed to initialize database", "error", err)
			// 数据库是核心依赖，初始化失败直接 Panic 终止应用
			panic(err)
		}
	})
	return db
}

// initDB 初始化 SQLite 数据库连接
// 根据配置文件中的路径打开数据库，并配置 GORM 的日志级别。
// 使用纯 Go 实现的 glebarez/sqlite 驱动，无需 CGO 支持。
func initDB() (*gorm.DB, error) {
	// 1. 获取数据库文件路径 (由 config 模块加载)
	dbPath := config.AppConfig.DBPath
	slog.Info("Opening database", "path", dbPath)

	// 2. 建立 GORM 连接
	// 默认开启 Info 级别日志，便于调试 SQL 语句
	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return database, nil
}

// Close 关闭数据库连接
// 主要是为了释放底层 sql.DB 的连接资源 (通常在应用退出时调用)
func Close() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
