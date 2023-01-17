/*
 * @Date: 2023-01-16 14:28:24
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-17 17:01:41
 * @FilePath: /goblog/main.go
 */
package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog home page !</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客由 Go 语言编写，使用 Gin 框架，使用 GORM 操作数据库，使用 PostgreSQL 数据库。联系方式："+
		"<a href='mailto:yaosenjun168@live.cn'>yaosenjun168@live.cn</a>。")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>Hello, 这里是 goblog 404 page !</h1>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "文章 ID："+id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "显示文章列表")
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "创建新的文章")
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 强制 Content-Type 为 HTML
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", defaultHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	// 404
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 文章详情
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")

	// 列表 or 创建
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesCreateHandler).Methods("POST").Name("articles.create")

	// 中间件：强制 Content-Type 为 HTML
	router.Use(forceHTMLMiddleware)

	// 生成 URL
	homeURL, _ := router.Get("home").URL()
	fmt.Println(homeURL)
	aboutURL, _ := router.Get("about").URL()
	fmt.Println(aboutURL)
	articlesShowURL, _ := router.Get("articles.show").URL("id", "1")
	fmt.Println(articlesShowURL)

	http.ListenAndServe(":3000", router)
}
