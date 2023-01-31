/*
 * @Date: 2023-01-16 14:28:24
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-31 14:13:22
 * @FilePath: /goblog/main.go
 */
package main

import (
	"database/sql"
	"fmt"
	"goblog/bootstrap"
	"goblog/pkg/database"
	"goblog/pkg/logger"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var router *mux.Router
var db *sql.DB

// Article 文章数据结构体
type Article struct {
	ID        int64
	Title     string
	Body      string
	CreatedAt mysql.NullTime
	UpdateAt  mysql.NullTime
}

func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body, &article.CreatedAt, &article.UpdateAt)
	return article, err
}

func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := getRouteVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 读物成功，显示表单
		updateURL, _ := router.Get("articles.update").URL("id", id)
		data := ArticlesFormData{
			Title:  article.Title,
			Body:   article.Body,
			URL:    updateURL,
			Errors: nil,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logger.LogError(err)

		err = tmpl.Execute(w, data)
		logger.LogError(err)
	}
}

// ArticlesFormData 文章表单数据
type ArticlesFormData struct {
	Title  string
	Body   string
	URL    *url.URL
	Errors map[string]string
}

// func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
// 	// 1. 获取 URL 参数
// 	id := getRouteVariable("id", r)

// 	// 2. 读取对应的文章数据
// 	_, err := getArticleByID(id)

// 	// 3. 如果出现错误
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			// 3.1 数据未找到
// 			w.WriteHeader(http.StatusNotFound)
// 			fmt.Fprintf(w, "404 文章未找到")
// 		} else {
// 			// 3.2 数据库错误
// 			logger.LogError(err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			fmt.Fprint(w, "500 服务器内部错误")
// 		}
// 	} else {
// 		// 4. 未出现错误

// 		// 4.1 表单验证
// 		title := r.PostFormValue("title")
// 		body := r.PostFormValue("body")

// 		errors := validateArticleFormData(title, body)

// 		// 4.2 如果有错误
// 		if len(errors) > 0 {
// 			// 4.2.1 显示重新编辑文章的表单
// 			updateURL, _ := router.Get("articles.update").URL("id", id)
// 			data := ArticlesFormData{
// 				Title:  title,
// 				Body:   body,
// 				URL:    updateURL,
// 				Errors: errors,
// 			}
// 			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
// 			logger.LogError(err)

// 			err = tmpl.Execute(w, data)
// 			logger.LogError(err)
// 		} else {
// 			// 4.3 如果没有错误
// 			// 4.3.1 更新文章数据
// 			query := "UPDATE articles SET title = ?, body = ? WHERE id = ?"

// 			rs, err := db.Exec(query, title, body, id)

// 			if err != nil {
// 				logger.LogError(err)
// 				w.WriteHeader(http.StatusInternalServerError)
// 				fmt.Fprintf(w, "500 服务器内部错误")
// 			}

// 			// 4.3.2 重定向到文章详情页
// 			if n, _ := rs.RowsAffected(); n > 0 {
// 				showURL, _ := router.Get("articles.show").URL("id", id)
// 				http.Redirect(w, r, showURL.String(), http.StatusFound)
// 			} else {
// 				fmt.Fprintf(w, "您未做任何更改！")
// 			}
// 		}
// 	}
// }

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 强制 Content-Type 为 HTML
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := getRouteVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		}
	} else {
		// 4. 未出现错误，执行删除操作
		rowsAffected, err := article.Delete()

		// 4.1 发生错误
		if err != nil {
			// 预计是 SQL 错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		} else {
			// 4.2 未发生错误
			if rowsAffected > 0 {
				// 4.2.1 影响行数大于 0，删除成功
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				// Edge case
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "404 文章未找到")
			}
		}
	}
}

func (a Article) Delete() (rowsAttected int64, err error) {
	rs, err := db.Exec("DELETE FROM articles WHERE id = " +
		strconv.FormatInt(a.ID, 10))

	if err != nil {
		return 0, err
	}

	// 删除成功
	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}

	return 0, nil
}

func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

func main() {
	// router := mux.NewRouter()
	database.Initialize()
	db = database.DB

	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()

	// 更新，编辑
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
	// router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")

	// 删除
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")
	// 中间件：强制 Content-Type 为 HTML
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
