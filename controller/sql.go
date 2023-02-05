package controller

import (
	"douyin-simple-version/service/middleware"
	"fmt"
	"sync/atomic"
)

// 验证登录用户的账户密码是否正确
func Query_login(username string, password string) (flag bool, userid int64) {
	db, err := middleware.InitDB() // 初始化数据库

	if err != nil {
		fmt.Printf("[DB ERR] Query_login\t%v\n", err)
		return
	}
	defer db.Close()

	var name_in_database, password_in_database string

	sqlStr := "select name, password, ID from user where name= ?" //构造查询的sql语句

	err = db.QueryRow(sqlStr, username).Scan(&name_in_database, &password_in_database, &userid)

	if err != nil {
		flag = false
	} else {
		flag = true
	}

	if password == password_in_database {
		return flag, userid
	} else {
		return flag, userid
	}
}

// 查找对应的账号是否存在
func Query_account(str string) (err error) {
	db, err := middleware.InitDB() // 初始化数据库

	if err != nil {
		fmt.Printf("[DB ERR] Query_account\t%v\n", err)
		return
	}
	defer db.Close()

	sqlStr := "select name from user where name=?" //构造查询的sql语句

	var tem string

	err = db.QueryRow(sqlStr, str).Scan(&tem)

	return err
}

// 插入新用户
func Insert_newuser(username string, password string, userIdSequence int64) (err error) {
	db, err := middleware.InitDB() // 初始化数据库

	if err != nil {
		fmt.Printf("[DB ERR] Insert_newuser\t%v\n", err)
		return
	}
	defer db.Close()

	atomic.AddInt64(&userIdSequence, 1) //用户ID安全的自增1
	sqlStr := "INSERT INTO user(ID, name, follow_num, fans_num, password, sex, token, other) VALUES (?, ?, 0, 0, ?, 'male', ?, '');"
	_, err = db.Exec(sqlStr, userIdSequence, username, password, username+password)
	if err != nil { //插入失败
		return err
	}
	return nil
}

// 查找token是否存在
func Query_token(str string) (user User, err error) {
	db, err := middleware.InitDB() // 初始化数据库

	if err != nil {
		fmt.Printf("[DB ERR] Query_token\t%v\n", err)
		return
	}
	defer db.Close()

	sqlStr := "select ID, name, follow_num, fans_num, token from user where token=?" //构造查询的sql语句

	var tem string

	err = db.QueryRow(sqlStr, str).Scan(&user.Id, &user.Name, &user.FollowCount, &user.FollowerCount, &tem)

	return user, err
}

// 查找feeds流
func Query_feeds(feeds *[]Video) (err error) {
	db, err := middleware.InitDB() // 初始化数据库

	if err != nil {
		fmt.Printf("[DB ERR] Query_feeds\t%v\n", err)
		return
	}
	defer db.Close()

	// todo

	return err
}
