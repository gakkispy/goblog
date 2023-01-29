/*
 * @Date: 2023-01-16 14:28:24
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-29 16:11:28
 * @FilePath: /goblog/main.go
 */
package main

import (
	"database/sql"
	"fmt"
	"goblog/pkg/route"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var router *mux.Router
var db *sql.DB

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog home page !</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客由 Go 语言编写，使用 Gin 框架，使用 GORM 操作数据库，使用 PostgreSQL 数据库。联系方式："+
		"<a href='mailto:yaosenjun168@live.cn'>yaosenjun168@live.cn</a>。")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>Hello, 这里是 goblog 404 page !</h1>")
}

// Article 文章数据结构体
type Article struct {
	ID        int64
	Title     string
	Body      string
	CreatedAt mysql.NullTime
	UpdateAt  mysql.NullTime
}

func (a Article) Link() string {
	showURL, err := router.Get("articles.show").URL("id", strconv.FormatInt(a.ID, 10))
	if err != nil {
		checkError(err)
		return ""
	}

	return showURL.String()
}

// Int64ToString int64 转 string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := getRouteVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		}
	} else {
		// 4. 读取成功，显示文章
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteName2URL": route.Name2URL,
				"Int64ToString": Int64ToString,
			}).
			ParseFiles("resources/views/articles/show.gohtml")
		checkError(err)

		err = tmpl.Execute(w, article)
		checkError(err)
	}
}

func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
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
			checkError(err)
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
		checkError(err)

		err = tmpl.Execute(w, data)
		checkError(err)
	}
}

func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := getRouteVariable("id", r)

	// 2. 读取对应的文章数据
	_, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 未出现错误

		// 4.1 表单验证
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := validateArticleFormData(title, body)

		// 4.2 如果有错误
		if len(errors) > 0 {
			// 4.2.1 显示重新编辑文章的表单
			updateURL, _ := router.Get("articles.update").URL("id", id)
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				URL:    updateURL,
				Errors: errors,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			checkError(err)

			err = tmpl.Execute(w, data)
			checkError(err)
		} else {
			// 4.3 如果没有错误
			// 4.3.1 更新文章数据
			query := "UPDATE articles SET title = ?, body = ? WHERE id = ?"

			rs, err := db.Exec(query, title, body, id)

			if err != nil {
				checkError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "500 服务器内部错误")
			}

			// 4.3.2 重定向到文章详情页
			if n, _ := rs.RowsAffected(); n > 0 {
				showURL, _ := router.Get("articles.show").URL("id", id)
				http.Redirect(w, r, showURL.String(), http.StatusFound)
			} else {
				fmt.Fprintf(w, "您未做任何更改！")
			}
		}
	}
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 执行查询语句，返回一个结果集
	rows, err := db.Query("SELECT id, title, body FROM articles")
	checkError(err)
	defer rows.Close()

	var articles []Article
	// 2. 遍历结果集
	for rows.Next() {
		var article Article
		// 2.1 读取每一行数据
		err := rows.Scan(&article.ID, &article.Title, &article.Body)
		checkError(err)
		// 2.2 添加到 articles 切片中
		articles = append(articles, article)
	}

	// 2.3 检查遍历过程中是否出现错误
	err = rows.Err()
	checkError(err)

	// 3. 加载模板文件
	tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
	checkError(err)

	// 4. 渲染模板，并将数据传递给模板
	err = tmpl.Execute(w, articles)
	checkError(err)
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {

	storeURL, _ := router.Get("articles.store").URL()
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}

	tmpl, err :=
		template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}

}

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title  string
	Body   string
	URL    *url.URL
	Errors map[string]string
}

func validateArticleFormData(title string, body string) map[string]string {
	errors := make(map[string]string)
	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 ||
		utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需要介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	return errors
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)

	// 如果有错误，重新显示表单
	if len(errors) > 0 {

		storeURL, _ := router.Get("articles.store").URL()

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}

	} else {
		// 保存到数据库
		lastInsertID, err := saveArticleToDB(title, body)
		if lastInsertID > 0 {
			fmt.Fprintf(w, "文章创建成功，ID 为 "+strconv.FormatInt(lastInsertID, 10))
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	}

}

func saveArticleToDB(title string, body string) (int64, error) {
	// 变量初始化
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)

	// 1. 获取一个 prepare 声明语句
	stmt, err = db.Prepare("INSERT INTO articles (title, body, created_at, updated_at) VALUES (?, ?, ?, ?)")
	// 例行错误检测
	if err != nil {
		return 0, err
	}

	// 2. 在此函数运行结束后关闭此语句，防止占用 SQL 连接
	defer stmt.Close()

	// 3. 执行 SQL 语句， 传参进入绑定的内容
	rs, err = stmt.Exec(title, body, time.Now(), time.Now())
	// 例行错误检测
	if err != nil {
		return 0, err
	}

	// 4. 获取插入的 ID，若没有插入则返回 0，若有则返回自增 ID
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}

	return 0, err
}

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

func initDB() {
	var err error

	config := mysql.Config{
		User:                 "gakkispy",
		Passwd:               "secret",
		Addr:                 "127.0.0.1:3306",
		Net:                  "tcp",
		DBName:               "goblog",
		AllowNativePasswords: true,
	}

	// 连接数据库
	db, err = sql.Open("mysql", config.FormatDSN())
	checkError(err)

	// 设置最大连接数
	db.SetMaxOpenConns(25)

	// 设置最大空闲连接数
	db.SetMaxIdleConns(25)

	// 设置每个连接的过期时间
	db.SetConnMaxLifetime(5 * time.Minute)

	// 尝试连接数据库
	err = db.Ping()
	checkError(err)

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles (
		id BIGINT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
		title VARCHAR(255) COLLATE utf8mb4_unicode_ci NOT NULL,
		body LONGTEXT COLLATE utf8mb4_unicode_ci NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createArticlesSQL)
	checkError(err)
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
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		}
	} else {
		// 4. 未出现错误，执行删除操作
		rowsAffected, err := article.Delete()

		// 4.1 发生错误
		if err != nil {
			// 预计是 SQL 错误
			checkError(err)
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

func main() {
	// router := mux.NewRouter()
	initDB()
	createTables()

	route.InitializeRouter()
	router = route.Router

	router.HandleFunc("/", defaultHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	// 404
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 文章详情
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")

	// 列表 or 创建
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	// 更新，编辑
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")

	// 删除
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")
	// 中间件：强制 Content-Type 为 HTML
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
