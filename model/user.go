package model

import (
	"golang.org/x/crypto/bcrypt"
	"qa/dao"
	"qa/util"
)

type User struct {
	GORMBase
	Username string `json:"username" gorm:"type:varchar(500);not null" validate:"required,min=3,max=12" label:"用户名"`
	Password string `json:"password" gorm:"type:varchar(500);not null" validate:"required,min=6,max=20" label:"密码"`
	Profile Profile `json:"profile" gorm:"foreignKey:ID;associationForeignKey:UserID"`
}

//查询用户是否存在
func CheckUser(name string) util.MyCode{
	var user User
	dao.DB.Select("id").Where("username=?", name).First(&user)
	if user.ID > 0 {
		return util.UserExist
	} else {
		return util.UserNotExist
	}
}

//查询用户
func (u *User) Get() (user User, code util.MyCode) {
	if err := dao.DB.Preload("Profile").Where(&u).First(&user).Error; err != nil {
		code = util.UserNotExist
	} else {
		code = util.CodeSuccess
	}
	return
}

//创建用户
func (u *User) Create() util.MyCode {
	if err := dao.DB.Create(&u).Error; err != nil {
		return util.UserDataBaseError
	}
	return util.CodeSuccess
}

//删除用户
func (u *User) Delete() util.MyCode {
	if err := dao.DB.Delete(&u).Error; err != nil {
		return util.UserDataBaseError
	}
	return util.CodeSuccess
}

//更新用户
func (u *User) Update() util.MyCode {
	if err := dao.DB.Update(&u).Error; err != nil {
		return util.UserDataBaseError
	}
	return util.CodeSuccess
}

//密码加密
func (u *User) BeforeSave() (err error) {
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hash)
	return
}

// 登录验证
func (u *User) CheckLogin() (user User, code util.MyCode) {
	err := dao.DB.Preload("Profile").Where("username = ?", u.Username).First(&user).Error

	if err != nil {
		code = util.UserNotExist
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		code = util.UserInvalidPassword
		return
	}
	code = util.CodeSuccess
	return
}
