package net

import (
	"GoStatusClient/logger"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func PostRequest(url string, params string) (string, error) {
	// 1. 创建http客户端实例
	client := &http.Client{}
	// 2. 创建请求实例
	req, err := http.NewRequest("POST", url, strings.NewReader(params))
	if err != nil {
		return "", err
	}
	// 3. 设置请求头，可以设置多个
	req.Header.Set("Host", " ")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 4. 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Warning("Close POST request error", err)
			return
		}
	}(resp.Body)

	// 5. 一次性读取请求到的数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
