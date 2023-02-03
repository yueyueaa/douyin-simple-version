package main

import (
	"douyin-simple-version/service"

	"github.com/gin-gonic/gin"
)

func main() {
	initService()

	go service.RunMessageServer()

	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// 初始化服务
func initService() {

}
