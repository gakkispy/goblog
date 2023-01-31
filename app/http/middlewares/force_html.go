/*
 * @Date: 2023-01-31 16:08:55
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-31 16:10:06
 * @FilePath: /goblog/app/http/middlewares/force_html.go
 */
package middlewares

import "net/http"

// ForceHTML, 强制浏览器使用 HTML 渲染
func ForceHTML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 强制 Content-Type 为 HTML
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}
