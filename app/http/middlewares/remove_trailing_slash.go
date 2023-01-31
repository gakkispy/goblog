/*
 * @Date: 2023-01-31 16:10:26
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-31 16:11:03
 * @FilePath: /goblog/app/http/middlewares/remove_trailing_slash.go
 */
package middlewares

import (
	"net/http"
	"strings"
)

// RemoveTrailingSlash, 去除 URL 结尾的斜杠
func RemoveTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}
