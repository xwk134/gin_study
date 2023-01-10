package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func LoggerWithFormatter(params gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	statusColor = params.StatusCodeColor()
	methodColor = params.MethodColor()
	resetColor = params.ResetColor()
	return fmt.Sprintf(
		"[ GIN ] %s  | %s %d  %s | \t %s | %s | %s %-7s %s \t  %s\n",
		params.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, params.StatusCode, resetColor,
		params.ClientIP,
		params.Latency,
		methodColor, params.Method, resetColor,
		params.Path,
	)
}

func main() {
	// 输出到文件
	f, _ := os.Create("gin.log")
	//gin.DefaultWriter = io.MultiWriter(f)
	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	//不显示默认日志，改为release模式
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(
		gin.LoggerWithConfig(
			gin.LoggerConfig{Formatter: LoggerWithFormatter},
		),
	)
	router.Run()

}
