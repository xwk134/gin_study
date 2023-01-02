package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func indexHandler(c *gin.Context) {
	fmt.Println("index......")
	c.JSON(http.StatusOK, gin.H{
		"msg": "index",
	})
}

// 单独注册中间件

// 定义一个中间件
func m1(c *gin.Context) {
	fmt.Println("m1 in .......")
	//中间件拦截
	c.Abort()

}

func main() {
	r := gin.Default()
	//m1处于indexHandler函数的前面，请求来之后，先走m1,再走index
	r.GET("/index", m1, indexHandler)
	r.Run(":8080")
}
