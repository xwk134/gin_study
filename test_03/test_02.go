package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// 多个中间件
func m2(c *gin.Context) {
	fmt.Println("m2 ...in")
	//中间件放行
	c.Next()
	fmt.Println("m2 ...out")
}

func m3(c *gin.Context) {
	fmt.Println("m3 ...in")
	//中间件放行
	c.Next()
	fmt.Println("m3 ...out")
}

func main() {
	router := gin.Default()
	router.GET("/index", m2, func(c *gin.Context) {
		fmt.Println("index......")
		c.JSON(200, gin.H{"msg": "响应数据"})
		c.Next()
		fmt.Println("index ...out")
	}, m3)
	router.Run(":8080")
}
