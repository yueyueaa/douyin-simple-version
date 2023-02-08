package controller

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Token_md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

// 注册函数
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password") // 读取用户给定的账号密码

	flag := Query_account(username)

	if flag  { //用户已经存在的情况
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		userinfo := Insert_newuser(username, password)
		token := Token_md5(fmt.Sprintf("%d", userinfo.Uid))
		//注册完成
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Sucessful Register"},
			UserId:   int64(userinfo.Uid),
			Token:    token,
		})

		usersLoginInfo[token] = User{
			Id:            int64(userinfo.Uid),
			Name:          userinfo.Name,
			FollowCount:   int64(userinfo.FollowCount),
			FollowerCount: int64(userinfo.FollowerCount),
		}
	}
}

// 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if exist, userinfo := Query_login(username, password); exist == 1 {

		token := Token_md5(fmt.Sprintf("%d", userinfo.Uid))

		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Successful Login"},
			UserId:   int64(userinfo.Uid),
			Token:    token,
		})

		usersLoginInfo[token] = User{
			Id:            int64(userinfo.Uid),
			Name:          userinfo.Name,
			FollowCount:   int64(userinfo.FollowCount),
			FollowerCount: int64(userinfo.FollowerCount),
		}

	} else if exist == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Password is wrong"},
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "user not exist"},
		})
	}
}

// 添加新用户
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
