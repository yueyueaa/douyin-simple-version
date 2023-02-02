package controller

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"io"
	"net"
	"os"
	"sync/atomic"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
)

type ViaSSHDialer struct {
	client *ssh.Client
}

func (sself *ViaSSHDialer) Dial(context context.Context, addr string) (net.Conn, error) {
	return sself.client.Dial("tcp", addr)
}

func GetKey() (account string, password string) {
	filePath := "./key.txt"
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("get key fail", err)
	}
	//关闭file句柄
	defer file.Close()

	reader := bufio.NewReader(file)
	i := 1
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			password = str
			break
		}
		if i == 1 {
			account = str
		}
		i++
	}
	account = account[:len(account)-2]
	return account, password
}

func InitSsh() {
	account, password := GetKey()
	// 一个ClientConfig指针,指向的对象需要包含ssh登录的信息
	config := &ssh.ClientConfig{
		User: account, //用户名
		Auth: []ssh.AuthMethod{
			ssh.Password(password), //密码
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //服务端验证
	}
	client, err := ssh.Dial("tcp", "60.255.139.184:20022", config)

	if err != nil {
		panic("connection errror") //抛出异常
	}

	mysql.RegisterDialContext("mysql+tcp", (&ViaSSHDialer{client}).Dial)
}

func InitDB() (db *sql.DB, err error) {

	InitSsh()

	dsn := "root@mysql+tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True"

	db, err = sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil { //ping通说明链接成功
		return nil, err
	}

	SetDB(db)

	return db, nil
}
func SetDB(db *sql.DB) {
	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.SetMaxOpenConns(20000)
	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	db.SetMaxIdleConns(0)
}

func QuerytoLogin(username string, password string) (flag bool, userid int64) {
	db, err := InitDB() // 初始化数据库

	if err != nil {
		fmt.Printf("QuerytoLogin DB ERROR ----， %v", err)
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

func Query_account(str string) (err error) { // 查找对应的账号是否存在
	db, err := InitDB() // 初始化数据库

	if err != nil {
		fmt.Printf("Query_account DB ERROR ----， %v", err)
		return
	}
	defer db.Close()

	sqlStr := "select name from user where name=?" //构造查询的sql语句

	var tem string

	err = db.QueryRow(sqlStr, str).Scan(&tem)

	return err
}

func Insert(username string, password string, userIdSequence int64) (err error) {

	db, err := InitDB() // 初始化数据库

	if err != nil {
		fmt.Printf("Insert DB ERROR ----， %v", err)
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
func Query_token(str string) (user User, err error) { // 查找token是否存在
	db, err := InitDB() // 初始化数据库

	if err != nil {
		fmt.Printf("Query_token DB ERROR ----， %v", err)
		return
	}
	defer db.Close()

	sqlStr := "select ID, name, follow_num, fans_num, token from user where token=?" //构造查询的sql语句

	var tem string

	err = db.QueryRow(sqlStr, str).Scan(&user.Id, &user.Name, &user.FollowCount, &user.FollowerCount, &tem)

	return user, err
}
