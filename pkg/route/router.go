/*
 * @Date: 2023-01-29 11:28:58
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-30 14:49:50
 * @FilePath: /goblog/pkg/route/router.go
 */
package route

import (
	"goblog/routes"
	"net/http"

	"github.com/gorilla/mux"
)

// Router is a global variable for mux.Router.
var Router *mux.Router

// InitializeRouter initializes the global Router variable.
func InitializeRouter() {
	Router = mux.NewRouter()
	routes.RegisterWebRoutes(Router)
}

// RouteName2URL is a map of route name to URL.
func Name2URL(routeName string, pairs ...string) string {
	url, err := Router.Get(routeName).URL(pairs...)
	if err != nil {
		return ""
	}

	return url.String()
}

// GetRouteVariable 获取 URI 路由参数
func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
