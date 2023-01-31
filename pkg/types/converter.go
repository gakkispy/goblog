/*
 * @Date: 2023-01-30 11:27:10
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-31 09:42:21
 * @FilePath: /goblog/pkg/types/converter.go
 */
package types

import (
	"goblog/pkg/logger"
	"strconv"
)

// Int64ToString converts int64 to string.
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// StringToUint64 converts string to uint64.
func StringToUint64(s string) uint64 {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		logger.LogError(err)
	}
	return i
}

// Uint64ToString converts uint64 to string.
func Uint64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}
