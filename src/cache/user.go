package cache

import (
	"db"
	"errors"
	"fmt"
	"log"
	_ "github.com/Go-SQL-Driver/MySQL"
	_ "image/jpeg"
	"sync"
)

const USER_CHAN_NUM = 1

// 定义用户
type User struct {
	Id   int    // 唯一标识
	Account string // 用户账号
	Password string //密码
	CreateTime  string //time.Time // 创建时间
	LoginTime string //
	LogoutTime string // 
	LoginFlag    string // 登录标志位  0:未登录 1：已登录
}

// 集合，只有id，用户代码和手机号是唯一的
type UserSet struct {
	allUser       []User         // 所有的用户
	id2User       map[int]int    // id到用户的映射
	account2User     map[string]int // 用户账号到用户的映射
	account2LoginFlag map[string]string //用户到登录标志位的映射
	
	UserChan chan User
	rwLock sync.RWMutex // 读写锁
}

// 初始化函数
func InitUserSet() *UserSet {
	pRet := new(UserSet)

	pRet = &UserSet{allUser: make([]User, 0), id2User: make(map[int]int, 0),account2User : make(map[string]int, 0),account2LoginFlag: make(map[string]string, 0),UserChan: make(chan User, USER_CHAN_NUM)}

	return pRet
}

//获得user表的长度
func (usrSet UserSet) GetLength() int{
	usrSet.rwLock.RLock()
	defer func(){
		usrSet.rwLock.RUnlock()
	}()
	
	return len(usrSet.allUser)
}

// 检查用户是否存在－－根据唯一的账号来判断
func (usr UserSet) IsExist(account string) (User, bool) {
	if i, flag := usr.account2User["account"]; flag {
		return usr.allUser[i], true
	}
	return User{Id: -1}, false
}

// 通过用户id读取用户数据
func (pSet *UserSet) GetUserById(id int) (User, error) {
	pSet.rwLock.RLock()
	defer func() {
		pSet.rwLock.RUnlock()
	}()

	if idx, flg := pSet.id2User[id]; flg {
		return pSet.allUser[idx], nil
	}

	return User{Id: -1}, errors.New("not exist related user data")
}

// 根据account获取用户
func (pSet *UserSet) GetUserByAccount(account string) (User, error) {
	pSet.rwLock.RLock()
	defer func() {
		pSet.rwLock.RUnlock()
	}()

	if idx, flg := pSet.account2User[account]; flg {
		return pSet.allUser[idx], nil
	}

	return User{Id: -1}, errors.New("not exist related user data")
}

// 加载用户数据
func (pUsersSet *UserSet) LoadData(pDb *db.DbOperation) {
	var user User
	fmt.Println("begining:load userset........")
	selectstr := "select userId,userAccount,userPassword,userCreatTime,userLoginTime,userLogoutTime,userLoginFlag from user"//此处应按照数据库字段顺序写查找内容
	rows := pDb.Find(selectstr)

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Account, &user.Password, &user.CreateTime,&user.LoginTime,&user.LogoutTime,&user.LoginFlag)

		if err != nil {
			panic("error in:scanning the table user")
		}
		

		fmt.Println("-------------------Load data to cache--------------------")
		fmt.Println(user)
		pUsersSet.WriteUserFromDb2Cache(user)
	}
}

// 从db写入数据到cache
func (pSet *UserSet) WriteUserFromDb2Cache(usr User) {
	fmt.Println("--------write to cache--------")
	fmt.Println(usr)

	pSet.rwLock.Lock()
	defer func() {
		pSet.rwLock.Unlock()
	}()

	size := len(pSet.allUser)
	//usr.Id = size
	fmt.Println("-----------size----------")
	fmt.Println(size)
	pSet.allUser = append(pSet.allUser, usr)

	pSet.id2User[usr.Id] = size 
	pSet.account2User[usr.Account] = size 
}

// 从服务器将数据写入cache
func (pSet *UserSet) WriteUserFromServer2Cache(usr User) {
	fmt.Println("--------write to cache--------")

	pSet.rwLock.Lock()
	defer func() {
		pSet.rwLock.Unlock()
	}()

	size := len(pSet.allUser)
	
	fmt.Println(usr)

	fmt.Println("-----------size----------")
	fmt.Println(size)
	pSet.allUser = append(pSet.allUser, usr)

	pSet.id2User[size] = size
	pSet.account2User[usr.Account] = size
}



//注册时 从服务器管道中将user信息写入db
func (pUsersSet *UserSet) Write2DbUser() {
	fmt.Println("------------entering write2db-----------")

	for user := range pUsersSet.UserChan {

		fmt.Println("------------------user in db-------------------")
		fmt.Println(user)
		
		var sql = "insert into user(userAccount,userPassword,registerTime,userLoginFlag) values ("  + user.Account + "','"  + user.Password + "','"  + user.CreateTime +"','"+user.LoginFlag+"');"
		log.Println("sql:", sql)
		if !db.G_db.Insert2Table(sql) {
			fmt.Println("user表插入不成功！")
		} else {
			fmt.Println("entering weite2db OK ! ! !")
		}
	}
}

//登录时 修改登录标志位到cache
func (pUsersSet *UserSet)ModifyLoginFlagInCache(u User,flag string){
	pUsersSet.account2LoginFlag[u.Account]=flag
	uId:=pUsersSet.account2User[u.Account]
	pUsersSet.allUser[uId].LoginFlag=flag
}

//登录时 修改登录标志位到db
func (pUsersSet *UserSet) ModifyLoginFlag2DbUser(u User,flag string)bool {
	var sql = "update user set userLoginFlag='"+flag+"'"+"where userAccount='"+u.Account+"';"
	log.Println("sql:", sql)
	if !db.G_db.Update(sql) {
		fmt.Println("user表更新不成功！")
		return false
		} else {
			fmt.Println("updating dbUser OK ! ! !")
			return true
		}
}	

//忘记密码时 修改新密码到cache
func (pUsersSet *UserSet)PasswordResetInCache(u User,password string){
	
	uId:=pUsersSet.account2User[u.Account]
	pUsersSet.allUser[uId].Password=password
}

//忘记密码时 修改新密码到db
func (pUsersSet *UserSet) PasswordReset2DbUser(u User,password string)bool {
	var sql = "update user set userPassword='"+password+"'"+"where userAccount='"+u.Account+"';"
	log.Println("sql:", sql)
	if !db.G_db.Update(sql) {
		fmt.Println("密码修改不成功！")
		return false
		} else {
			fmt.Println("密码修改成功 ! ! !")
			return true
		}
}
