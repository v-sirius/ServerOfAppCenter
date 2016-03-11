package handler

import (
	"cache"
	"net/http"
)

func CheckUserLoginFlag(content map[string]interface{}, w http.ResponseWriter) {
	res := ResponseInfo{ResponseType: "ret_check_user_loginflag", Is_success: false}

	if _, exist := content["useraccount"]; !exist {
		res.Content = map[string]interface{}{"msg": "登录账号为空！", "data": nil}
		write2Client(w, res)
		return
	}

	//通过账号获取用户，若用户存在则err为空
	this_user, err := cache.G_CacheData.GetUserByAccount(content["useraccount"].(string))

	if err == nil { //先检查登录状态
		if this_user.LoginFlag == "1" {
			//用户已登录
				res.Is_success = true
				res.Content = map[string]interface{}{"msg": "用户已登录！", "data": map[string]interface{}{"loginflag":this_user.LoginFlag}}
				write2Client(w, res)
				return
			
		} else { ////若未登录提示未登录
			res.Is_success = false
			res.Content = map[string]interface{}{"msg": "用户未登录！", "data": map[string]interface{}{"loginflag":this_user.LoginFlag}}
			write2Client(w, res)
			return
		}
	} else {
		res.Content = map[string]interface{}{"msg": "用户不存在！", "data": nil}
		res.Is_success = false
		write2Client(w, res)
		return
	}
}
