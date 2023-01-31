/*
 * @Date: 2023-01-29 11:28:58
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-31 10:06:20
 * @FilePath: /goblog/pkg/route/router.go
 */
package route

import (
	"goblog/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

// route is a global variable for *mux.Router.
var route *mux.Router

// SetupRoute registers routes.
func SetupRoute(r *mux.Router) {
	route = r
}

// RouteName2URL is a map of route name to URL.
func Name2URL(routeName string, pairs ...string) string {
	url, err := route.Get(routeName).URL(pairs...)
	if err != nil {
		logger.LogError(err)
		return ""
	}

	return url.String()
}

// GetRouteVariable 获取 URI 路由参数
func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
