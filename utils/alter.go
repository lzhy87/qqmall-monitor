package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const webhook_url = "https://oapi.dingtalk.com/robot/send?access_token=015791f1acee09602b550414e49da59e19262c60e637169641b4087d825a352d"

func DingToInfo(s string) bool {
	content, data := make(map[string]string), make(map[string]interface{})
	content["content"] = s
	data["msgtype"] = "text"
	data["text"] = content
	b, _ := json.Marshal(data)
	resp, err := http.Post(webhook_url,
		"application/json",
		bytes.NewBuffer(b))
	fmt.Println(resp)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return true
}
