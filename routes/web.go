/*
 * @Date: 2023-01-30 14:13:54
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-30 15:21:12
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
}
