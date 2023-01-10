package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"reflect"
)

// 绑定json 绑定查询参数form 绑定动态参数uri
// binding:"required" 不能为空，并且不能没有这个字段
// 针对字符串的长度
//min 最小长度，如：binding:"min=5"
//max 最大长度，如：binding:"max=10"

type ArticleModel struct {
	Title   string `json:"title" form:"title" uri:"name" binding:"required" msg:"标题不能为空"`
	Content string `json:"content" binding:"required" msg:"文章内容不能为空"`
}

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func _bindJson(c *gin.Context, obj any) (err error) {
	body, _ := c.GetRawData()
	contentType := c.GetHeader("Content-Type")
	switch contentType {
	case "application/json":
		err = json.Unmarshal(body, &obj)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

	}
	return nil
}

// GetValidMsg 返回结构体中的msg参数
func GetValidMsg(err error, obj any) string {
	// 使用的时候，需要传obj的指针
	getObj := reflect.TypeOf(obj)
	// 将err接口断言为具体类型
	if errs, ok := err.(validator.ValidationErrors); ok {
		// 断言成功
		for _, e := range errs {
			// 循环每一个错误信息
			// 根据报错字段名，获取结构体的具体字段
			if f, exits := getObj.Elem().FieldByName(e.Field()); exits {
				msg := f.Tag.Get("msg")
				return msg
			}
		}
	}

	return err.Error()
}

// _getList 文章列表页面
func _getList(c *gin.Context) {
	// 包含搜索、分页
	articleList := []ArticleModel{
		{"Go语言入门", "这个是go语言入门"},
		{"java语言入门", "这个是java语言入门"},
		{"Python语言入门", "这个是python语言入门"},
	}

	c.JSON(200, Response{0, articleList, "成功"})
	logrus.Info("成功")
}

// _getDetail 文章详情
func _getDetail(c *gin.Context) {
	// 获取param中的id
	fmt.Println(c.Param("id"))
	article := ArticleModel{
		"Go语言入门", "这篇文章是《Go语言入门》",
	}
	c.JSON(200, Response{0, article, "成功"})
}

// _create 创建文章
func _create(c *gin.Context) {
	// 接收前端传递过来的json数据
	var article ArticleModel
	err := _bindJson(c, &article)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = c.ShouldBindUri(&article)
	if err != nil {
		// 显示自定义的错误信息
		msg := GetValidMsg(err, &article)
		c.JSON(200, Response{1, article, msg})
		return
	}
	c.JSON(200, Response{0, article, "添加成功"})
}

// _update 编辑文章
func _update(c *gin.Context) {
	fmt.Println(c.Param("id"))
	var article ArticleModel
	err := _bindJson(c, &article)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.ShouldBindUri(&article)
	if err != nil {
		// 显示自定义的错误信息
		msg := GetValidMsg(err, &article)
		c.JSON(200, Response{1, article, msg})
		return
	}
	c.JSON(200, Response{0, article, "修改成功"})
}

// _delete 删除文章
func _delete(c *gin.Context) {
	// 获取请求头
	fmt.Println(c.GetHeader("User-Agent"))
	// 设置响应头
	c.Header("Content-Type", "application/json; charset=utf-8")
	fmt.Println(c.Param("id"))
	c.JSON(200, Response{0, map[string]string{}, "删除成功"})
}

func main() {
	router := gin.Default()
	//添加中间件，主要实现log日志的生成

	router.GET("/articles", _getList)        // 文章列表
	router.GET("/articles/:id", _getDetail)  // 文章详情
	router.POST("/articles/create", _create) // 创建文章
	router.PUT("/articles/:id", _update)     // 修改文章
	router.DELETE("/articles/:id", _delete)  // 删除文章
	router.Run("0.0.0.0:8080")
}
