package handler

import (
	"cache"
	"net/http"
)

func Login(content map[string]interface{}, w http.ResponseWriter) {
	res := ResponseInfo{ResponseType: "ret_login", Is_success: false}

	if _, exist := content["useraccount"]; !exist {
		res.Content = map[string]interface{}{"msg": "登录账号为空！", "data": nil}
		write2Client(w, res)
		return
	}

	if _, exist := content["password"]; !exist {
		res.Content = map[string]interface{}{"msg": "登录时密码为空！", "data": nil}
		write2Client(w, res)
		return
	}

	//通过账号获取用户，若用户存在则err为空
	this_user, err := cache.G_CacheData.GetUserByAccount(content["useraccount"].(string))

	if err == nil { //先检查登录状态
		if this_user.LoginFlag == "1" {
			res.Is_success = true
			res.Content = map[string]interface{}{"msg": "该账号已经登录！", "data": map[string]interface{}{"userid": this_user.Account}}
			write2Client(w, res)
			return
		} else { ////若未登录检查密码是否匹配，匹配则登录成功，否则不成功
			if this_user.Password == content["password"].(string) {
				//密码匹配，此时修改登录标识位为1
				//todo:修改登录时间，修改登录标志位为1
				//修改cache中用户的登录标志位
				cache.G_CacheData.UserSet.ModifyLoginFlagInCache(this_user, "1")
				//修改数据库中用户的登录标志位
				flg := cache.G_CacheData.UserSet.ModifyLoginFlag2DbUser(this_user, "1")
				if flg {
					res.Is_success = true
					res.Content = map[string]interface{}{"msg": "登录成功！", "data": map[string]interface{}{"useraccount": this_user.Account}}
					write2Client(w, res)
					return
				} else {
					res.Is_success = false
					res.Content = map[string]interface{}{"msg": "数据库更新失败！", "data": nil}
					write2Client(w, res)
					return
				}
			} else {
				//密码不匹配
				res.Is_success = false
				res.Content = map[string]interface{}{"msg": "密码不匹配！", "data": nil}
				write2Client(w, res)
				return
			}
		}
	} else {
		res.Content = map[string]interface{}{"msg": "用户不存在！", "data": nil}
		res.Is_success = false
		write2Client(w, res)
		return
	}

}
