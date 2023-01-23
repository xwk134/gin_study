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

type User struct {
	ID     uint    `gorm:"size:10;comment:主键" json:"id"`
	Name   string  `gorm:"size:16;comment:名字" json:"name"`
	Age    int     `gorm:"size:3;comment:年龄" json:"age"`
	Email  *string `gorm:"size:128;comment:邮箱" json:"email"`
	Gender bool
}

func main() {
	// 连接成功
	fmt.Println(DB)
	DB.AutoMigrate(&User{})
	var userList = []User{
		{ID: 1, Name: "李元芳", Age: 32, Email: PtrString("lyf@yf.com"), Gender: true},
		{ID: 2, Name: "张武", Age: 18, Email: PtrString("zhangwu@lly.cn"), Gender: true},
		{ID: 3, Name: "枫枫", Age: 23, Email: PtrString("ff@yahoo.com"), Gender: true},
		{ID: 4, Name: "刘大", Age: 54, Email: PtrString("liuda@qq.com"), Gender: true},
		{ID: 5, Name: "李武", Age: 23, Email: PtrString("liwu@lly.cn"), Gender: true},
		{ID: 6, Name: "李琦", Age: 14, Email: PtrString("liqi@lly.cn"), Gender: false},
		{ID: 7, Name: "晓梅", Age: 25, Email: PtrString("xiaomeo@sl.com"), Gender: false},
		{ID: 8, Name: "如燕", Age: 26, Email: PtrString("ruyan@yf.com"), Gender: false},
		{ID: 9, Name: "魔灵", Age: 21, Email: PtrString("moling@sl.com"), Gender: true},
	}
	DB.Create(&userList)

	var users []User
	// 查询用户名是枫枫的
	DB.Where("name = ?", "枫枫").Find(&users)
	fmt.Println(users)
	// 查询用户名不是枫枫的
	DB.Where("name <> ?", "枫枫").Find(&users)
	fmt.Println(users)
	// 查询用户名包含 如燕，李元芳的
	DB.Where("name in ?", []string{"如燕", "李元芳"}).Find(&users)
	fmt.Println(users)
	// 查询姓李的
	DB.Where("name like ?", "李%").Find(&users)
	fmt.Println(users)
	// 查询年龄大于23，是qq邮箱的
	DB.Where("age > ? and email like ?", "23", "%@qq.com").Find(&users)
	fmt.Println(users)
	// 查询是qq邮箱的，或者是女的
	DB.Where("gender = ? or email like ?", false, "%@qq.com").Find(&users)
	fmt.Println(users)

	// 使用map查询
	DB.Where(map[string]any{"name": "李元芳", "age": 32}).Find(&users)
	// SELECT * FROM `students` WHERE `age` = 0 AND `name` = '李元芳'
	fmt.Println(users)

	// 排除年龄大于23的
	DB.Not("age > 23").Find(&users)
	fmt.Println(users)
	DB.Or("gender = ?", false).Or(" email like ?", "%@qq.com").Find(&users)
	fmt.Println(users)

	DB.Table("f_users").Select("name", "age").Scan(&users)
	fmt.Println(users)

	// 根据年龄排序
	DB.Order("age desc").Find(&users)
	fmt.Println(users)
	// desc    降序
	// asc     升序

	// 分页查询
	// 一页多少条
	limit := 10
	// 第几页
	page := 1
	offset := (page - 1) * limit
	DB.Limit(limit).Offset(offset).Find(&users)
	fmt.Println(users)

	//去重
	var ageList []int
	DB.Table("f_users").Select("distinct age").Scan(&ageList)
	fmt.Println(ageList)

	// 分组查询
	type AggeGroup struct {
		Gender int
		Count  int
		Name   string
	}

	var agge []AggeGroup
	// 查询男生的个数和女生的个数
	//DB.Table("f_users").Select("count(id) as count", "gender", "group_concat(name) as name").Group("gender").Scan(&agge)
	// 执行sql
	DB.Raw(`SELECT count(id) as count, gender, group_concat(name) as name FROM f_users GROUP BY gender`).Scan(&agge)
	fmt.Println(agge)
	DB.Raw(`select * from f_users where age > (select avg(age) from f_users)`).Scan(&users)
	fmt.Println(users)
}

func PtrString(email string) *string {
	return &email
}
