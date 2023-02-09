package controller

import (
	"douyin-simple-version/service/middleware"
	"fmt"
)

// 验证登录用户的账户密码是否正确
func Query_login(username string, password string) (status int64, userinfo middleware.User_info) {
	//status -1 用户不存在   0 密码错误  1 成功
	db, err := middleware.InitDB() // 初始化数据库

	if err != nil {
		fmt.Printf("[DB ERR] Query_login\t%v\n", err)
		return
	}
	var (
		user middleware.User
	)
	db.Select([]string{"UID", "name", "FollowCount", "FollowerCount"}).Where("name = ?", username).Take(&userinfo)
	if userinfo.Uid == 0 {
		status = -1
		return status, userinfo
	}
	db.Select("password").Where("UID = ?", userinfo.Uid).Take(&user)
	if user.Password != password {
		return 0, userinfo
	} else {
		return 1, userinfo
	}
}

// 根据username查找对应的账号是否存在
func Query_account(str string) (flag bool) {
	db, err := middleware.InitDB() // 初始化数据库

	if err != nil {
		fmt.Printf("[DB ERR] Query_account\t%v\n", err)
		return true
	}

	sqlStr := "select name from user_info where name=?" //构造查询的sql语句

	var userinfo middleware.User_info

	db.Raw(sqlStr, str).Scan(&userinfo)

	if userinfo.Name == str {
		return true
	} else {
		return false
	}
}

// 插入新用户
func Insert_newuser(username string, password string) (userinfo middleware.User_info) {
	db, err := middleware.InitDB() // 初始化数据库

	if err != nil {
		fmt.Printf("[DB ERR] Insert_newuser\t%v\n", err)
		return
	}
	user := middleware.User{
		Password: password,
	}
	userinfo = middleware.User_info{
		Name:          username,
		FollowCount:   0,
		FollowerCount: 0,
	}

	db.Create(&userinfo)
	user.Uid = userinfo.Uid//获取自增主键
	db.Create(&user)

	return userinfo
}

// 查找feeds流
func Query_feeds(feeds *[]Video) (err error) {
	// db, err := middleware.InitDB() // 初始化数据库

	// if err != nil {
	// 	fmt.Printf("[DB ERR] Query_feeds\t%v\n", err)
	// 	return
	// }
	// // defer db.Close()

	// var maxID int64
	// err = db.QueryRow("SELECT MAX(ID) FROM video_info").Scan(&maxID)

	// if err != nil {
	// 	fmt.Printf("[DB ERR] Query_feeds\t%v\n", err)
	// 	return
	// }
	// // defer db.Close()

	// var tempVideo Video

	// var i int64
	// for i = 0; (i < 30) && (maxID-i > 0); i++ {
	// 	tempVideo.Id = maxID - i
	// 	tempVideo.IsFavorite = false

	// 	// 查询视频基本信息
	// 	err = db.QueryRow("SELECT author,like_num,comment_num FROM video_info WHERE ID = ?", tempVideo.Id).Scan(&tempVideo.Author.Name, &tempVideo.FavoriteCount, &tempVideo.CommentCount)
	// 	if err != nil {
	// 		fmt.Printf("[DB ERR] Query_feeds_info\t%v\n", err)
	// 		return
	// 	}
	// 	// defer db.Close()

	// 	// 查询视频url信息
	// 	err = db.QueryRow("SELECT play_url,cover_url FROM video_url WHERE ID = ?", tempVideo.Id).Scan(&tempVideo.PlayUrl, &tempVideo.CoverUrl)
	// 	if err != nil {
	// 		fmt.Printf("[DB ERR] Query_feeds_url\t%v\n", err)
	// 		return
	// 	}
	// 	// defer db.Close()

	// 	*feeds = append(*feeds, tempVideo)
	// }

	return err
}
