package controller

import (
	"douyin-simple-version/service/middleware"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")

	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	token := c.Query("token")
	user, exist := usersLoginInfo[token]
	if !exist || user.Id != userId {
		return
	}
	var videos []Video
	db, err := middleware.InitDB()
	if err != nil {
		return
	}
	rows, err := db.Raw("SELECT VID,like_num,comment_num,title FROM video_info WHERE author = ?", user.Id).Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		var vid, like_num, comment_num int64
		var play_url, cover_url, title string
		var isFavorite int
		rows.Scan(&vid, &like_num, &comment_num, &title)
		//查找对应视频url
		err := db.Raw("SELECT play_url,cover_url FROM video_url WHERE VID = ?", vid).Row().Scan(&play_url, &cover_url)
		if err != nil {
			return
		}
		//查找点赞信息
		err = db.Raw("SELECT flag FROM likes WHERE VID = ? AND UID = ?", vid, user.Id).Row().Scan(&isFavorite)
		if err != nil {
			isFavorite=0
		}
		videos = append(videos, Video{
			vid,
			user,
			play_url,
			cover_url,
			like_num,
			comment_num,
			isFavorite != 0,
		})
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
