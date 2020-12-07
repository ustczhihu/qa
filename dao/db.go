package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"qa/config"
)

var DB *gorm.DB

// 连接数据库
func Init() (err error) {

	DB, err = gorm.Open(config.Conf.DB.Driver, config.Conf.DB.Addr)
	DB.LogMode(true)

	if err != nil {
		log.Println(err)
		panic("failed to connect database !")
	}

	DB.SingularTable(true)
	return
}
