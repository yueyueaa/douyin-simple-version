package middleware

import (
	"time"
)

// 用户信息
// 用户表
type User struct {
	UID      int64  `gorm:"column:UID;primaryKey"`
	Password string `gorm:"type:VARCHAR(512) NOT NULL"`
}

func (User) TableName() string {
	return "user"
}

type User_info struct {
	UID           int64  `gorm:"column:UID;type:INT PRIMARY KEY AUTO_INCREMENT;primaryKey"`
	Name          string `gorm:"type:VARCHAR(512) NOT NULL"`
	FollowCount   int64  `gorm:"type:INT"`
	FollowerCount int64  `gorm:"type:INT"`
}

func (User_info) TableName() string {
	return "user_info"
}

// 视频信息
// 视频表
type Video_url struct {
	VID      int64  `gorm:"column:VID"`
	PlayUrl  string `gorm:"type:CHAR(255)"`
	CoverUrl string `gorm:"type:CHAR(255)"`
}

func (Video_url) TableName() string {
	return "video_url"
}

type Video_info struct {
	VID           int64     `gorm:"column:VID;type:INT PRIMARY KEY AUTO_INCREMENT;primaryKey"`
	AuthorID      int64     `gorm:"column:author_id;type:INT"`
	Title         string    `gorm:"type:VARCHAR(255)"`
	FavoriteCount int64     `gorm:"column:favorite_count;type:INT"`
	CommentCount  int64     `gorm:"column:comment_count;type:INT"`
	PublishTime   time.Time `gorm:"column:publish_date;type:timestamp"`
}

func (Video_info) TableName() string {
	return "video_info"
}

// 视频对用户记录的信息
// 点赞表
type Favorite struct {
	VID  int64 `gorm:"column:VID;type:INT NOT NULL"`
	UID  int64 `gorm:"column:UID;type:INT NOT NULL"`
	Flag int64 `gorm:"type:INT NOT NULL"`
}

func (Favorite) TableName() string {
	return "favorites"
}

type Comment struct {
	CID        int64     `gorm:"column:CID;type:INT PRIMARY KEY AUTO_INCREMENT;primaryKey"`
	VID        int64     `gorm:"column:VID;type:INT NOT NULL"`
	UID        int64     `gorm:"column:UID;type:INT NOT NULL"`
	Content    string    `gorm:"type:VARCHAR(512)"`
	CreateDate time.Time `gorm:"type:timestamp"`
}

func (Comment) TableName() string {
	return "comments"
}

// 用户对用户记录的信息
// 关注表
type Follow struct {
	UID      int64 `gorm:"column:UID;type:INT NOT NULL"`
	FollowID int64 `gorm:"column:FOLLOW_ID;type:INT NOT NULL"`
	Flag     int64 `gorm:"type:INT"`
}

func (Follow) TableName() string {
	return "follows"
}

// 粉丝表/被关注表
type Follower struct {
	UID        int64 `gorm:"column:UID;type:INT NOT NULL"`
	FollowerID int64 `gorm:"column:FOLLOWER_ID;type:INT NOT NULL"`
	Flag       int64 `gorm:"type:INT"`
}

func (Follower) TableName() string {
	return "followers"
}
