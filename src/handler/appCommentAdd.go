package handler

import (
	"time"
	"strconv"
	"net/http"
	"cache"
)

func AppComment(content map[string]interface{}, w http.ResponseWriter){
	res := ResponseInfo{ResponseType: "ret_app_comment", Is_success: false}
	
	if _, exist := content["useraccount"]; !exist {
		res.Content = map[string]interface{}{"msg": "账号为空！", "data": nil}
		write2Client(w, res)
		return
	}
	
	if _, exist := content["appid"]; !exist {
		res.Content = map[string]interface{}{"msg": "应用号为空！", "data": nil}
		write2Client(w, res)
		return
	}
	
	if _, exist := content["appcomment"]; !exist {
		res.Content = map[string]interface{}{"msg": "评价为空！", "data": nil}
		write2Client(w, res)
		return
	}
	
	if _, exist := content["appscore"]; !exist {
		res.Content = map[string]interface{}{"msg": "评分为空！", "data": nil}
		write2Client(w, res)
		return
	}
	
	var appCmt cache.AppComment
	appCmt.UserAccount=content["useraccount"].(string)
	appCmt.AppId,_=strconv.Atoi(content["appid"].(string))
	appCmt.AppComment=content["appcomment"].(string)
	appCmt.AppScore=content["appscore"].(float64)
	currentTime := time.Now().Format(longForm)
	appCmt.CommentTime=currentTime
	appCmt.CommentId=cache.G_CacheData.AppCommentSet.GetLength()
	
	//将用户评论和评分写入到cache
	cache.G_CacheData.AppCommentSet.WriteAppCommentFromServer2Cache(appCmt)
	//将用户评论和评分写入到数据库
	cache.G_CacheData.AppCommentChan<-appCmt
	
	res.Content = map[string]interface{}{"msg": "评价成功！", "data": map[string]interface{}{"useraccount": appCmt.UserAccount,"appid":appCmt.AppId,"appcommentid":appCmt.CommentId,"appcomment":appCmt.AppComment,"appscore":appCmt.AppScore,"appcommenttime":appCmt.CommentTime}}
	res.Is_success = true
	write2Client(w, res)
	return
}