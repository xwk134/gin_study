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
	Dbname := "test1"   //数据库名
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

type Tag struct {
	ID       uint
	Name     string
	Articles []Article `gorm:"many2many:article_tags;"` // 用于反向引用
}

type Article struct {
	ID    uint
	Title string
	Tags  []Tag `gorm:"many2many:article_tags;"`
}

func main() {
	// 连接成功
	fmt.Println(DB)
	// DB.AutoMigrate(&Tag{}, &Article{})
	// // 多对多添加
	// // 添加文章，并创建标签
	// DB.Create(&Article{
	// 	Title: "python基础课程",
	// 	Tags: []Tag{
	// 		{Name: "python"},
	// 		{Name: "基础课程"},
	// 	},
	// })

	// // 添加文章，选择标签
	// var tags []Tag
	// DB.Find(&tags, "name = ?", "基础课程")
	// DB.Create(&Article{
	// 	Title: "golang基础",
	// 	Tags:  tags,
	// })

	// 多对多查询
	// 查询文章，显示文章的标签列表
	var article Article
	DB.Preload("Tags").Take(&article, 2)
	fmt.Println(article)

	//查询标签，显示文章列表
	var tag Tag
	DB.Preload("Articles").Take(&tag, 2)
	fmt.Println(tag)

	// 多对多更新
	// 移除文章的标签
	DB.Preload("Tags").Take(&article, 2)
	DB.Model(&article).Association("Tags").Delete(article.Tags)
	fmt.Println(article)
	// 更新文章标签
	var tags []Tag
	DB.Find(&tags, []int{2, 3, 4})
	DB.Preload("Tags").Take(&article, 2)
	DB.Model(&article).Association("Tags").Replace(tags)
	fmt.Println(article)

}
