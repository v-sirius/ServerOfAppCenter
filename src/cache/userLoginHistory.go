package cache

import (
	"fmt"
	"db"
	"sync"
)

// LoginHistory
type LoginHistory struct {
	dataId         int    // 
	Account string //用户账号
	//loginWay     int //0:客户端	1：web
	LoginTime string //
	LogoutTime string // 
	Login    bool // 登录标志位  0:未登录 1：已登录
	
}

// 登录历史表
type LoginHistorySet struct {
	allLoginHistory []LoginHistory
	account2dataId map[string]int 

	rwLock sync.RWMutex // 读写锁
}

// 初始化登陆历史表缓存
func InitLoginHistorySet() *LoginHistorySet {
	pRet := new(LoginHistorySet)
	pRet = &LoginHistorySet{allLoginHistory: make([]LoginHistory, 0),account2dataId:make(map[string]int,0)}

	return pRet
}

// 加载数据
func (loginHistorySet *LoginHistorySet) LoadData(pDb *db.DbOperation) {
	var loginHistory LoginHistory
	selectstr := ""
	rows := pDb.Find(selectstr)

	for rows.Next() {
		err := rows.Scan(&loginHistory.dataId,&loginHistory.Account,&loginHistory.LoginTime,&loginHistory.LogoutTime,&loginHistory.Login)
		if err != nil {
			panic("error in:scanning the table LoginHistorySet")
		}

		loginHistorySet.WriteLoginHistory2Cache(loginHistory)
	}
}

// 从数据库写入用户登录历史数据到cache
func (loginHistorySet *LoginHistorySet) WriteLoginHistory2Cache(loginHistory LoginHistory) {
	loginHistorySet.rwLock.Lock()
	size := len(loginHistorySet.allLoginHistory)
	
	fmt.Println("-----------size----------")
	fmt.Println(size)
	loginHistorySet.allLoginHistory = append(loginHistorySet.allLoginHistory, loginHistory)
	loginHistorySet.account2dataId[loginHistory.Account]=size
	
	func() {
		loginHistorySet.rwLock.Unlock()
	}()
}

////读取搜索应用匹配的app
//func (appSet *AppSet) GetAllApps(str string) ([]App, error) {
//	var partApp []App
	

//	for _, app := range appSet.allAgenda {
//		//todo:比较操作
//		//……
//		if bool {
//			partApp = append(partApp, app)
//		}
//	}
//	if len(partApp) == 0 {
//		return partApp, errors.New("not exist related app data")
//	}
//	return partApp, nil
//}



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
