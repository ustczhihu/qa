package main

import (
	"fmt"
	"qa/dao"
	"qa/model"
	"qa/util"
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
	code := u0.Create()
	if code != util.CodeSuccess {
		fmt.Println("user create error!!!")
	}
	var p0 model.Profile
	p0.UserID = u0.ID
	code = p0.Create()
	if code != util.CodeSuccess {
		fmt.Println("profile create error!!!")
	}

	fmt.Println("restarted success !")
	return
}
