/*
 * @Date: 2023-01-29 09:37:45
 * @LastEditors: gakkispy && yaosenjun168@live.cn
 * @LastEditTime: 2023-01-29 10:50:16
 * @FilePath: /goblog/tests/pages_test.go
 */
package tests

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllPages(t *testing.T) {
	baseURL := "http://localhost:3000"

	// 1. 声明加初始化测试数据
	var tests = []struct {
		url      string
		method   string
		expected int
	}{
		{"/", "GET", http.StatusOK},
		{"/about", "GET", http.StatusOK},
		{"/notfound", "GET", http.StatusNotFound},
		{"/articles", "GET", http.StatusOK},
		{"/articles/2", "GET", http.StatusOK},
		{"/articles/2/edit", "GET", http.StatusOK},
		{"/articles/create", "GET", http.StatusOK},
		{"/articles/2", "POST", http.StatusOK},
		{"/articles/", "POST", http.StatusOK},
		{"/articles/1/delete", "POST", http.StatusNotFound},
	}

	// 2. 遍历测试数据
	for _, test := range tests {
		t.Logf("当前测试的URL: %v \n", test.url)
		var (
			resp *http.Response
			err  error
		)

		// 2.1 发起请求
		switch {
		case test.method == "POST":
			data := make(map[string][]string)
			resp, err = http.PostForm(baseURL+test.url, data)

		default:
			resp, err = http.Get(baseURL + test.url)
		}

		// 2.2 断言
		assert.NoError(t, err, "请求 "+test.url+" 时有错误")
		assert.Equal(t, test.expected, resp.StatusCode, "请求 "+test.url+" 时状态码错误，应为 "+strconv.Itoa(test.expected)+"")
	}
}
