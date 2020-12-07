package model

import (
	"encoding/base64"
	"golang.org/x/crypto/scrypt"
	"log"
	"qa/dao"
	"qa/util"
)

type User struct {
	GORMBase
	Username string `form:"username" json:"username" gorm:"type:varchar(500)"`
	Password string `form:"password" json:"password" gorm:"type:varchar(500)"`
}

// 查询用户
func (u *User) Get() (user User, code util.MyCode) {
	if err := dao.DB.Where(&u).First(&user).Error; err != nil {
		code = util.UserNotExist
	} else {
		code = util.CodeSuccess
	}
	return
}

// 创建用户
func (u *User) Create() util.MyCode {
	if err := dao.DB.Create(&u).Error; err != nil {
		return util.UserDataBaseError
	}
	return util.CodeSuccess
}

// 密码加密
func (u *User) BeforeSave() (err error) {
	u.Password = ScryptPw(u.Password)
	return nil
}

func ScryptPw(password string) string {
	const KeyLen = 10
	salt := make([]byte, 8)
	salt = []byte{12, 32, 4, 6, 66, 22, 222, 11}

	HashPw, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen)
	if err != nil {
		log.Fatal(err)
	}
	fpw := base64.StdEncoding.EncodeToString(HashPw)
	return fpw
}

// 登录验证
func (u *User) CheckLogin() (user User, code util.MyCode) {
	err := dao.DB.Where("username = ?", u.Username).First(&user).Error

	if err != nil {
		code = util.UserNotExist
		return
	}
	if ScryptPw(u.Password) != user.Password {
		code = util.UserInvalidPassword
		return
	}
	code = util.CodeSuccess
	return
}
