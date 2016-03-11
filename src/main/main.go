package main

import (
	"cache"
	"db"
	"handler"
	"log"
	"net/http"
)

func main() {
	db.G_db = db.InitDbOperation(`/test?charset=utf8`)
	db.G_db.Open()
	defer db.G_db.Close()
	//db.G_db.CreateTable()

	cache.G_CacheData = cache.InitAllCacheData()

	go cache.G_CacheData.UserSet.Write2DbUser()
	go cache.G_CacheData.AppCommentSet.Write2DbAppComment()
	
	http.HandleFunc("/software_center", handler.HandleRequest)

	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
