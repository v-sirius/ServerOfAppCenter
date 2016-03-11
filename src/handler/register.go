package handler

import (
	"time"
	"net/http"
	"cache"
	
)

const longForm = "2006-01-02 15:04:05"

func Register(content map[string]interface{}, w http.ResponseWriter){
	res := ResponseInfo{ResponseType: "ret_register", Is_success: false}

	if _,exist:=content["useraccount"];!exist{
		res.Content=map[string]interface{}{"msg": "注册账号为空！","data": nil}
		write2Client(w, res)
		return
	}
	
	if _,exist:=content["password"];!exist{
		res.Content=map[string]interface{}{"msg": "注册时密码为空！", "data":nil}
		write2Client(w, res)
		return
	}
	
	
	if _, bExist := cache.G_CacheData.UserSet.IsExist(content["userid"].(string)); bExist {
		res.Content = map[string]interface{}{"msg": "此号码已注册","data":nil}
		write2Client(w, res)
		return
	}
	
	var tmpUser cache.User
	tmpUser.Account=content["useraccount"].(string)
	tmpUser.Password=content["password"].(string)
	// 当前时间
	currentTime := time.Now().Format(longForm)
	tmpUser.CreateTime=currentTime
	tmpUser.LoginFlag="0"
	tmpUser.Id=cache.G_CacheData.UserSet.GetLength()
	
	//将用户信息写入到cache
	cache.G_CacheData.UserSet.WriteUserFromServer2Cache(tmpUser)
	//将用户写入到数据库
	cache.G_CacheData.UserChan<-tmpUser
	
	res.Is_success = true
	res.Content = map[string]interface{}{"msg": "注册成功！", "data":map[string]interface{}{"useraccount":tmpUser.Account}}
	
	write2Client(w, res)
	return
}