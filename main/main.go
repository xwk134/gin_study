package main

import (
	_ "gin_study/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	logger "gin_study/test_05"

	"gin_study/logrus/log"
	"gin_study/logrus/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title API文档
// @version 1.0
// @description API文档
// @host 127.0.0.1:8081
// @BasePath /
func main() {
	logger.Log.Info("This is a logger")
	logger.Log.Error("is logger")
	log.InitFile("log", "server")
	router := gin.New()
	router.Use(middleware.LogMiddleware())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/", func(c *gin.Context) {
		logrus.Infoln("来了")
		c.JSON(200, gin.H{"msg": "你好"})
	})
	router.GET("/api/userlist", UserList)
	router.Run(":8081")
}

type Response struct {
	Code int    `json:"code"` //响应码
	Msg  string `json:"msg"`  //描述
	Data any    `json:"data"` //具体的数据
}

// UserList 用户列表
// @Tags 用户管理
// @Sumnary 用户列表
// @Description 返回一个用户列表，可根据查询参数指定
// @Param limit query string false "返回多少条"
// @Router /api/users [get]
// @Produce json
// @Success 200 {object} Response
func UserList(c *gin.Context) {
	c.JSON(200, Response{0, "成功", 21})
}
