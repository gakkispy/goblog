/*
 * @Date: 2023-01-29 11:28:58
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-29 16:10:43
 * @FilePath: /goblog/pkg/route/router.go
 */
package route

import "github.com/gorilla/mux"

// Router is a global variable for mux.Router.
var Router *mux.Router

// InitializeRouter initializes the global Router variable.
func InitializeRouter() {
	Router = mux.NewRouter()
}

// RouteName2URL is a map of route name to URL.
func Name2URL(routeName string, pairs ...string) string {
	url, err := Router.Get(routeName).URL(pairs...)
	if err != nil {
		return ""
	}

	return url.String()
}
