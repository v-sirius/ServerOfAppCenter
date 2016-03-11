package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/Go-SQL-Driver/MySQL"
)

// 定义数据库及各类操作
type DbOperation struct {
	Db      *sql.DB // 数据库
	DataSrc string  // 数据源
}

var G_db *DbOperation

// 初始化数据库
func InitDbOperation(dsr string) *DbOperation {
	pRet := new(DbOperation)
	pRet.DataSrc = dsr
	return pRet
}

// 打开数据库
func (db *DbOperation) Open() {
	database, err := sql.Open("mysql", db.DataSrc)
	db.Db = database
	if err != nil {
		panic("cannot open database!")
	}
}

// 关闭数据库
func (db *DbOperation) Close() {
	if db.Db != nil {
		if db.Db.Close() != nil {
			panic("cannot close database!")
		}
	}
}

// 创建数据库表格
func (db *DbOperation) CreateTable() {
	// 创建用户表
	userTableSqlStr := "create table user(userId int(20),userAccount varchar(20),userPassword varchar(50),userCreattime varchar(50)),userLoginTime varcaar(50),userLogoutTime varchar(50),userLoginFlag int;"
	smt, err := db.Db.Prepare(userTableSqlStr)
	checkErr(err)
	smt.Exec()

	// 创建应用信息表
	appTableSqlStr := "create table app(appId int(20),appName varchar(20),appLogoUrl varchar(20)," +
		"appClass varchar(20),appInfo varchar(50),appImgUrl varchar(20),appSize float,appUrl varchar(20)," +
		"appScore float,appVersion varchar(20),appProvider varchar(20),appUpTime varchar(20),appDownloadTime int," +
		"appOnlineFlg int" + ");"
	smt, err = db.Db.Prepare(appTableSqlStr)
	checkErr(err)
	smt.Exec()

	// 创建用户登录记录表
//	userLoginTableSqlStr := "create table loginhistory(dataId int(20),userId int(20),loginWay varchar(20),longinTime varchar(20),logoutTime varchar(20),login int);"
//	smt, err = db.Db.Prepare(userLoginTableSqlStr)
//	checkErr(err)
//	smt.Exec()

	// 创建应用评价表
	appCommentTableSqlStr := "create table appComment(commentId int(20),userAccount char(20),appId int(20),appScore float,appComment varchar(50),commentTime varchar(20),appCommentDelFlag int);"
	smt, err = db.Db.Prepare(appCommentTableSqlStr)
	checkErr(err)
	smt.Exec()

}

//插入到数据库表中
func (db *DbOperation) Insert2Table(str string) bool {

	stmt, err := db.Db.Prepare(str)
	checkErr(err)
	res, err := stmt.Exec()
	defer stmt.Close()

	checkErr(err)

	//可以获得插入的id
	id, err2id := res.LastInsertId()
	if err2id == nil {
		fmt.Println("插入后的id:", id)
	}

	//可以获得影响的行数
	i, err := res.RowsAffected()
	if i > 0 && err == nil {
		return true
	} else {
		return false
	}

}

//从数据库表中删除数据
func (db *DbOperation) DelFromTable(str string) bool {

	stmt, err := db.Db.Prepare(str)
	checkErr(err)
	res, err := stmt.Exec()
	defer stmt.Close()

	checkErr(err)
	i, err := res.RowsAffected()
	if i > 0 && err == nil {
		return true
	} else {
		return false
	}

}

//查找数据
func (db *DbOperation) Find(str string) *sql.Rows {
	fmt.Println(str)

	rows, err := db.Db.Query(str)
	fmt.Println("--------------------err", err)
	if err != nil {
		panic("error in: selecting in table")
	}

	return rows
}

//更新数据
func (db *DbOperation) Update(str string) bool {

	stmt, err := db.Db.Prepare(str)
	checkErr(err)
	res, err := stmt.Exec()
	defer stmt.Close()

	checkErr(err)
	i, err := res.RowsAffected()
	if i > 0 && err == nil {
		return true
	} else {
		return false
	}

}

//错误检查
func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
