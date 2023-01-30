/*
 * @Date: 2023-01-30 16:54:54
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-30 16:57:35
 * @FilePath: /goblog/pkg/model/model.go
 */
package model

import (
	"goblog/pkg/logger"

	"github.com/jinzhu/gorm"

	// mysql driver
	"gorm.io/driver/mysql"
)

// DB is a global variable for *gorm.DB.
var DB *gorm.DB

// ConnectDB 初始化模型
func ConnectDB() *gorm.DB {
	var err error

	config := mysql.New(mysql.Config{
		DSN: "gakkispy:secret@tcp(127.0.0.1:3306)/goblog?charset=utf8mb4&parseTime=True&loc=Local",
	})

	// 准备数据库连接池，连接数据库
	DB, err = gorm.Open(config, &gorm.Config{})
	logger.LogError(err)

	return DB
}
