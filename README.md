# Gee

## 下载Gee

```go
go get github.com/AYO-Al/myPackage/gee
```
## 运行Gee
A basic example:

```go

import (
	"github.com/AYO-Al/myPackage/gee"
	"net/http"
)

func main() {
	r := gee.Default()

	r.GET("/", func(context *gee.Context) {
		context.Data(http.StatusOK, "Hello Gee")
	})

	r.Run(":8080")
}
```
## 使用HTTP方法

```go

func main() {
	r := gee.Default()

	r.GET("/", HelloGee)
	r.POST("/", HelloGee)

	r.Run(":8080")
}
```
## 分组控制

```go
func main() {
	r := gee.Default()

	// 创建路由组
	v1 := r.Group("/hello")
	{
		v1.GET("/gee", HelloGee)
	}
	
	// 创建子路由组 
	v2 := v1.Group("/v2")
	v2.GET("/gee", HelloGee)

	r.Run(":8080")
}
```
## 中间件

```go
// 自定义中间件
func HelloGee() gee.HandlerFunc {
	return func(context *gee.Context) {
		fmt.Println("Hello")
		context.Next()
		fmt.Println("Gee")
	}
}

func main() {
	r := gee.Default()
	
	// 注册中间件
	r.Use(HelloGee())
	// 创建路由组
	v1 := r.Group("/hello")
	{
		v1.GET("/gee", Hello)
	}

	r.Run(":8080")
}

```
## 渲染HTML

```go
func main() {
	r := gee.Default()
	r.LoadHTMLGlob("gee/templates/*")
	r.Static("/assets", "gee/static")

	r.GET("/gee", func(context *gee.Context) {
		context.HTML(http.StatusOK, "gee.html", nil)
	})

	r.Run(":8080")
}
```
## 获取请求参数

```go
func main() {
	r := gee.Default()
	r.LoadHTMLGlob("gee/templates/*")
	r.Static("/assets", "gee/static")

	// /gee/hello
	r.GET("/gee/:name", func(context *gee.Context) {
		context.Data(http.StatusOK, context.Param("name"))
	})

	// /hello?name=gee
	r.GET("/hello", func(context *gee.Context) {
		context.Data(http.StatusOK, context.Query("name"))
	})

	r.Run(":8080")
}
```
## 修改环境

```go
func main() {
	r := gee.Default()
	// 默认为gee.DebugMode
	gee.SetMode(gee.ReleaseMode)
	
	// /gee/hello
	r.GET("/gee/:name", func(context *gee.Context) {
		context.Data(http.StatusOK, context.Param("name"))
	})
	
	r.Run(":8080")
}
```

