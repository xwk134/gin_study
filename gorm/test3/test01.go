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
	student.Age = 24
	// 全字段更新
	//DB.Save(&student)

	// 更新指定字段
	DB.Select("age").Save(&student)

	// 批量更新
	DB.Model(&Student{}).Where("age = ?", 21).Update("email", "777777@qq.com")
	// UPDATE `students` SET `email`='is22@qq.com' WHERE age = 21

	// 更新多列
	name := "峰峰"
	email := "88888@qq.com"
	DB.Model(&Student{}).Where("age = ?", 21).Updates(map[string]any{
		"name":  name,
		"email": &email,
	})

	// 删除单个
	DB.Delete(&Student{}, 1)

	student = Student{} // 重新赋值
	DB.Take(&student)
	fmt.Println(student)
	DB.Delete(&student)

}
