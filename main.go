/*
 * @Date: 2023-01-16 14:28:24
 * @LastEditors: gakkispy && yaosenjun@cii.com
 * @LastEditTime: 2023-01-17 12:59:58
 * @FilePath: /goblog/main.go
 */
package main

import (
	"fmt"
	"net/http"
	"strings"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog home page !</h1>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>Hello, 这里是 goblog 404 page !</h1>")
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "此博客由 Go 语言编写，使用 Gin 框架，使用 GORM 操作数据库，使用 PostgreSQL 数据库。联系方式："+
		"<a href='mailto:yaosenjun168@live.cn'>yaosenjun168@live.cn</a>。")
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/", defaultHandler)
	router.HandleFunc("/about", aboutHandler)

	// 文章详情
	router.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		id := strings.SplitN(r.URL.Path, "/", 3)[2]
		fmt.Fprint(w, "文章 ID："+id)
	})

	// 列表 or 创建
	router.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fmt.Fprint(w, "显示文章列表")
		case "POST":
			fmt.Fprint(w, "创建新的文章")
		}
	})

	http.ListenAndServe(":3000", router)
}
