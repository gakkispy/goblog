/*
 * @Date: 2023-01-30 12:35:41
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-30 12:49:29
 * @FilePath: /goblog/pkg/database/database.go
 */
package database

import (
	"database/sql"
	"goblog/pkg/logger"
	"time"

	"github.com/go-sql-driver/mysql"
)

// DB is a global variable for *sql.DB.
var DB *sql.DB

// Initialize 初始化数据库
func Initialize() {
	DB = initDB()
	createTables()
}

// initDB 初始化数据库
func initDB() *sql.DB {
	var err error

	// 1. 准备数据库配置
	config := mysql.Config{
		User:                 "gakkispy",
		Passwd:               "secret",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "goblog",
		AllowNativePasswords: true,
	}

	// 2. 准备数据库连接池，连接数据库
	DB, err := sql.Open("mysql", config.FormatDSN())
	logger.LogError(err)

	// 2.1 设置数据库连接池最大连接数
	DB.SetMaxOpenConns(25)
	// 2.2 设置数据库连接池最大空闲连接数
	DB.SetMaxIdleConns(25)
	// 2.3 设置每个连接的过期时间
	DB.SetConnMaxLifetime(5 * time.Minute)

	// 3. 尝试连接数据库，设置错误日志
	err = DB.Ping()
	logger.LogError(err)

	return DB
}

// createTables 创建数据表
func createTables() {
	createArticlesSQL := `
	CREATE TABLE IF NOT EXISTS articles(
		id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
		title varchar(255) NOT NULL DEFAULT '',
		body longtext NOT NULL,
		created_at datetime NOT NULL,
		updated_at datetime NOT NULL,
		deleted_at datetime NOT NULL,
		PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	_, err := DB.Exec(createArticlesSQL)
	logger.LogError(err)
}
