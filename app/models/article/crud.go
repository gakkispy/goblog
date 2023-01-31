/*
 * @Date: 2023-01-31 09:08:50
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-31 13:39:56
 * @FilePath: /goblog/app/models/article/crud.go
 */
package article

import (
	"goblog/pkg/logger"
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

// GetAll 获取所有文章
func GetAll() ([]Article, error) {
	var articles []Article
	var err error

	if err = model.DB.Find(&articles).Error; err != nil {
		return articles, err
	}

	return articles, nil
}

// Create 创建文章
func (article *Article) Create() (err error) {
	result := model.DB.Create(&article).Error
	if err != nil {
		logger.LogError(err)
		return err
	}
	return result
}
