package main

import (
	"encoding/json"
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

type Student struct {
	ID    uint    `gorm:"size:10;comment:主键" json:"id"`
	Name  string  `gorm:"size:16;comment:名字" json:"name"`
	Age   int     `gorm:"size:3;comment:年龄" json:"age"`
	Email *string `gorm:"size:128;comment:邮箱" json:"email"`
}

func main() {
	// 连接成功
	fmt.Println(DB)
	var student Student
	DB.Take(&student)
	fmt.Println(student)
	// SELECT * FROM `students` LIMIT 1
	student = Student{}
	DB.First(&student)
	fmt.Println(student)
	// SELECT * FROM `students` ORDER BY `students`.`id` LIMIT 1
	student = Student{}
	DB.Last(&student)
	fmt.Println(student)
	// SELECT * FROM `students` ORDER BY `students`.`id` DESC LIMIT 1

	//根据主键查询
	student = Student{} // 重新赋值
	DB.Take(&student, "4")
	fmt.Println(student)
	//根据其他条件查询
	student = Student{} // 重新赋值
	DB.Take(&student, "name = ?", "测试99号")
	fmt.Println(student)

	student = Student{}
	err := DB.Take(&student, "name = ?", "测试9号").Error
	switch err {
	case gorm.ErrRecordNotFound:
		fmt.Println("查询不到")
	default:
		fmt.Println(student)
	}

	//查询多条记录
	var studentList []Student
	count := DB.Find(&studentList).RowsAffected
	fmt.Print(count)
	DB.Find(&studentList)
	for _, student := range studentList {
		data, _ := json.Marshal(student)
		fmt.Println(string(data))
	}

}
