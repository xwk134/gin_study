package main

import (
	"fmt"
	"time"

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
	Dbname := "test2"   //数据库名
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

type Article struct {
	ID    uint
	Title string
	Tags  []Tag `gorm:"many2many:article_tags"`
}

type Tag struct {
	ID   uint
	Name string
}

type ArticleTag struct {
	ArticleID uint `gorm:"primaryKey"`
	TagID     uint `gorm:"primaryKey"`
	CreatedAt time.Time
}

func main() {
	// 连接成功
	fmt.Println(DB)
	// 设置Article的Tags表为ArticleTag
	DB.SetupJoinTable(&Article{}, "Tags", &ArticleTag{})
	err := DB.AutoMigrate(&Article{}, &Tag{}, &ArticleTag{})
	fmt.Println(err)
	DB.Create(&Article{
		Title: "flask零基础入门",
		Tags: []Tag{
			{Name: "python"},
			{Name: "后端"},
			{Name: "web"},
		},
	})

	// 添加文章，关联已有标签
	var tags []Tag
	DB.Find(&tags, "name in ?", []string{"python", "web"})
	DB.Create(&Article{
		Title: "flask请求对象",
		Tags:  tags,
	})

	// 给已有文章关联标签
	article := Article{
		Title: "djiango基础",
	}
	DB.Create(&article)
	var at Article

	DB.Find(&tags, "name in ?", []string{"python", "web"})
	DB.Take(&at, article.ID).Association("Tags").Append(tags)

	// 替换已有文章的标签

	DB.Find(&tags, "name in ?", []string{"后端"})
	DB.Take(&article, "title = ?", "django基础")
	DB.Model(&article).Association("Tags").Replace(tags)

	var articles []Article
	// 查询文章列表，显示标签
	DB.Preload("Tags").Find(&articles)
	fmt.Println(articles)

}
