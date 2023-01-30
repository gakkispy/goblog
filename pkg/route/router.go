/*
 * @Date: 2023-01-29 11:28:58
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-30 16:12:00
 * @FilePath: /goblog/pkg/route/router.go
 */
package route

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RouteName2URL is a map of route name to URL.
func Name2URL(routeName string, pairs ...string) string {
	var router = mux.NewRouter()
	url, err := router.Get(routeName).URL(pairs...)
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
