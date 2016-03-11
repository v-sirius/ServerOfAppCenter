package cache

import (
	"db"
	"errors"
	"fmt"
	"strings"
	"sync"
)

// app表
type App struct {
	AppId           int    //
	AppName         string //
	AppLogoUrl      string //
	AppClass        string //
	AppInfo         string //
	AppImgUrl       string //
	AppSize         float32
	AppUrl          string
	AppScore        float32
	AppVersion      string
	AppProvider     string
	AppUpTime       string
	AppDownloadTime int
	AppOnlineFlg    bool
}

// 所有应用表
type AppSet struct {
	allApp []App
	//appId2App map[int]int
	appName2App  map[string]int //
	appClass2App map[string][]int

	rwLock sync.RWMutex // 读写锁
}

// 初始化应用表缓存
func InitAppSet() *AppSet {
	pRet := new(AppSet)
	pRet = &AppSet{allApp: make([]App, 0), appName2App: make(map[string]int, 0), appClass2App: make(map[string][]int, 0)}

	return pRet
}

// 加载数据
func (appSet *AppSet) LoadData(pDb *db.DbOperation) {
	var app App
	selectstr := "select appId,appName,appLogoUrl,appClass,appInfo,appImgUrl,appSize,appUrl,appScore,appVersion,appProvider,appUpTime,appDownloadTime,appOnlineFlg from app"
	rows := pDb.Find(selectstr)

	for rows.Next() {
		err := rows.Scan(&app.AppId, &app.AppName, &app.AppLogoUrl, &app.AppClass, &app.AppInfo, &app.AppImgUrl, &app.AppSize, &app.AppUrl, &app.AppScore, &app.AppVersion, &app.AppProvider, &app.AppUpTime, &app.AppDownloadTime, &app.AppOnlineFlg)
		if err != nil {
			panic("error in:scanning the table app")
		}

		appSet.WriteApp2Cache(app)
	}
}

// 从数据库读取数据写入到cache
func (appSet *AppSet) WriteApp2Cache(app App) {
	appSet.rwLock.Lock()

	size := len(appSet.allApp)

	fmt.Println("-----------size----------")
	fmt.Println(size)
	appSet.allApp = append(appSet.allApp, app)
	appSet.appName2App[app.AppName] = size
	appSet.appClass2App[app.AppClass] = append(appSet.appClass2App[app.AppClass], size)

	func() {
		appSet.rwLock.Unlock()
	}()
}

//通过id获取app
func (appSet *AppSet) GetAppById(id int) (App, error) {

	if appSet.allApp[id].AppOnlineFlg == true {
		return appSet.allApp[id], nil
	} else {
		return App{AppId: -1}, errors.New("该应用已下线！")
	}

}

//读取搜索应用匹配的app
func (appSet *AppSet) SearchApps(t_class string, t_content string, t_count string, t_order string) ([]App, error) {
	var partApp []App
	//todo:比较操作
	//在全部分类下查找应用：all+空（表示在首页查看全部应用）或者all+content（表示全部分类下查找某应用）
	//在某个分类下查找应用：class+空（表示查找某个分类的应用）或者class+content（表示某分类下查找某应用）
	if t_class == "all" {
		//全部分类下查找应用,名称中包含关键字、简介里包含关键字
		for _, app := range appSet.allApp {
			if strings.Contains(app.AppName, t_content) {
				partApp = append(partApp, app)
				fmt.Println(app)
			}
		}
	} else { //某分类下查找应用
		for _, app_id := range appSet.appClass2App[t_class] {
			//判断应用名称中是否含有相似关键字
			if strings.Contains(appSet.allApp[app_id].AppName, t_content) {
				partApp = append(partApp, appSet.allApp[app_id])
			}
		}
	}

	if len(partApp) == 0 {
		return partApp, errors.New("not exist related app data")
	}
	return partApp, nil
}

////channel写入数据库
//func (appSet *AppSet) Write2DbAgenda() {
//	var agendachan = G_CacheData.AppSet.AgendaChan
//	for app := range agendachan {
//		var sql = "insert into scheduleinfo(id,uid,stitle,stime,scontent,screattime)" +
//			" values ('" + strconv.Itoa(app.Id) + "','" + app.Title + "','" + app.AgendaTime + "','" +
//			app.Content + "','" + app.CreateTime + "')"

//		if !db.G_db.Insert2Table(sql) {
//			fmt.Println("agenda表插入不成功！")
//		}

//	}

//}
