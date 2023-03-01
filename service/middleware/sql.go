package middleware

import (
	"context"
	"database/sql"
	"douyin-simple-version/config"
	"fmt"
	"net"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"

	mysql_gorm "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ViaSSHDialer struct {
	client *ssh.Client
}

func (sself *ViaSSHDialer) Dial(context context.Context, addr string) (net.Conn, error) {
	return sself.client.Dial("tcp", addr)
}

// 初始化数据库
func InitDB() (db_gorm *gorm.DB, err error) {
	var db *sql.DB

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

	dsn := "test:123456@mysql+tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True"

	db, err = sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil { //ping通说明链接成功
		return nil, err
	}

	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	// 	logger.Config{
	// 		SlowThreshold: time.Second, // 慢 SQL 阈值
	// 		LogLevel:      logger.Info, // Log level
	// 		Colorful:      true,        // 禁用彩色打印
	// 	},
	// )

	gormDB, err := gorm.Open(mysql_gorm.New(mysql_gorm.Config{Conn: db}), &gorm.Config{
		// Logger: newLogger,
	})
	//去除62-69行，以及72行注释，可以进行sql语句可视化调试(即在控制台输出每一句执行的sql语句)

	if err != nil {
		fmt.Printf("[GORM ERR] Connection error\t%v\n", err)
	}

	sqldb, _ := gormDB.DB()
	sqldb.SetMaxOpenConns(config.MySQLMaxOpenConns)
	sqldb.SetMaxIdleConns(config.MySQLMaxIdleConns)

	return gormDB, nil
}
