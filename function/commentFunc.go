package function

import (
	"douyin-simple-version/public"
	"douyin-simple-version/service/middleware"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Publish_comment(user public.User, c *gin.Context) (CommentList []public.Comment, status public.Response) {
	db, err := middleware.InitDB()

	if err != nil {
		log.Println(err)
		return nil, public.Response{StatusCode: 1, StatusMsg: "Server Error"}
	}

	comment := middleware.Comment{
		UID:        user.Id,
		Content:    c.Query("comment_text"),
		CreateDate: time.Now(),
	}
	comment.VID, _ = strconv.ParseInt(c.Query("video_id"), 10, 64)
	db.Create(&comment)

	var video_info middleware.Video_info
	video_info.VID = comment.VID
	db.Model(&middleware.Video_info{}).Find(&video_info)
	db.Model(&video_info).Update("comment_count", video_info.CommentCount+1)

	return nil, public.Response{StatusCode: 0}
}

func Delete_comment(user public.User, c *gin.Context) (CommentList []public.Comment, status public.Response) {
	db, err := middleware.InitDB()

	if err != nil {
		log.Println(err)
		return nil, public.Response{StatusCode: 1, StatusMsg: "Server Error"}
	}

	var comment middleware.Comment
	comment.CID, _ = strconv.ParseInt(c.Query("comment_id"), 10, 64)
	db.Delete(&comment)

	var video_info middleware.Video_info
	video_info.VID = comment.VID
	db.Model(&middleware.Video_info{}).Find(&video_info)
	db.Model(&video_info).Update("comment_count", video_info.CommentCount-1)

	return nil, public.Response{StatusCode: 0}
}

func Query_commentList(user public.User, c *gin.Context) (CommentList []public.Comment, status public.Response) {
	db, err := middleware.InitDB()

	if err != nil {
		log.Println(err)
		return CommentList, public.Response{StatusCode: 1, StatusMsg: "Server Error"}
	}

	VID, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)

	var comments []middleware.Comment
	db.Model(&middleware.Comment{}).Where("VID = ?", VID).Find(&comments)

	for _, comment := range comments {
		var tempComment public.Comment
		var tempUser public.User
		tempComment.Id = comment.CID
		tempComment.CreateDate = comment.CreateDate.String()
		tempComment.Content = comment.Content

		// UserInfo
		{
			var user_info middleware.User_info
			db.Model(&middleware.User_info{}).Where("UID = ?", comment.UID).Find(&user_info)
			tempUser.Id = comment.UID
			tempUser.FollowCount = user_info.FollowCount
			tempUser.FollowerCount = user_info.FollowerCount
			tempUser.Name = user_info.Name
		}

		// Follow Info
		{
			var follow_info middleware.Follow
			db.Model(&middleware.Follow{}).Where("UID = ?", user.Id).Where("follow_id = ?", tempUser.Id).Find(&follow_info)
			tempUser.IsFollow = (follow_info.Flag == 1)
		}

		tempComment.User = tempUser
		CommentList = append(CommentList, tempComment)

	}

	return CommentList, public.Response{StatusCode: 0}
}
