package controller

import (
	"bytes"
	"douyin-simple-version/public"
	"douyin-simple-version/service/middleware"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"path"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/sftp"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"golang.org/x/crypto/ssh"
)

const (
	Domain         = "douyin.yoitsu-holo.top:20080/"
	VideoUrlPrefix = "http://" + Domain
	host           = ""
	port           = 20022
	username       = ""
	password       = ""
)

func connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey),
	}
	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

func uploadFile(sftpClient *sftp.Client, localFilePath string, remotePath string) {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("os.Open error : ", localFilePath)
		log.Fatal(err)
	}
	defer srcFile.Close()
	var remoteFileName = path.Base(localFilePath)
	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		fmt.Println("sftpClient.Create error : ", path.Join(remotePath, remoteFileName))
		log.Fatal(err)
	}
	defer dstFile.Close()
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("ReadAll error : ", localFilePath)
		log.Fatal(err)
	}
	dstFile.Write(ff)
	fmt.Println(localFilePath + " copy file to remote server finished!")
}
func DoBackup(host string, port int, userName string, password string, localPath string, remotePath string) {
	var (
		err        error
		sftpClient *sftp.Client
	)
	start := time.Now()
	sftpClient, err = connect(userName, password, host, port)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()
	_, errStat := sftpClient.Stat(remotePath)
	if errStat != nil {
		log.Fatal(remotePath + " remote path not exists!")
	}
	uploadFile(sftpClient, localPath, remotePath)
	elapsed := time.Since(start)
	fmt.Println("elapsed time : ", elapsed)
}

func Insert_newvideo(Title string, Cover_Url string, Play_Url string, user_uid int64, publishdate time.Time) (publishvideo_info middleware.Video_info) {
	db, err := middleware.InitDB()
	if err != nil {
		fmt.Printf("[DB ERR] Insert_video\t%v\n", err)
		return
	}
	publishvideo_url := middleware.Video_url{
		PlayUrl:  Play_Url,
		CoverUrl: Cover_Url,
	}
	publishvideo_info = middleware.Video_info{
		AuthorID:      user_uid,
		Title:         Title,
		FavoriteCount: 0,
		CommentCount:  0,
		PublishTime:   publishdate,
	}
	db.Create(&publishvideo_info)
	publishvideo_url.VID = publishvideo_info.VID
	db.Create(&publishvideo_url)
	return publishvideo_info
}
func ErrorResponse(err error) public.Response {
	return public.Response{
		StatusCode: 1,
		StatusMsg:  err.Error(),
	}
}

type VideoListResponse struct {
	public.Response
	VideoList []public.Video `json:"video_list"`
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
		c.JSON(http.StatusOK, public.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	if user, exist := usersLoginInfo[token]; exist {
		author := user.Id
		c.JSON(http.StatusOK, PublishFunc(token, Title, data, c, author))
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: public.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
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
func PublishFunc(token, Title string, data *multipart.FileHeader, c *gin.Context, auth int64) public.Response {
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
	_, err = GetSnapshot("./publish/"+videoFileName, "./publish/"+"prc"+videoFileName, 1)
	if err != nil {
		return ErrorResponse(err)
	}
	//生成基本信息
	// var author int64
	//本地文件上传到服务器
	var CoverUrl string
	// var authorid =12
	CoverUrl = VideoUrlPrefix + "prc" + videoFileName + ".png"
	publish_time := time.Now()
	Insert_newvideo(Title, CoverUrl, PlayUrl, auth, publish_time)
	DoBackup(host, port, username, password, "./publish/"+videoFileName, "/home/team/")
	DoBackup(host, port, username, password, "./publish/"+"prc"+videoFileName+".png", "/home/team/")
	fmt.Println("完成")
	// cliConf.Upload("./publish/"+"-"+videoFileName, "/root")
	if err != nil {
		return ErrorResponse(err)
	}
	return public.Response{
		StatusCode: 0,
		StatusMsg:  "success",
	}
}
func PublishList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	token := c.Query("token")
	user, exist := usersLoginInfo[token]
	if !exist || user.Id != userId {
		return
	}
	var videos []public.Video
	db, err := middleware.InitDB()
	if err != nil {
		return
	}
	//查找作者所有稿件信息
	var video_infos []middleware.Video_info
	err = db.Select("*").Where("author_id=?", user.Id).Find(&video_infos).Error
	if err != nil {
		return
	}
	for _, info := range video_infos {
		var isFavorite int64
		var url middleware.Video_url
		//查找对应视频url
		err := db.Select("*").Where("VID=?", info.VID).Take(&url).Error
		if err != nil {
			return
		}
		//查找点赞信息
		var favorite middleware.Favorite
		err = db.Select("*").Where("VID=? AND UID=?", info.VID, user.Id).Find(&favorite).Error
		if err != nil {
			isFavorite = 0
		} else {
			isFavorite = favorite.Flag
		}
		videos = append(videos, public.Video{
			int64(info.VID),
			user,
			url.PlayUrl,
			url.CoverUrl,
			int64(info.FavoriteCount),
			int64(info.CommentCount),
			isFavorite != 0,
		})
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: public.Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
