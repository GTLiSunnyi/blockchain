package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Request(urls string) {
	var data = map[string]interface{}{}
	res := 1
	var resp *http.Response
	var err error

	//-------------------------发起请求------------------------------------
	switch res {
	case 1:
		// get请求
		resp, err = http.Get(urls)
	case 2:
		// post发送json
		datas, _ := json.Marshal(data)
		resp, err = http.Post(urls, "application/json", bytes.NewBuffer(datas))
	case 3:
		// post发送from表单
		datas := make(url.Values)
		for key, value := range data {
			datas[string(key)] = []string{value.(string)}
		}
		resp, err = http.PostForm(urls, datas)
	case 4:
		// 创建请求设置参数
		datas := url.Values{
			"name": {"fei"},
		}
		req, err := http.NewRequest("POST", urls, strings.NewReader(datas.Encode()))
		if err != nil {
			// handle error
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")

		client := &http.Client{}
		resp, err = client.Do(req)

	}
	//-------------------------发起请求成功------------------------------------

	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(body))
}
