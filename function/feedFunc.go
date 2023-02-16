package function

import (
	"douyin-simple-version/public"
	"douyin-simple-version/service/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

// 查找feeds流
func Query_feeds(user public.User, c *gin.Context) (feeds []public.Video, status public.Response) {
	db, err := middleware.InitDB() // 初始化数据库

	if err != nil {
		fmt.Printf("[DB ERR] Query_feeds\t%v\n", err)
		return
	}

	// all videos (limit 30)
	var videos []middleware.Video_info
	db.Limit(30).Order("VID desc").Find(&videos)
	fmt.Println(videos)

	for _, video := range videos {
		tempVideo, status := video_info(video.VID, user.Id)
		if status.StatusCode != 0 {
			return nil, status
		}
		feeds = append(feeds, tempVideo)
	}

	return feeds, public.Response{StatusCode: 0}
}

func video_info(VID int64, UID int64) (video public.Video, status public.Response) {
	db, err := middleware.InitDB() // 初始化数据库
	if err != nil {
		fmt.Printf("[DB ERR] Query_feeds\t%v\n", err)
		return video, public.Response{StatusCode: 1, StatusMsg: "Server error"}
	}

	var tempVideo middleware.Video_info
	tempVideo.VID = VID
	db.Model(&middleware.Video_info{}).Find(&tempVideo)

	var author public.User

	video.Id = VID
	video.CommentCount = tempVideo.CommentCount
	video.FavoriteCount = tempVideo.FavoriteCount

	// Author
	var author_info middleware.User_info
	author_info.UID = tempVideo.AuthorID
	db.Model(&middleware.User_info{}).Take(&author_info)
	author.FollowCount = author_info.FollowCount
	author.FollowerCount = author_info.FollowerCount
	author.Name = author_info.Name

	// crr User info
	{
		// Follow
		var follow_info middleware.Follow
		follow_info.UID = UID
		err := db.Model(&middleware.Follow{}).Where("Follow_ID = ?", author.Id).Find(&follow_info)
		if err == nil {
			author.IsFollow = true
		} else {
			author.IsFollow = false
		}

		// Favorite
		var favorite_info middleware.Favorite
		db.Model(&middleware.Favorite{}).Where("VID = ?", video.Id).Where("UID = ?", UID).Find(&favorite_info)
		video.IsFavorite = (favorite_info.Flag == 1)
	}

	// Video url
	var url_info middleware.Video_url
	db.Model(&middleware.Video_url{}).Where("VID = ?", video.Id).Find(&url_info)
	video.CoverUrl = url_info.CoverUrl
	video.PlayUrl = url_info.PlayUrl

	return video, public.Response{StatusCode: 0}
}
