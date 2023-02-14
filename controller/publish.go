package controller

import (
	"bytes"
	"douyin-simple-version/service/middleware"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	// "github.com/sirupsen/logrus"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	Domain         = "douyin.yoitsu-holo.top:20080"
	VideoUrlPrefix = "http://" + Domain
)

func Insert_newvideo(Title string, Cover_Url string, Play_Url string, user_uid uint, publishdate time.Time) (err error) {
	db, err := middleware.InitDB()
	if err != nil {
		fmt.Printf("[DB ERR] Insert_newuser\t%v\n", err)
		return
	}
	publishvideo_url := middleware.Video_url{
		PlayUrl:  Play_Url,
		CoverUrl: Cover_Url,
	}
	publishvideo_info := middleware.Video_info{
		Title:       Title,
		PlayNum:     0,
		LikeNum:     0,
		PublishTime: publishdate,
		Author:      user_uid,
		CommentNum:  0,
	}
	db.Create(&publishvideo_info)
	publishvideo_url.VID = publishvideo_info.VID
	db.Create(&publishvideo_url)
	return nil
}
func ErrorResponse(err error) Response {
	return Response{
		StatusCode: 1,
		StatusMsg:  err.Error(),
	}
}

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
type Error struct {
	Msg string
}

func (e Error) Error() string {
	return e.Msg
}
func GenerateUUID() string {
	return uuid.NewString()
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {

	token := c.PostForm("token")
	Title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, PublishFunc(token, Title, data, c))
}
func GetSnapshot(videoPath, snapshotPath string, frameNum int) (snapshotName string, err error) {

	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}
	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	err = imaging.Save(img, snapshotPath+".png")
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	names := strings.Split(snapshotPath, "\\")
	snapshotName = names[len(names)-1] + ".png"
	return
}
func PublishFunc(token, Title string, data *multipart.FileHeader, c *gin.Context) Response {
	//检查文件是否为空
	if data == nil {
		return ErrorResponse(Error{Msg: "empty data file"})
	}
	//检查后缀名
	ext := filepath.Ext(data.Filename)
	if ext != ".mp4" {
		return ErrorResponse(Error{Msg: "unsupported file extension"})
	}
	//存文件
	filepath.Base(data.Filename)
	fileName := GenerateUUID()
	videoFileName := fmt.Sprintf("%s%s", fileName, ext)

	// 判断文件夹是否存在
	var dir = "./publish/"
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}

	saveFile := filepath.Join(dir, videoFileName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		return ErrorResponse(err)
	}
	//生成视频信息
	PlayUrl := VideoUrlPrefix + videoFileName
	//封面
	_, err = GetSnapshot("./publish/"+videoFileName, "./publish/"+"-"+videoFileName, 1)
	if err != nil {
		return ErrorResponse(err)
	}
	//生成基本信息
	var author uint
	if user, exist := usersLoginInfo[token]; exist {
		author = uint(user.Id)
	}
	var CoverUrl string
	// var authorid =12
	CoverUrl = VideoUrlPrefix + videoFileName + ".png"
	publish_time := time.Now()
	err = Insert_newvideo(Title, CoverUrl, PlayUrl, author, publish_time)
	if err != nil {
		return ErrorResponse(err)
	}
	return Response{
		StatusCode: 0,
		StatusMsg:  "success",
	}

}
func PublishList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
