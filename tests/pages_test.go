/*
 * @Date: 2023-01-29 09:37:45
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-29 09:41:36
 * @FilePath: /goblog/tests/pages_test.go
 */
package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomePage(t *testing.T) {
	baseURL := "http://localhost:3000"

	// 1. 请求 -- 模拟用户访问浏览器
	var (
		resp *http.Response
		err  error
	)
	resp, err = http.Get(baseURL + "/")

	// 2. 断言 -- 检查响应是否符合预期
	assert.NoError(t, err, "有错误发生，错误信息为：%v", err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "响应状态码不是 200")
}
