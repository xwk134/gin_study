package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// 路由分组注册中间件

func middle(c *gin.Context) {
	fmt.Println("middle ...in")
}

func main() {
	router := gin.Default()

	r := router.Group("/api").Use(middle) // 可以链式，也可以直接r.Use(middle)
	r.GET("/index", func(c *gin.Context) {
		c.String(200, "index")
	})
	r.GET("/home", func(c *gin.Context) {
		c.String(200, "home")
	})

	router.Run(":8080")
}
