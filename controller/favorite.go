package controller

import (
	"douyin-simple-version/function"
	"douyin-simple-version/public"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FavoriteActionResponse struct {
	public.Response
}

type FavoriteListResponse struct {
	public.Response
	Favoritelist []public.Video `json:"video_list,omitempty"`
}

// Author: godhands
// 对应点赞与取消赞操作
func FavoriteAction(c *gin.Context) {
	if user, exist := usersLoginInfo[c.Query("token")]; exist {
		_, status := function.Set_favorite(user, c)
		c.JSON(http.StatusOK, FavoriteActionResponse{
			Response: status,
		})
	} else {
		c.JSON(http.StatusOK, public.Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't login",
		})
	}
}

// Author: godhands
// 获取点赞列表
func FavoriteList(c *gin.Context) {
	if user, exist := usersLoginInfo[c.Query("token")]; exist {
		favorites, status := function.Favorite_list(user, c)
		c.JSON(http.StatusOK, FavoriteListResponse{
			Response:     status,
			Favoritelist: favorites,
		})
	} else {
		c.JSON(http.StatusOK, public.Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't login",
		})
	}
}
