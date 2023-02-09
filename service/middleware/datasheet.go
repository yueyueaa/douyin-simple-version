package middleware

import (
	"time"
)

/*
  Author : cuizhu12138
  date : 2/8
*/
//用户信息
//用户表
type User struct {
	Uid      uint   `gorm:"column:UID;primaryKey"`
	Password string `gorm:"type:VARCHAR(512) NOT NULL"`
}

func (User) TableName() string {
	return "user"
}

type User_info struct {
	Uid           uint   `gorm:"column:UID;type:INT PRIMARY KEY AUTO_INCREMENT;primaryKey"`
	Name          string `gorm:"type:VARCHAR(512) NOT NULL"`
	FollowCount   uint   `gorm:"type:INTEGER"`
	FollowerCount uint   `gorm:"type:INTEGER"`
}

func (User_info) TableName() string {
	return "user_info"
}

// 视频信息
// 视频表
type Video_url struct {
	VID      uint   `gorm:"column:VID"`
	PlayUrl  string `gorm:"type:CHAR(255)"`
	CoverUrl string `gorm:"type:CHAR(255)"`
}

func (Video_url) TableName() string {
	return "video_url"
}

type Video_info struct {
	VID         uint      `gorm:"column:VID;type:INT PRIMARY KEY AUTO_INCREMENT;primaryKey"`
	Title       string    `gorm:"type:VARCHAR(255)"`
	PlayNum     uint      `gorm:"type:INTEGER"`
	LikeNum     uint      `gorm:"type:INTEGER"`
	PublishTime time.Time `gorm:"type:timestamp"`
	Author      string    `gorm:"type:VARCHAR(255)"`
	CommentNum  uint      `gorm:"type:INTEGER"`
}

func (Video_info) TableName() string {
	return "video_info"
}

// 视频对用户记录的信息
// 点赞表
type Like struct {
	VID  uint  `gorm:"column:VID;type:INT NOT NULL"`
	UID  uint  `gorm:"column:UID;type:INT NOT NULL"`
	Flag int32 `gorm:"type:INTEGER"`
}
type Comment struct {
	VID         uint      `gorm:"column:VID;type:INT NOT NULL"`
	UID         uint      `gorm:"column:UID;type:INT NOT NULL"`
	CommentText string    `gorm:"type:VARCHAR(512)"`
	CommentTime time.Time `gorm:"type:timestamp"`
}

// 用户对用户记录的信息
// 关注表
type Follow struct {
	UID      uint `gorm:"column:UID;type:INT NOT NULL"`
	FollowID uint `gorm:"column:FOLLOW_ID;type:INT NOT NULL"`
}

// 粉丝表/被关注表
type Follower struct {
	UID        uint `gorm:"column:UID;type:INT NOT NULL"`
	FollowerID uint `gorm:"column:FOLLOWER_ID;type:INT NOT NULL"`
}
