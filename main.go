/*
 * @Date: 2023-01-16 14:28:24
 * @LastEditors: gakkispy && yaosenjun@cii.com
 * @LastEditTime: 2023-01-16 15:06:29
 * @FilePath: /goblog/main.go
 */
package main

import (
	"fmt"
	"net/http"
)

func handleFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, 这里是 goblog home page !</h1>")
	} else if r.URL.Path == "/about" {
		fmt.Fprint(w, "此博客由 Go 语言编写，使用 Gin 框架，使用 GORM 操作数据库，使用 PostgreSQL 数据库。联系方式："+
			"<a href='mailto:yaosenjun168@live.cn'>yaosenjun168@live.cn</a>。")
	} else {
		fmt.Fprint(w, "<h1>Hello, 这里是 goblog 404 page !</h1>")
	}
}

func main() {
	http.HandleFunc("/", handleFunc)
	http.ListenAndServe(":3000", nil)
}
