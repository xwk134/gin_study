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

type User struct {
	ID       uint
	Name     string
	Age      int
	Gender   bool
	UserInfo UserInfo //用户详情信息

}

type UserInfo struct {
	User   *User
	UserID uint //外键
	ID     uint
	Addr   string
	Like   string
}

func main() {
	// 连接成功
	fmt.Println(DB)
	DB.AutoMigrate(&User{}, &UserInfo{})
	// 添加用户，自动添加用户详情
	DB.Create(&User{
		Name:   "测试",
		Age:    32,
		Gender: true,
		UserInfo: UserInfo{
			Addr: "湖南省",
			Like: "打篮球",
		},
	})

	var user User
	DB.Take(&user, 2)
	DB.Create(&UserInfo{
		User: &user,
		Addr: "南京市",
		Like: "吃饭",
	})

	//通过主表查副表
	DB.Preload("UserInfo").Take(&user)
	fmt.Println(user)
}
