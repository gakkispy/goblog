/*
 * @Date: 2023-01-30 11:27:10
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-30 11:28:01
 * @FilePath: /goblog/pkg/types/converter.go
 */
package types

import "strconv"

// Int64ToString converts int64 to string.
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
