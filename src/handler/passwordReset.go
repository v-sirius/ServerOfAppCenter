package handler

import (
	"cache"
	"fmt"
	"net/http"
)

func PasswordReset(content map[string]interface{}, w http.ResponseWriter) {
	res := ResponseInfo{ResponseType: "ret_check_user_loginflag", Is_success: false}

	if _, exist := content["useraccount"]; !exist {
		res.Content = map[string]interface{}{"msg": "账号为空！", "data": nil}
		write2Client(w, res)
		return
	}

	if _, exist := content["newpassword"]; !exist {
		res.Content = map[string]interface{}{"msg": "新密码为空！", "data": nil}
		write2Client(w, res)
		return
	}

	//通过账号获取用户，若用户存在则err为空
	//用户存在，用新密码修改旧密码
	this_user, err := cache.G_CacheData.GetUserByAccount(content["useraccount"].(string))

	fmt.Println(this_user)
	fmt.Println(err)

	if err == nil {
		cache.G_CacheData.UserSet.PasswordResetInCache(this_user, content["newpassword"].(string))
		flg := cache.G_CacheData.UserSet.PasswordReset2DbUser(this_user, content["newpassword"].(string))
		if flg {
			res.Is_success = true
			res.Content = map[string]interface{}{"msg": "密码修改成功！", "data": map[string]interface{}{"useraccount": this_user.Account}}
			write2Client(w, res)
			return
		} else {
			res.Is_success = false
			res.Content = map[string]interface{}{"msg": "密码修改失败！", "data": map[string]interface{}{"useraccount": this_user.Account}}
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
