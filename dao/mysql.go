package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"qa/config"
	"time"
)

var DB *gorm.DB

// 连接数据库
func InitDB() (err error) {

	DB, err = gorm.Open(config.Conf.DB.Driver, config.Conf.DB.Addr)
	DB.LogMode(true)

	if err != nil {
		log.Println(err)
		panic("failed to connect database !")
	}

	DB.SingularTable(true) //以实现结构体名为非复数形式

	//设置最大闲置连接数
	DB.DB().SetMaxIdleConns(10)

	//设置最大连接数
	DB.DB().SetMaxOpenConns(100)

	//设置连接的最大可复用时间
	DB.DB().SetConnMaxLifetime(10 * time.Second)

	return
}
