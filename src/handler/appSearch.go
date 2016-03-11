package handler

import (
	"cache"
	"net/http"
)

func AppSearch(content map[string]interface{}, w http.ResponseWriter) {
	res := ResponseInfo{ResponseType: "ret_app_search", Is_success: false}

	//查找应用的策略
	//在全部分类下查找应用：all+空（表示在首页查看全部应用）或者all+content（表示全部分类下查找某应用）
	//在某个分类下查找应用：class+空（表示查找某个分类的应用）或者class+content（表示某分类下查找某应用）
	if _, exist := content["appclass"]; !exist {
		res.Content = map[string]interface{}{"msg": "分类不能为空！", "data": nil}
		write2Client(w, res)
		return
	}

	tmp_content := content["searchcontent"].(string)
	tmp_class := content["appclass"].(string)
	tmp_count := content["appcount"].(string)
	tmp_order := content["apporder"].(string)

	tmp_apps, err := cache.G_CacheData.AppSet.SearchApps(tmp_class, tmp_content, tmp_count, tmp_order)
	if err == nil {
		res.Content = map[string]interface{}{"msg": "查询成功！ ", "data": map[string]interface{}{"apps": tmp_apps}}
		res.Is_success = true
		write2Client(w, res)
		return
	} else {
		res.Content = map[string]interface{}{"msg": "未能查找到相关应用！ ", "data": map[string]interface{}{"apps": tmp_apps}}
		res.Is_success = false
		write2Client(w, res)
		return
	}

}
