package controller

import (
	"douyin-simple-version/function"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Author: godhands
// 互动接口的点赞与取消模块的response类型
type tiktokResponse struct {
	// 状态码, 0代表成功, 其他代表失败
	StatusCode int32 `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg,omitempty"`
}

// Author: godhands
// 对应点赞与取消赞操作
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	strvideo_id := c.Query("video_id")
	video_id, _ := strconv.ParseInt(strvideo_id, 10, 64)
	straction_type := c.Query("action_type")
	action_type, _ := strconv.ParseInt(straction_type, 10, 32)

	user, exist := usersLoginInfo[token]
	if exist {
		like := new(function.FavoriteFunction)
		err := like.UpdateLikes(user.Id, video_id, int32(action_type))
		if err != nil {
			c.JSON(http.StatusOK, tiktokResponse{
				StatusCode: 2,
				StatusMsg:  "Action Failed!",
			})
		} else {
			c.JSON(http.StatusOK, tiktokResponse{
				StatusCode: 0,
				StatusMsg:  "Action success!",
			})
		}
	} else {
		c.JSON(http.StatusOK, tiktokResponse{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
	}
}

// Author: godhands
// 获取点赞列表
func FavoriteList(c *gin.Context) {
	strUseId := c.Query("user_id")
	user_id, _ := strconv.ParseInt(strUseId, 10, 64)
	like := new(function.FavoriteFunction)
	var VideoList []Video
	video_idList, err := like.FindVideoIdFromUserId(user_id)

	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 2,
			},
			VideoList: VideoList,
		})
		return
	}

	for _, id := range video_idList {
		tmpVideo, err := like.VideoIdToVideo(id)
		if err != nil {
			continue
		}

		user := User{
			Id:            tmpVideo.Author.Id,
			Name:          tmpVideo.Author.Name,
			FollowCount:   tmpVideo.Author.FollowCount,
			FollowerCount: tmpVideo.Author.FollowerCount,
			IsFollow:      tmpVideo.Author.IsFollow,
		}

		result := Video{
			Id:            tmpVideo.Id,
			Author:        user,
			PlayUrl:       tmpVideo.PlayUrl,
			CoverUrl:      tmpVideo.CoverUrl,
			FavoriteCount: tmpVideo.FavoriteCount,
			CommentCount:  tmpVideo.CommentCount,
			IsFavorite:    tmpVideo.IsFavorite,
		}

		VideoList = append(VideoList, result)
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: VideoList,
	})
}
