/*
 * @Date: 2023-01-31 09:02:33
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-31 09:07:51
 * @FilePath: /goblog/app/models/article/article.go
 */
package article

import (
	"github.com/go-sql-driver/mysql"
)

// Article is a struct for articles.
type Article struct {
	ID        uint64
	Title     string
	Body      string
	CreatedAt mysql.NullTime
	UpdateAt  mysql.NullTime
}
