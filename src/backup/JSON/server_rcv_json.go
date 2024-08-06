package backup

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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

func SQLserver_test() {
	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", MySQLUsername, MySQLPassword, MySQLHost, MySQLPort, MySQLDatabase)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, error=" + err.Error())
	}
	// ////////////////////
	r := gin.Default()
	r.POST("/api/register", func(c *gin.Context) {
		var user User
		if err := c.ShouldBind(&user); err != nil {
			fmt.Println("gin.Context.ShouldBind ERROR,", err)
		}
		fmt.Println(user)
		if err = db.Create(&user).Error; err != nil {
			fmt.Println("插入失败", err)
		}
	})
	r.Run("127.0.0.1:8080")
}
