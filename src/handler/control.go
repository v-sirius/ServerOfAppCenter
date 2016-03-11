package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RequestInfo struct {
	RequestType string                 `json:"type"`
	Content     map[string]interface{} `json:"content"`
}

type ResponseInfo struct {
	ResponseType string                 `json:"type"`
	Is_success   bool                   `json:"is_success"`
	Content      map[string]interface{} `json:"content"`
}

// 处理接收到的所有请求
func HandleRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)

		fmt.Println("request body:", string(body))

		if err != nil {
			// todo:需添加日志
			panic("error in reading request's body when handling request")
		}

		// 定义一个接收的结构变量
		var requestBody RequestInfo
		err = json.Unmarshal(body, &requestBody)

		fmt.Println("json Unmashal 错误:", err)
		fmt.Println("json Unmashal 内容：", requestBody)

		if err != nil {
			// todo:需要添加日志
			panic("format error in request body")
		}

		// 根据type种类分别进行处理
		switch requestBody.RequestType {

		// 1.新用户注册
		case "register":
			Register(requestBody.Content, w)

		//2.用户登录
		case "login":
			Login(requestBody.Content, w)

		//3.用户注销
		case "logout":
			Logout(requestBody.Content, w)

		//4.评价应用
		case "app_comment":
			AppComment(requestBody.Content, w)

		//5.删除评价

		case "app_comment_del":
			AppCommentDel(requestBody.Content, w)

		//6.查找软件
		case "app_search":
			AppSearch(requestBody.Content, w)
		//7.查看应用详情
		case "app_more_info":
			AppMoreInfo(requestBody.Content, w)

		//8.下载应用
		case "app_download":
			AppDownload(requestBody.Content, w)
			
		//9.忘记密码，更改密码
		case "password_reset":
			PasswordReset(requestBody.Content, w)
		
		//10.检查登录状态
		case "check_user_loginflag":
			CheckUserLoginFlag(requestBody.Content, w)
		}
	}
}

// 写返回给客户端的数据
func write2Client(w http.ResponseWriter, val ResponseInfo) {
	bytes, err := json.Marshal(val)
	if err != nil {
		// todo: add log
		panic("error in encoding json -- func:write2Client")
	}
	if _, err = w.Write(bytes); err != nil {
		// todo: add log
		panic("error in writing response -- func:write2Client")
	}
}
