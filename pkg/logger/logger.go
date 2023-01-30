/*
 * @Date: 2023-01-29 16:32:51
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-29 16:33:36
 * @FilePath: /goblog/pkg/logger/logger.go
 */
package logger

import "log"

// Logger is a global variable for *log.Logger.
func LogError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
