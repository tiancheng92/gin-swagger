# gin-swagger (swagger-ui版本：v3.48.0)

* 为减小包体积，删除了[swagger-ui](https://github.com/swagger-api/swagger-ui)中的Source Map文件,仅保留必要的图标、js、css文件。
* 使用了golang 1.16 的新特性 File Embed，取代原[swaggo/gin-swagger](https://github.com/swaggo/gin-swagger)使用的webdav
* 预计比原[swaggo/gin-swagger](https://github.com/swaggo/gin-swagger)内存占用少35M+

## 用法
1. 下载 [Swag](https://github.com/swaggo/swag)
```sh
$ go get -u github.com/swaggo/swag/cmd/swag
```
2. 初始化 [Swag](https://github.com/swaggo/swag) 
```sh
$ swag init
```
3. 下载 [gin-swagger](https://github.com/tiancheng92/gin-swagger)
```sh
$ go get -u github.com/tiancheng92/gin-swagger
```
4. 代码中引入
```go
import "github.com/tiancheng92/gin-swagger"
```

## 实例
```go
package main

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/tiancheng92/gin-swagger"

	_ "./docs"
)

// @title Swagger Example API
// @version 1.0
// @BasePath /api
func main() {
	r := gin.New()
	r.GET("/swagger/*any", ginSwagger.WrapHandler())
	r.Run()
}
```

## 注意
1. golang版本必须大于等于1.16.0
