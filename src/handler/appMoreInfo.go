package handler

import (
	"cache"
	"net/http"
	"strconv"
)

func AppMoreInfo(content map[string]interface{}, w http.ResponseWriter) {
	res := ResponseInfo{ResponseType: "ret_app_more_info", Is_success: false}

	if _, exist := content["appid"]; !exist {
		res.Content = map[string]interface{}{"msg": "应用id为空！", "data": nil}
		write2Client(w, res)
		return
	}

	app_id, _ := strconv.Atoi(content["appid"].(string))
	this_app, err := cache.G_CacheData.AppSet.GetAppById(app_id)
	if err == nil {
		res.Content = map[string]interface{}{"msg": "详情获取成功！", "data": map[string]interface{}{"appid": this_app.AppId,"appname":this_app.AppName,"applogo":this_app.AppLogoUrl,"appclass":this_app.AppClass,"appinfo":this_app.AppInfo,"appscreenshot":this_app.AppImgUrl,"appversion":this_app.AppVersion,"appscore":this_app.AppScore}}
		res.Is_success = true
		write2Client(w, res)
		return
	}else{
		res.Is_success=false
		res.Content= map[string]interface{}{"msg": "详情获取失败！", "data": nil}
	}

}
