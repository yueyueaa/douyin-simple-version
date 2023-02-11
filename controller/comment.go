package controller

import (
	"douyin-simple-version/service/middleware"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	strvideo_id := c.Query("video_id")
	actionType := c.Query("action_type")
	videoId, _ := strconv.ParseInt(strvideo_id, 10, 64)
	createDate := time.Now()
	if user, exist := usersLoginInfo[token]; exist {
		if actionType == "1" {
			commentText := c.Query("comment_text")
			err := Insert_comments(uint(user.Id), uint(videoId), commentText, createDate)
			if err != nil {
				c.JSON(http.StatusOK, Response{
					StatusCode: 1,
				})
			}
			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
				Comment: Comment{
					Id:         1,
					User:       user,
					Content:    commentText,
					CreateDate: createDate.String(),
				}})
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
		})
		if actionType == "2" {
			comment_id := c.Query("comment_id")
			commentId, _ := strconv.ParseInt(comment_id, 10, 64)
			err := Delete_comments(uint(user.Id), uint(videoId), uint(commentId))
			if err != nil {
				c.JSON(http.StatusOK, Response{
					StatusCode: 1,
				})
			}
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
		})

	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	token := c.Query("token")
	user, exist := usersLoginInfo[token]
	if !exist {
		return
	}
	db, err := middleware.InitDB()
	if err != nil {
		return
	}
	var comments []Comment
	rows, err := db.Raw("SELECT CID, UID, CommentText, CommentTime FROM comments WHERE VID = ?", videoId).Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		var user_id, comment_id int64
		var comment_text string
		var comment_time string
		rows.Scan(&comment_id, &user_id, &comment_text, &comment_time)

		comments = append(comments, Comment{
			comment_id,
			user,
			comment_text,
			comment_time,
		})
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{
			StatusCode: 0,
		},
		CommentList: comments,
	})
}
