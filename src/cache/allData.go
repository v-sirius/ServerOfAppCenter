package cache

import (
	"strings"
	"db"
	"log"
)

// cache中所有的数据结构
// main程序启动后，须初始化这些结构
type AllCacheData struct {
	AppSet
	AppCommentSet
	UserSet
	LoginHistorySet
	
//	ChString chan string

}

// 所有数据结构
var G_CacheData *AllCacheData

// 初始化函数－－初始化所有的缓存数据结构
func InitAllCacheData() *AllCacheData {
	pRet := new(AllCacheData)

	pRet.AppSet = *InitAppSet()
	pRet.AppCommentSet = *InitAppCommentSet()
	pRet.UserSet = *InitUserSet()
	pRet.LoginHistorySet = *InitLoginHistorySet()
	
	//pRet.ChString = make(chan string, 1000)

	//pRet.Stop = false

	return pRet
}

// 加载所有用户数据
func LoadAllData(pRet *AllCacheData, pDb *db.DbOperation) {
	pRet.AppSet.LoadData(pDb)
	pRet.AppCommentSet.LoadData(pDb)
	pRet.UserSet.LoadData(pDb)
	pRet.LoginHistorySet.LoadData(pDb)
}

// 开启goroute
func (pAll *AllCacheData) GoCache() {
//	fmt.Println("entering GoChe!")
//	//log.Println(pAll.RealTimePositionSet)
//	log.Println(pAll.ChPos)

//	//go pAll.Buffer4CachePosition.WriteData2Buffer(&pAll.RealTimePositionSet, pAll.ChPos)
//	//go pAll.Buffer4CachePosition.Write2CacheAndDbBuffer(&pAll.Buffer4Db, &pAll.RealTimePositionSet)
//	go pAll.Write2CacheAndDbBuffer()
//	go pAll.WritePosDb()
//	//go WriteDb() //位置写入数据库
//	log.Println("out Gocache")
}

// 从缓存写入数据－－同时写入cache和db缓存
func (pAll *AllCacheData) Write2CacheAndDbBuffer(/*chPos <-chan UserPosition*/) {
	
}

// 写入数据库
func (pAll *AllCacheData)WritePosDb(){
	
}

// 处理数据库更新或删除
//func (pAll *AllCacheData) UpdateDb(){
//	for v := range pAll.ChString{
//		if db.G_db.Insert2Table(v) == false {
//			log.Println("update error: ", v)
//		}
//	}
//}

//接口：sql语句中特殊字符处理--encode
func SqlEncode(msg string)string{
	log.Println("----------input of sqlEncode():",msg)
	
	msgAfterEncode:=strings.Replace(msg,"'","&#39",-1)
	
	msgAfterEncode=strings.Replace(msgAfterEncode,"\"","&#34",-1)
	msgAfterEncode=strings.Replace(msgAfterEncode,"=","&#61",-1)
	msgAfterEncode=strings.Replace(msgAfterEncode,"-","&#45",-1)
	msgAfterEncode=strings.Replace(msgAfterEncode,";","&#59",-1)
	msgAfterEncode=strings.Replace(msgAfterEncode,"exec","ＥＸＥＣ",-1)
	msgAfterEncode=strings.Replace(msgAfterEncode,"or","ＯＲ",-1)
	msgAfterEncode=strings.Replace(msgAfterEncode,"and","ＡＮＤ",-1)
	
	log.Println("----------ouput after encode:",msgAfterEncode)
	return msgAfterEncode
}

//接口：sql语句中特殊字符处理--decode
func SqlDecode(msg string)string{
	log.Println("----------input of sqlDecode():",msg)
	
	msgAfterDecode:=strings.Replace(msg,"&#39","'",-1)
	msgAfterDecode=strings.Replace(msgAfterDecode,"&#34","\"",-1)
	msgAfterDecode=strings.Replace(msgAfterDecode,"&#61","=",-1)
	msgAfterDecode=strings.Replace(msgAfterDecode,"&#45","-",-1)
	msgAfterDecode=strings.Replace(msgAfterDecode,"&#59",";",-1)
	msgAfterDecode=strings.Replace(msgAfterDecode,"ＥＸＥＣ","exec",-1)
	msgAfterDecode=strings.Replace(msgAfterDecode,"ＯＲ","or",-1)
	msgAfterDecode=strings.Replace(msgAfterDecode,"ＡＮＤ","and",-1)
	
	log.Println("----------output after decode:",msgAfterDecode)
	return msgAfterDecode
}