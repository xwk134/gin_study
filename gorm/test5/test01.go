package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB
var mysqlLogger logger.Interface

func init() {
	username := "root"  //账号
	password := "admin" //密码
	host := "127.0.0.1" //数据库地址，可以是Ip或者域名
	port := 3306        //数据库端口
	Dbname := "test"    //数据库名
	timeout := "10s"    //连接超时，10秒

	// 要显示的日志等级
	mysqlLogger = logger.Default.LogMode(logger.Info)

	// root:root@tcp(127.0.0.1:3306)/gorm?
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//跳过默认事务
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "f_",  // 表名前缀
			SingularTable: false, // 单数表名
			NoLowerCase:   false, // 关闭小写转换
		},
		Logger: mysqlLogger,
	})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	DB = db
}

type User struct {
	ID       uint
	Name     string    `gorm:"size:8"`
	Articles []Article //用户拥有的文章列表
}

type Article struct {
	ID     uint
	Title  string `gorm:"size:16"`
	UserID uint
	User   User
}

func main() {
	// 连接成功
	fmt.Println(DB)
	// DB.AutoMigrate(&User{}, &Article{})
	// // 创建用户,并且创建文章
	// a1 := Article{Title: "python"}
	// a2 := Article{Title: "golang"}
	// user := User{Name: "峰峰", Articles: []Article{a1, a2}}
	// DB.Create(&user)
	// // 创建文章，关联已有用户
	// a3 := Article{Title: "golang零基础入门", UserID: 1}
	// DB.Create(&a3)

	// 预加载
	// 加载文章列表
	var user User
	DB.Preload("Articles").Take(&user, 1)
	fmt.Println(user)

	// 查询文章，显示文章用户的信息
	var article Article
	DB.Preload("User").Take(&article, 1)
	fmt.Println(article)

	//嵌套预加载
	//查询文章，显示用户，并显示用户关联的所有文章
	DB.Preload("User.Articles").Take(&article, 1)
	fmt.Println(article)

	//带条件的预加载
	//查询用户下的所有文章列表，过滤某些文章
	DB.Preload("Articles", "id = ?", 1).Take(&user, 1)
	fmt.Println(user)

	//级联删除
	//删除用户，与用户关联的文章也会删除
	DB.Take(&user, 1)
	DB.Select("Articles").Delete(&user)

}
