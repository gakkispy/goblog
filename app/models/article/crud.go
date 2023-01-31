/*
 * @Date: 2023-01-31 09:08:50
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-31 09:10:27
 * @FilePath: /goblog/app/models/article/crud.go
 */
package article

import (
	"goblog/pkg/model"
	"goblog/pkg/types"
)

// Get 通过 ID 获取文章
func Get(idstr string) (Article, error) {
	var article Article
	var err error

	if idstr != "" {
		id := types.StringToUint64(idstr)
		if err = model.DB.First(&article, id).Error; err != nil {
			return article, err
		}
	}

	return article, nil
}
