package handler

import (
	"fmt"
	"cache"
	"net/http"
	"strconv"
)

func AppCommentDel(content map[string]interface{}, w http.ResponseWriter) {
	res := ResponseInfo{ResponseType: "ret_app_comment_del", Is_success: false}
	if _, exist := content["appcommentid"]; !exist {
		res.Content = map[string]interface{}{"msg": "应用评价id为空！", "data": nil}
		write2Client(w, res)
		return
	}

	cmtId, _ := strconv.Atoi(content["appcommentid"].(string))
	//通过评价id获取评价
	this_cmt, err := cache.G_CacheData.GetCommentById(cmtId)
	fmt.Println("this_cmt:",this_cmt)
	
	if err == nil {
		cache.G_CacheData.AppCommentSet.ModifyCmtDelFlagInCache(cmtId, "1")
		V := cache.G_CacheData.AppCommentSet.ModifyCmtDelFlag2DbAppCmt(cmtId, "1")
		if V {
			res.Is_success = true
			res.Content = map[string]interface{}{"msg": "删除成功！ ", "data": nil}
			write2Client(w, res)
			return
		}else{
			res.Is_success = false
			res.Content = map[string]interface{}{"msg": "删除失败！ ", "data": nil}
			write2Client(w, res)
			return
		}
	}
}
