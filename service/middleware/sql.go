package middleware

import (
	"context"
	"database/sql"
	"douyin-simple-version/config"
	"net"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
)

type ViaSSHDialer struct {
	client *ssh.Client
}

func (sself *ViaSSHDialer) Dial(context context.Context, addr string) (net.Conn, error) {
	return sself.client.Dial("tcp", addr)
}

// 初始化数据库
func InitDB() (db *sql.DB, err error) {
	account, password := config.MySQLUser, config.MySQLPwd
	// 一个ClientConfig指针,指向的对象需要包含ssh登录的信息
	sshClientConfig := &ssh.ClientConfig{
		User: account, //用户名
		Auth: []ssh.AuthMethod{
			ssh.Password(password), //密码
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //服务端验证
	}
	client, err := ssh.Dial("tcp", config.MySQLIP, sshClientConfig)

	if err != nil {
		panic("connection errror") //抛出异常
	}

	mysql.RegisterDialContext("mysql+tcp", (&ViaSSHDialer{client}).Dial)

	dsn := "root@mysql+tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True"

	db, err = sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil { //ping通说明链接成功
		return nil, err
	}

	db.SetMaxOpenConns(config.MySQLMaxOpenConns)
	db.SetMaxIdleConns(config.MySQLMaxIdleConns)

	return db, nil
}
