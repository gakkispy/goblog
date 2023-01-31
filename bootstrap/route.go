/*
 * @Date: 2023-01-30 15:57:46
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-31 09:27:40
 * @FilePath: /goblog/bootstrap/route.go
 */
package bootstrap

import (
	"goblog/pkg/route"
	"goblog/routes"

	"github.com/gorilla/mux"
)

// SetupRoute registers routes.
func SetupRoute() *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)

	route.SetupRoute(router)

	return router
}
