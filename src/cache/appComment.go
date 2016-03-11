package cache

import (
	"errors"
	"log"
	"strconv"
	"fmt"
	"db"
	"sync"
)

const APPCOMMENT_CHAN_NUM=1
// app评论表
type AppComment struct {
	CommentId         int    // 
	UserAccount string
	AppId     int //
	AppScore float64 //
	AppComment string // 
	CommentTime    string //
	DelFlag int //是否删除标识位 0：未删除 1：删除 
}

//type AppComments struct{
//	CommentIds []int
//	//UserId int
//}

// 所有app评论的表
type AppCommentSet struct {
	allAppComment []AppComment
	appId2AppComment map[int][]int
	appCommentId2AppComment map[int]int
	
	AppCommentChan chan AppComment
	rwLock sync.RWMutex // 读写锁
}

// 初始化app评论表缓存
func InitAppCommentSet() *AppCommentSet {
	pRet := new(AppCommentSet)
	pRet = &AppCommentSet{allAppComment: make([]AppComment, 0),appId2AppComment:make(map[int][]int,0),AppCommentChan: make(chan AppComment, APPCOMMENT_CHAN_NUM)}

	return pRet
}

//获得comment表的长度
func (cmtSet AppCommentSet) GetLength() int{
	cmtSet.rwLock.RLock()
	defer func(){
		cmtSet.rwLock.RUnlock()
	}()
	
	return len(cmtSet.allAppComment)
}

// 通过评价id读取评价数据
func (cmtSet *AppCommentSet) GetCommentById(id int) (AppComment, error) {
	cmtSet.rwLock.RLock()
	defer func() {
		cmtSet.rwLock.RUnlock()
	}()

	if idx, flg := cmtSet.appCommentId2AppComment[id]; flg {
		return cmtSet.allAppComment[idx], nil
	}

	return AppComment{CommentId: -1}, errors.New("not exist related user data")
}

// 加载app评论数据
func (appCommentSet *AppCommentSet) LoadData(pDb *db.DbOperation) {
	var appComment AppComment
	selectstr := "select commentId,userAccount,appId,appScore,appComment,commentTime,delFlag from appcomment"
	rows := pDb.Find(selectstr)

	for rows.Next() {
		err := rows.Scan(&appComment.CommentId,&appComment.UserAccount,&appComment.AppId,&appComment.AppScore,&appComment.AppComment,&appComment.CommentTime,&appComment.DelFlag)
		if err != nil {
			panic("error in:scanning the table appCommentSet")
		}

		appCommentSet.WriteApp2Cache(appComment)
	}
}

// 从数据库读取数据写入到cache
func (appCommentSet *AppCommentSet) WriteApp2Cache(appComment AppComment) {
	appCommentSet.rwLock.Lock()
	appCommentSet.allAppComment = append(appCommentSet.allAppComment, appComment)
	appCommentSet.appId2AppComment[appComment.AppId]=append(appCommentSet.appId2AppComment[appComment.AppId],appComment.CommentId)
	
	func() {
		appCommentSet.rwLock.Unlock()
	}()
}

// 从服务器将数据写入cache
func (appCommentSet *AppCommentSet) WriteAppCommentFromServer2Cache(app_cmt AppComment) {
	fmt.Println("--------write to cache--------")

	appCommentSet.rwLock.Lock()
	defer func() {
		appCommentSet.rwLock.Unlock()
	}()

	size := len(appCommentSet.allAppComment)
	
	fmt.Println(app_cmt)

	fmt.Println("-----------size----------")
	fmt.Println(size)
	
	appCommentSet.allAppComment = append(appCommentSet.allAppComment, app_cmt)
	appCommentSet.appId2AppComment[app_cmt.AppId] = append(appCommentSet.appId2AppComment[app_cmt.AppId],size)
	appCommentSet.appCommentId2AppComment[size]=size
}

//用户评价时 从服务器管道中将评价信息写入db中
func (appCommentSet *AppCommentSet) Write2DbAppComment() {
	fmt.Println("------------entering write2db appcomment-----------")

	for comment := range appCommentSet.AppCommentChan {

		fmt.Println("------------------comment -------------------")
		fmt.Println(comment)
		
		var sql = "insert into appComment(userAccount,appId,appScore,appComment,commentTime,delFlag) values ('"  + comment.UserAccount + "','"  + strconv.Itoa(comment.AppId) + "','"  + strconv.FormatFloat(comment.AppScore,'f',-1,32) +"','"+comment.AppComment+"','"+comment.CommentTime+"','"+strconv.Itoa(0)+"');"
		log.Println("sql:", sql)
		if !db.G_db.Insert2Table(sql) {
			fmt.Println("应用评价插入不成功！")
		} else {
			fmt.Println("entering weiteComment2db OK ! ! !")
		}
	}
}

//修改cache中的删除评论标志位
func (appCommentSet *AppCommentSet)ModifyCmtDelFlagInCache(AppCmtId int,flag string){
	f,_:=strconv.Atoi(flag)
	appCommentSet.allAppComment[AppCmtId].DelFlag=f
	
}

// 修改数据库中删除评论标志位
func (appCommentSet *AppCommentSet) ModifyCmtDelFlag2DbAppCmt(AppCmtId int,flag string)bool {
	var sql = "update appComment set appCommentDelFlag='"+flag+"'"+"where appCommentId='"+strconv.Itoa(AppCmtId)+"';"
	log.Println("sql:", sql)
	if !db.G_db.Update(sql) {
		fmt.Println("删除评论标志位更改不成功！")
		return false
		} else {
			fmt.Println("成功修改删除评论标志位 ! ! !")
			return true
		}
}