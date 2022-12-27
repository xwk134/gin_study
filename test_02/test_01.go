package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// 单文件上传
func main() {
	router := gin.Default()
	router.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		fmt.Println(file.Filename)
		fmt.Println(file.Size / 1024)
		c.SaveUploadedFile(file, "./uploads/"+file.Filename)
		c.JSON(200, gin.H{"msg": "上传成功"})
	})
	router.Run(":8080")
}
