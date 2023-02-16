package controller

import (
	"douyin-simple-version/function"
	"douyin-simple-version/public"
	"douyin-simple-version/service/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	public.Response
	CommentList []public.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	public.Response
	Comment middleware.Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	actionType := c.Query("action_type")
	if user, exist := usersLoginInfo[c.Query("token")]; exist {
		if actionType == "1" {
			_, status := function.Publish_comment(user, c)
			if status.StatusCode == 1 {
				c.JSON(http.StatusOK, status)
			} else {
				CommentList(c)
			}
		} else if actionType == "2" {
			_, status := function.Delete_comment(user, c)
			if status.StatusCode == 1 {
				c.JSON(http.StatusOK, status)
			} else {
				CommentList(c)
			}
		} else {
			c.JSON(http.StatusOK, public.Response{
				StatusCode: 1,
				StatusMsg:  "No operation",
			})
		}
	} else {
		c.JSON(http.StatusOK, public.Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't login",
		})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	if user, exits := usersLoginInfo[c.Query("token")]; exits {
		comments, status := function.Query_commentList(user, c)
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    status,
			CommentList: comments,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: public.Response{StatusCode: 1, StatusMsg: "User doesn't login"},
		})
	}
}
