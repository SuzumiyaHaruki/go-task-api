/*
database.go 负责初始化数据库连接。

本文件从环境变量读取 MySQL 连接配置，使用 GORM 建立数据库连接，
配置连接池，并在应用启动时自动迁移用户和任务表结构。
*/
package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
databaseConfig 表示连接 MySQL 所需的基础配置。

Host 和 Port 指向数据库服务地址，User、Password 和 Name 分别表示
数据库用户名、密码和数据库名称。
*/
type databaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

/*
loadDatabaseConfig 从环境变量读取数据库配置。

当环境变量不存在时，会使用适合本地开发的默认值；Docker Compose
运行时会通过 .env 文件覆盖这些值。
*/
func loadDatabaseConfig() databaseConfig {
	return databaseConfig{
		Host:     getenv("DB_HOST", "127.0.0.1"),
		Port:     getenv("DB_PORT", "3306"),
		User:     getenv("DB_USER", "task_user"),
		Password: getenv("DB_PASSWORD", "task_password"),
		Name:     getenv("DB_NAME", "task_api"),
	}
}

/*
dsn 生成 MySQL 连接字符串。

parseTime=True 让 MySQL 的时间字段正确映射为 Go 的 time.Time；
charset=utf8mb4 用于支持完整 Unicode 字符集。
*/
func (cfg databaseConfig) dsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
}

/*
openDatabase 创建 GORM 数据库连接并迁移表结构。

为了适配 Docker Compose 中 MySQL 首次启动需要初始化的情况，
连接失败时会短暂重试。连接成功后会设置连接池参数，并自动创建或更新
users 和 tasks 表。
*/
func openDatabase() (*gorm.DB, error) {
	cfg := loadDatabaseConfig()

	var db *gorm.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = gorm.Open(mysql.Open(cfg.dsn()), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get database handle: %w", err)
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.AutoMigrate(&user{}, &task{}); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	return db, nil
}
