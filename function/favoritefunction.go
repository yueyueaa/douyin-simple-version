// Author: godhands
package function

import (
	"douyin-simple-version/service/middleware"
	"log"
)

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type FavortieInterface interface {
	// 用于更新LIKE表，进行点赞和取消赞的操作
	UpdateLikes(UserId, VideoId int64, ActionType int32) (err error)

	// 通过视频ID获取视频信息
	VideoIdToVideo(video_id int64) (video Video, err error)

	// 通过我们的用户ID找到所有点赞的视频ID
	FindVideoIdFromUserId(user_id int64) (video_idList []int64, err error)

	// 判断当前的登录用户对于当前视频的点赞状态
	IsFavorite(user_id, video_id int64) (okk bool, err error)

	// 获取当前视频的点赞数量
	FavoriteCount(video_id int64) (count int64, err error)
}

type FavoriteFunction struct {
	FavortieInterface
}

// 用于更新LIKE表，进行点赞和取消赞的操作
func (t *FavoriteFunction) UpdateLikes(UserId, VideoId int64, ActionType int32) (err error) {
	db, err := middleware.InitDB() // 初始化数据库
	if err != nil {
		log.Println(err)
		return err
	}

	like := middleware.Like{}
	db.Model(&middleware.Like{}).Where("UID = ? AND VID = ?", uint(UserId), uint(VideoId)).First(&like)

	if like == (middleware.Like{}) {
		like = middleware.Like{
			UID:  uint(UserId),
			VID:  uint(VideoId),
			Flag: ActionType,
		}
		result := db.Create(&like)
		return result.Error
	} else {
		like = middleware.Like{}
		db.Model(&middleware.Like{}).Where("UID = ? AND VID = ?", uint(UserId), uint(VideoId)).Take(&like)
		like.Flag = ActionType
		db.Save(like)
	}

	return nil
}

// 通过视频ID获取视频信息
func (t *FavoriteFunction) VideoIdToVideo(video_id int64) (video Video, err error) {
	db, err := middleware.InitDB()

	if err != nil {
		return Video{}, err
	}

	video_info := middleware.Video_info{}
	db.Model(&middleware.Video_info{}).Where("VID = ?", uint(video_id)).Take(&video_info)

	user_info := middleware.User_info{}
	db.Model(&middleware.User_info{}).Where("UID = ?", video_info.Author).Take(&user_info)

	video_url := middleware.Video_url{}
	db.Model(&middleware.Video_url{}).Where("VID = ?", uint(video_id)).Take(&video_url)

	user := User{
		Id:            int64(user_info.Uid),
		Name:          user_info.Name,
		FollowCount:   int64(user_info.FollowCount),
		FollowerCount: int64(user_info.FollowerCount),
		IsFollow:      true,
	}

	favorite_count, err := t.FavoriteCount(video_id)
	if err != nil {
		log.Println(err)
		return Video{}, err
	}

	is_favorite, err := t.IsFavorite(user.Id, video_id)
	if err != nil {
		log.Println(err)
		return Video{}, err
	}

	video = Video{
		Id:            int64(video_info.VID),
		Author:        user,
		PlayUrl:       video_url.PlayUrl,
		CoverUrl:      video_url.CoverUrl,
		FavoriteCount: favorite_count,
		CommentCount:  int64(0),
		IsFavorite:    is_favorite,
	}

	return video, nil
}

// 通过我们的用户ID找到所有点赞的视频ID
func (t *FavoriteFunction) FindVideoIdFromUserId(user_id int64) (video_idList []int64, err error) {
	db, err := middleware.InitDB()

	if err != nil {
		log.Println(err)
		return video_idList, err
	}

	var LikeList []middleware.Like
	db.Model(&middleware.Like{}).Where("UID = ? AND Flag = ?", uint(user_id), 1).Take(&LikeList)

	for _, id := range LikeList {
		video_idList = append(video_idList, int64(id.VID))
	}
	return video_idList, nil
}

// 判断当前的登录用户对于当前视频的点赞状态
func (t *FavoriteFunction) IsFavorite(user_id, video_id int64) (okk bool, err error) {
	db, err := middleware.InitDB()

	if err != nil {
		log.Println(err)
		return false, err
	}

	like := middleware.Like{}
	db.Model(&middleware.Like{}).Where("UID = ? AND VID = ?", uint(user_id), uint(video_id)).First(&like)

	if like == (middleware.Like{}) {
		return false, nil
	} else {
		like = middleware.Like{}
		db.Model(&middleware.Like{}).Where("UID = ? AND VID = ?", uint(user_id), uint(video_id)).Take(&like)
		return like.Flag == int32(1), nil
	}
}

// 获取当前视频的点赞数量
func (t *FavoriteFunction) FavoriteCount(video_id int64) (count int64, err error) {
	db, err := middleware.InitDB()

	if err != nil {
		log.Println(err)
		return 0, err
	}

	db.Model(&middleware.Like{}).Where("VID = ?", uint(video_id)).Count(&count)
	return count, err
}
