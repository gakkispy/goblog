/*
 * @Date: 2023-01-30 14:13:54
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-31 14:29:41
 * @FilePath: /goblog/routes/web.go
 */
package routes

import (
	"goblog/app/http/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterWebRoutes registers web routes.
func RegisterWebRoutes(r *mux.Router) {
	// Static pages
	pc := new(controllers.PagesController)
	r.HandleFunc("/", pc.HomeHandler).Methods("GET").Name("home")
	r.HandleFunc("/about", pc.AboutHandler).Methods("GET").Name("about")
	r.NotFoundHandler = http.HandlerFunc(pc.NotFoundHandler)

	// Articles
	ac := new(controllers.ArticlesController)
	r.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")
	r.HandleFunc("/articles", ac.Store).Methods("POST").Name("articles.store")
	r.HandleFunc("/articles/create", ac.Create).Methods("GET").Name("articles.create")

	r.HandleFunc("/articles/{id:[0-9]+}/edit", ac.Edit).Methods("GET").Name("articles.edit")
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Update).Methods("POST").Name("articles.update")

}
