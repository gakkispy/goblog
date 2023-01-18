/*
 * @Date: 2023-01-16 14:28:24
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-18 13:39:51
 * @FilePath: /goblog/main.go
 */
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var router = mux.NewRouter()

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

	storeURL, _ := router.Get("articles.store").URL()
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}

	tmpl, err :=
		template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}

}

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title  string
	Body   string
	URL    *url.URL
	Errors map[string]string
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := map[string]string{}

	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度必须介于 3 - 40 个字符之间"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if len(body) < 10 {
		errors["body"] = "内容不能少于 10 个字符"
	}

	// 如果有错误，重新显示表单
	if len(errors) > 0 {

		storeURL, _ := router.Get("articles.store").URL()

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}

	} else {
		// 显示验证通过
		fmt.Fprintf(w, "验证通过<br>")
		fmt.Fprintf(w, "标题：%s<br>", title)
		fmt.Fprintf(w, "内容：%s<br>", body)
		fmt.Fprintf(w, "标题的长度：%d<br>", len(title))
		fmt.Fprintf(w, "内容的长度：%d<br>", len(body))
	}

}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 强制 Content-Type 为 HTML
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

func main() {
	// router := mux.NewRouter()

	router.HandleFunc("/", defaultHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	// 404
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 文章详情
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")

	// 列表 or 创建
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	// 中间件：强制 Content-Type 为 HTML
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
