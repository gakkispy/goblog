/*
 * @Date: 2023-01-30 15:29:30
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-31 10:37:23
 * @FilePath: /goblog/app/http/controllers/articles_controller.go
 */
package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"html/template"
	"net/http"

	"gorm.io/gorm"
)

// ArticlesController is a struct for articles.
type ArticlesController struct{}

// IndexHandler handles requests to the /articles route
// func (ac ArticlesController) IndexHandler(w http.ResponseWriter, r *http.Request) {
// 	// 1. 执行查询语句，返回一个结果集
// 	rows, err := db.Query("SELECT id, title, body FROM articles")
// 	logger.LogError(err)
// 	defer rows.Close()

// 	var articles []Article
// 	// 2. 遍历结果集
// 	for rows.Next() {
// 		var article Article
// 		// 2.1 读取每一行数据
// 		err := rows.Scan(&article.ID, &article.Title, &article.Body)
// 		logger.LogError(err)
// 		// 2.2 添加到 articles 切片中
// 		articles = append(articles, article)
// 	}

// 	// 2.3 检查遍历过程中是否出现错误
// 	err = rows.Err()
// 	logger.LogError(err)

// 	// 3. 加载模板文件
// 	tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
// 	logger.LogError(err)

// 	// 4. 渲染模板，并将数据传递给模板
// 	err = tmpl.Execute(w, articles)
// 	logger.LogError(err)
// }

// ShowHandler handles requests to the /articles/{id} route
func (ac ArticlesController) ShowHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		}
	} else {
		// 4. 读取成功，显示文章
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteName2URL":  route.Name2URL,
				"Uint64ToString": types.Uint64ToString,
			}).
			ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)

		err = tmpl.Execute(w, article)
		logger.LogError(err)
	}
}
