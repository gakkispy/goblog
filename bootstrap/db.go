/*
 * @Date: 2023-01-30 17:00:35
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-30 17:02:46
 * @FilePath: /goblog/bootstrap/db.go
 */
package bootstrap

import (
	"goblog/pkg/model"
	"time"
)

// SetupDB initializes the database.
func SetupDB() {
	// 1. 建立数据库连接
	db := model.ConnectDB()

	// 2. 设置命令行打印 SQL
	sqlDB, _ := db.DB()

	// 3. 设置最大连接数
	sqlDB.SetMaxOpenConns(100)
	// 4. 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(25)
	// 5. 设置每个链接的过期时间
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
}
