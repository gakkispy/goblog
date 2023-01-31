/*
 * @Date: 2023-01-30 16:54:54
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-31 11:28:51
 * @FilePath: /goblog/pkg/model/model.go
 */
package model

import (
	"goblog/pkg/logger"

	"gorm.io/gorm"

	// mysql driver
	"gorm.io/driver/mysql"

	//gorm logger
	gormlogger "gorm.io/gorm/logger"
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
	DB, err = gorm.Open(config, &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	logger.LogError(err)

	return DB
}
