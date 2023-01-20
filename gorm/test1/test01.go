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
	Dbname := "gorm"    //数据库名
	timeout := "10s"    //连接超时，10秒

	// 要显示的日志等级
	mysqlLogger = logger.Default.LogMode(logger.Info)

	// root:root@tcp(127.0.0.1:3306)/gorm?
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//跳过默认事务
		SkipDefaultTransaction: true,
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

type Student struct {
	ID    uint    `gorm:"size:10;comment:主键"`
	Name  string  `gorm:"size:16;comment:名字"`
	Age   int     `gorm:"size:3;comment:年龄"`
	Email *string `gorm:"size:128;comment:邮箱"`
}

func main() {
	// 连接成功
	fmt.Println(DB)
	DB.AutoMigrate(&Student{})
	email := "5345345433@qq.com"
	// 创建记录
	student := Student{
		Name:  "枫枫",
		Age:   21,
		Email: &email,
	}
	DB.Create(&student)
	//批量插入
	var studentList []Student
	for i := 0; i < 100; i++ {
		studentList = append(studentList, Student{
			Name:  fmt.Sprintf("测试%d号", i+1),
			Age:   21,
			Email: &email,
		})
	}
	DB.Create(&studentList)

}
