/*
 * @Date: 2023-01-30 14:52:32
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-30 14:54:59
 * @FilePath: /goblog/app/http/controllers/pages_controller.go
 */
package controllers

import (
	"fmt"
	"net/http"
)

// PagesController is a struct for static pages.
type PagesController struct{}

// HomeHandler handles requests to the / route
func (pc PagesController) HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog home page !</h1>")
}

// AboutHandler handles requests to the /about route
func (pc PagesController) AboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客由 Go 语言编写，使用 Gin 框架，使用 GORM 操作数据库，使用 PostgreSQL 数据库。联系方式："+
		"<a href='mailto:yaosenjun168@live.cn'>yaosenjun168@live.cn</a>。")
}

// NotFoundHandler handles requests to the 404 route
func (pc PagesController) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>Hello, 这里是 goblog 404 page !</h1>")
}
