package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// 全局注册中间件

type User struct {
	Name string
	Age  int
}

func m10(c *gin.Context) {
	fmt.Println("m1 ...in")
	c.Set("name", User{"枫枫", 21})
	c.Next()
	fmt.Println("m1 ...out")
}

func main() {
	router := gin.Default()

	router.Use(m10)
	router.GET("/index", func(c *gin.Context) {
		fmt.Println("index ...in")
		name, _ := c.Get("name")
		user := name.(User)
		fmt.Println(user.Name, user.Age)
		c.JSON(200, gin.H{"msg": "index"})
		fmt.Println("index ...out")
	})

	router.Run(":8080")

}
