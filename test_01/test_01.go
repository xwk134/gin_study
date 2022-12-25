package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(context *gin.Context) {
	context.String(http.StatusOK, "helloworld")
}
func _json(c *gin.Context) {
	//结构体转json
	type UserInfo struct {
		UserName string `json:"user_name"`
		Age      int    `json:"age"`
		Password string `json:"-"` //忽略转换为json
	}
	User := UserInfo{"张三", 24, "123456"}
	c.JSON(http.StatusOK, User)
}

func main() {
	//创建一个默认的路由
	router := gin.Default()
	//加载目录下静态资源
	router.StaticFS("/static", http.Dir("static/static"))
	// 配置单个文件， 网页请求的路由，文件的路径
	router.StaticFile("/dy", "static/dy.png")
	//绑定路由规则和路由函数
	router.GET("/index", Index)
	router.GET("/json", _json)
	//启动监听 修改ip为内网ip
	router.Run("0.0.0.0:8080")

}
