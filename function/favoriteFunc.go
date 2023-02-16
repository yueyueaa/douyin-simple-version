// Author: godhands
package function

import (
	"douyin-simple-version/public"
	"douyin-simple-version/service/middleware"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Set_favorite(user public.User, c *gin.Context) (favoriteList []public.Video, status public.Response) {
	action_type, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if action_type != 1 && action_type != 2 {
		return nil, public.Response{StatusCode: 1, StatusMsg: "No operation"}
	}

	db, err := middleware.InitDB()
	if err != nil {
		log.Println(err)
		return nil, public.Response{StatusCode: 1, StatusMsg: "Server Error"}
	}

	var favorite middleware.Favorite
	favorite.UID = user.Id
	favorite.VID, _ = strconv.ParseInt(c.Query("video_id"), 10, 64)

	var video_info middleware.Video_info
	video_info.VID = favorite.VID
	db.Model(&middleware.Video_info{}).Find(&video_info)

	if action_type == 1 {
		favorite.Flag = 1
		db.Model(&middleware.Favorite{}).Create(&favorite)
		db.Model(&video_info).Update("favorite_count", video_info.FavoriteCount+1)
	} else {
		db.Model(&middleware.Favorite{}).
			Where("UID = ?", favorite.UID).Where("VID = ?", favorite.VID).
			Delete(&favorite)
		db.Model(&video_info).Update("favorite_count", video_info.FavoriteCount-1)
	}

	return nil, public.Response{StatusCode: 0}
}

func Favorite_list(user public.User, c *gin.Context) (favoriteList []public.Video, status public.Response) {
	db, err := middleware.InitDB()
	if err != nil {
		log.Println(err)
		return nil, public.Response{StatusCode: 1, StatusMsg: "Server Error"}
	}

	var favorites []middleware.Favorite
	db.Model(&middleware.Favorite{}).Where("UID = ?", user.Id).Find(&favorites)

	for _, favorite := range favorites {
		tempFavorite, status := video_info(favorite.VID, favorite.VID)
		if status.StatusCode != 0 {
			return nil, status
		}
		favoriteList = append(favoriteList, tempFavorite)
	}

	return favoriteList, public.Response{StatusCode: 0}
}
