package handler

import (
	"cache"
	"net/http"
)

func Logout(content map[string]interface{}, w http.ResponseWriter) {
	res := ResponseInfo{ResponseType: "ret_login", Is_success: false}

	if _, exist := content["useraccount"]; !exist {
		res.Content = map[string]interface{}{"msg": "登录账号为空！", "data": nil}
		write2Client(w, res)
		return
	}

	//通过账号获取用户，若用户存在则err为空
	this_user, err := cache.G_CacheData.GetUserByAccount(content["useraccount"].(string))

	if err == nil { //先检查登录状态
		if this_user.LoginFlag == "1" {
			//用户已登录，此时可以注销，修改登录标志位为0
			//todo:修改注销时间
			//修改标识位为0
			//修改cache中用户的登录标志位
			cache.G_CacheData.UserSet.ModifyLoginFlagInCache(this_user, "0")
			//修改数据库中用户的登录标志位
			flg := cache.G_CacheData.UserSet.ModifyLoginFlag2DbUser(this_user, "0")
			
			//是否需要返回客户端信息？？
			if flg {
				res.Is_success = true
				res.Content = map[string]interface{}{"msg": "注销成功！", "data": nil}
				write2Client(w, res)
				return
			} else {
				res.Is_success = false
				res.Content = map[string]interface{}{"msg": "数据库更新失败！", "data": nil}
				write2Client(w, res)
				return
			}
		} else { ////若未登录提示未登录
			res.Is_success = false
			res.Content = map[string]interface{}{"msg": "用户未登录！", "data": nil}
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
