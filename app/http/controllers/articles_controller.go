/*
 * @Date: 2023-01-30 15:29:30
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-31 14:25:08
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
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// ArticlesController is a struct for articles.
type ArticlesController struct{}

// Article.Index handles requests to the /articles route
func (ac ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	// 1. ArticlesModel 取出所有文章数据
	articles, err := article.GetAll()

	if err != nil {
		// 1.1 如果出错，记录错误信息
		logger.LogError(err)
		// 1.2 返回 500 响应
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		// 1.3 加载模板文件
		tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
		logger.LogError(err)
		// 1.4 渲染模板，并将数据传递给模板
		err = tmpl.Execute(w, articles)
	}

}

// Article.Show handles requests to the /articles/{id} route
func (ac ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
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

// ArticlesFormData is a struct for articles form data.
type ArticlesFormData struct {
	Title  string
	Body   string
	URL    string
	Errors map[string]string
}

// Article.Create handles requests to the /articles/create route
func (ac ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	storeURL := route.Name2URL("articles.store")
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

// validateArticleFormData validates the form data from the /articles/create route
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

// Article.Store handles requests to the /articles/store route
func (ac ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)

	// 如果有错误，重新显示表单
	if len(errors) > 0 {

		storeURL := route.Name2URL("articles.store")

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
		_article := article.Article{
			Title:     title,
			Body:      body,
			CreatedAt: mysql.NullTime{Time: time.Now(), Valid: true},
			UpdatedAt: mysql.NullTime{Time: time.Now(), Valid: true},
		}
		_article.Create()

		if _article.ID > 0 {
			fmt.Fprintf(w, "文章创建成功，ID 为 "+strconv.FormatInt(int64(_article.ID), 10))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "文章创建失败，请联系管理员.")
		}
	}

}
