# gin-swagger

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

	"github.com/swaggo/gin-swagger/example/basic/api"

	_ "xxx/docs"
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
