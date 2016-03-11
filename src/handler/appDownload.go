package handler

import (
	"net/http"
	
)

func AppDownload(content map[string]interface{}, w http.ResponseWriter){
	res := ResponseInfo{ResponseType: "ret_app_download", Is_success: false}

	
	res.Content = map[string]interface{}{"msg": " ", "data": map[string]interface{}{"data": nil}}
	res.Is_success = true
	write2Client(w, res)
	return
}