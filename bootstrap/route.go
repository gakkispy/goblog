/*
 * @Date: 2023-01-30 15:57:46
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-30 15:59:58
 * @FilePath: /goblog/bootstrap/route.go
 */
package bootstrap

import (
	"goblog/routes"

	"github.com/gorilla/mux"
)

// SetupRoute registers routes.
func SetupRoute() *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)

	return router
}
