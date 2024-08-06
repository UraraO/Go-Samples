package backup

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User_ struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(32);index:user_name;not null;unique"`
	Password string `gorm:"column:password;type:varchar(255);not null"`
}

type User struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(32);index:user_name;not null;unique"`
	Password string `gorm:"column:password;type:varchar(255);not null"`
}

func (u User) TableName() string {
	return "users"
}

const (
	MySQLUsername = "root"
	MySQLPassword = "123456"
	MySQLHost     = "localhost"
	MySQLPort     = 3306
	MySQLDatabase = "K_file_server"
)

func GormTest() {
	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", MySQLUsername, MySQLPassword, MySQLHost, MySQLPort, MySQLDatabase)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, error=" + err.Error())
	}
	// 插入数据
	user := &User{
		Username: "UraraO1",
		Password: "Password",
	}
	if err := db.Create(user).Error; err != nil {
		fmt.Println("插入失败", err)
	}
	// 查询
	var auser User
	db.First(&auser, "username = ?", "UraraO") // 手动条件
	fmt.Println(auser)
	db.First(&auser, 13) // 按主键查询
	fmt.Println(auser)
	db.Where("username = ?", "UraraO").First(&auser) // where
	fmt.Println(auser)
	db.Where(&User{Username: "UraraO"}).Find(&auser) // Find 查询多个
	fmt.Println(auser)

	// 更新
	auser = User{
		Username: "UraraO",
		Password: "NewPassword2",
	}
	db.Table("users").Where("username = ?", "UraraO").Update("password", "NewPassword") // 更新一个字段
	db.Table("users").Where("username = ?", "UraraO").Updates(auser)                    // 更新一条数据

	// 删除
	db.Where("username = ?", "UraraO").Delete(&User{})
}
