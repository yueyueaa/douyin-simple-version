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
	user.Uid = userinfo.Uid //获取自增主键
	db.Create(&user)

	return userinfo
}

// 查找feeds流
func Query_feeds(token string) (feeds []Video) {
	db, err := middleware.InitDB() // 初始化数据库

	if err != nil {
		fmt.Printf("[DB ERR] Query_feeds\t%v\n", err)
		return
	}

	var video_info middleware.Video_info
	var author_info middleware.User_info
	var url_info middleware.Video_url
	var follow_info middleware.Follow
	var favorite_info middleware.Like

	db.Last(&video_info)
	maxID := video_info.VID

	var tempVideo Video
	var tempAuthor User

	for i := 0; (i < 30) && (int(maxID)-i > 0); i++ {

		tempVideo.Id = int64(maxID) - int64(i)

		video_info.VID = uint(tempVideo.Id)
		db.Find(&video_info)
		tempAuthor.Id = int64(video_info.Author)
		tempVideo.FavoriteCount = int64(video_info.LikeNum)
		tempVideo.CommentCount = int64(video_info.CommentNum)

		// Author
		{
			author_info.Uid = uint(tempAuthor.Id)
			db.Find(&author_info)
			tempAuthor.FollowCount = int64(author_info.FollowCount)
			tempAuthor.FollowerCount = int64(author_info.FollowerCount)
			tempAuthor.Name = author_info.Name
		}

		// crr User info
		{
			UID := usersLoginInfo[token].Id
			// Follow
			{
				follow_info.UID = uint(UID)
				err := db.Where("Follow_ID = ?", tempAuthor.Id).Find(&follow_info)

				if err != nil {
					tempAuthor.IsFollow = false
				} else {
					tempAuthor.IsFollow = true
				}
			}

			// Like
			{
				favorite_info.VID = uint(tempVideo.Id)
				err := db.Where("UID = ?", UID).Find(&favorite_info)
				if (err != nil) || (favorite_info.Flag == 0) {
					tempVideo.IsFavorite = false
				} else {
					tempVideo.IsFavorite = true
				}
			}
		}

		url_info.VID = uint(tempVideo.Id)
		db.Take(&url_info)
		tempVideo.CoverUrl = url_info.CoverUrl
		tempVideo.PlayUrl = url_info.PlayUrl

		fmt.Println(tempVideo)

		feeds = append(feeds, tempVideo)
	}

	return feeds
}
