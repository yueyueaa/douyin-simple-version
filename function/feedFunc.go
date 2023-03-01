package function

import (
	"douyin-simple-version/public"
	"douyin-simple-version/service/middleware"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var author_id int64

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

		var tmpvideo public.Video
		//已有数据先录入，避免二次查询
		{
			tmpvideo.Id = video.VID
			tmpvideo.Title = video.Title
			tmpvideo.FavoriteCount = video.FavoriteCount
			tmpvideo.CommentCount = video.CommentCount
		}
		author_id = video.AuthorID
		fmt.Println("*****author_id = :", author_id)
		tempvideo2, flag := video_info(video.VID, user.Id)

		if !flag {
			return nil, public.Response{StatusCode: 1, StatusMsg: "sever error"}
		} else {
			tmpvideo.Author = tempvideo2.Author
			tmpvideo.PlayUrl = tempvideo2.PlayUrl
			tmpvideo.CoverUrl = tempvideo2.CoverUrl
			tmpvideo.IsFavorite = tempvideo2.IsFavorite
		}
		fmt.Println("show video", tmpvideo)
		feeds = append(feeds, tmpvideo)
	}

	return feeds, public.Response{StatusCode: 0, StatusMsg: "ok"}
}

func video_info(VID int64, UID int64) (video public.Video, flag bool) {
	db, err := middleware.InitDB() // 初始化数据库
	if err != nil {
		fmt.Printf("[DB ERR] Query_feeds\t%v\n", err)
		return video, false
	}

	var url_info middleware.Video_url
	var author public.User
	var author_info middleware.User_info

	video.Id = VID

	// Author info
	fmt.Println("函数内测试author：",author_id)
	author_info.UID = 1
	db.Take(&author_info)
	fmt.Println("测试搜索到的作者信息",author_info)
	author.Id = author_info.UID
	author.FollowCount = author_info.FollowCount
	author.FollowerCount = author_info.FollowerCount
	author.Name = author_info.Name

	// crr User info
	{
		// Follow
		var tmpfollow int
		var follow_info middleware.Follow
		follow_info.UID = UID
		tmp2err := db.Where("Follow_ID = ?", author_id).Find(&follow_info).Error
		if errors.Is(tmp2err, gorm.ErrRecordNotFound) {
			fmt.Println("查询不到数据，未关注")
			tmpfollow = 0
		} else if err != nil {
			fmt.Println("查询失败", err)
		} else {
			tmpfollow = 1
			fmt.Println("查找到目标数据，已关注")
		}
		author.IsFollow = (tmpfollow == 1)
		video.Author = author

		// Favorite
		var favorite int
		var favorite_info middleware.Favorite
		tmperr := db.Where("VID = ?", video.Id).Where("UID = ?", UID).First(&favorite_info).Error

		if errors.Is(tmperr, gorm.ErrRecordNotFound) {
			fmt.Println("查询不到数据，未点赞")
			favorite = 0
		} else if err != nil {
			// 如果err不等于record not found错误，又不等于nil，那说明sql执行失败了。
			fmt.Println("查询失败", err)
		} else {
			favorite = 1
			fmt.Println("查找到目标数据，已点赞")
		}

		video.IsFavorite = (favorite == 1)
	}

	// Video url
	db.Model(&middleware.Video_url{}).Where("VID = ?", VID).Take(&url_info)
	video.CoverUrl = url_info.CoverUrl
	video.PlayUrl = url_info.PlayUrl
	//暂时不考虑一个视频对应多个URL的情况
	return video, true
}
