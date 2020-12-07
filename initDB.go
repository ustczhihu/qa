package main

import (
	"fmt"
	"qa/dao"
	"qa/model"
)

func Init() (err error) {
	if (dao.DB.HasTable(&model.User{})) {
		fmt.Println("db has the table user, so drop it.")
		if err = dao.DB.DropTable(&model.User{}).Error; err != nil {
			return
		}
	}

	if (dao.DB.HasTable(&model.Profile{})) {
		fmt.Println("db has the table profile, so drop it.")
		if err = dao.DB.DropTable(&model.Profile{}).Error; err != nil {
			return
		}
	}

	if (dao.DB.HasTable(&model.Question{})) {
		fmt.Println("db has the table question, so drop it.")
		if err = dao.DB.DropTable(&model.Question{}).Error; err != nil {
			return
		}
	}

	if err = dao.DB.AutoMigrate(&model.User{}).Error; err != nil {
		return
	}
	if err = dao.DB.AutoMigrate(&model.Profile{}).Error; err != nil {
		return
	}
	if err = dao.DB.AutoMigrate(&model.Question{}).Error; err != nil {
		return
	}

	u0 := model.User{Username: "lily", Password: "123456"}
	u0.Create()

	fmt.Println("restarted success !")
	return
}
