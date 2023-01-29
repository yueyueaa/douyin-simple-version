# douyin-simple-version

## 框架文档

Web框架: [Gin框架](https://gin-gonic.com/zh-cn/docs/)

Gorm文档: [Gorm文档](https://gorm.io/zh_CN/docs/)

gRPC文档: [gRPC文档](https://grpc.io/docs/languages/go/)

项目功能更新文档: [项目功能更新文档](https://shimo.im/file-invite/HNg9n5RRAcCPBXexFmYxKSDoPEBy6/)

## 功能简介

接口功能不完善，仅作为示例

* 用户登录数据保存在内存中，单次运行过程中有效
* 视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可

## 测试

test 目录下为不同场景的功能测试case，可用于验证功能实现正确性

其中 common.go 中的 _serverAddr_ 为服务部署的地址，默认为本机地址，可以根据实际情况修改

测试数据写在 demo_data.go 中，用于列表接口的 mock 测试