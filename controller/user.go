package controller

import (
	"database/sql"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	// DB
	db, err := sql.Open("sqlite3", "./user.db")
	checkErr(err)

	// 若不存在，则创建表单
	sql_table := `
	CREATE TABLE IF NOT EXISTS "userinfo"(
		"name" VARCHAR(64) NULL,
		"password" VARCHAR(64) NULL,
		"id" int NULL
	)`
	db.Exec(sql_table)

	addUser, err := db.Prepare("INSERT INTO userinfo(name, password, id) values(?,?,?)")
	checkErr(err)

	queryUser, err := db.Prepare("select * from userinfo where name = ?")
	checkErr(err)

	var a string
	var b int

	if queryUser.QueryRow(username).Scan(&a, &a, &b) == nil { //检测userLoginInfo，如果用户存在，返回错误
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else { //用户不存在，创建用户
		//添加用户
		atomic.AddInt64(&userIdSequence, 1)
		newUser := User{
			Id:   userIdSequence,
			Name: username,
		}
		usersLoginInfo[token] = newUser

		_, err := addUser.Exec(username, password, userIdSequence)
		checkErr(err)

		//返回状态码
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userIdSequence,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")

	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
