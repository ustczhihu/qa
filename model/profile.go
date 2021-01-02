package model

import (
	"errors"
	"math/rand"
	"qa/dao"
	"qa/util"
	"strconv"
)

// Profile 用户简介
type Profile struct {
	GORMBase
	Nickname  string `form:"nickname" json:"nickname" gorm:"type:varchar(500)"`
	Gender    int    `form:"gender" json:"gender" gorm:"type:int;DEFAULT:0;"`
	Desc      string `form:"desc" json:"desc" gorm:"type:varchar(1000)"`
	AvatarUrl string `form:"avatarUrl" json:"avatarUrl" gorm:"type:varchar(1000)"`
	UserID    uint64 `form:"userId" json:"userId,string"`
}

// 查询porfile by userid
func GetByUserID(userID uint64) (profile Profile, code util.MyCode) {
	if err := dao.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		code = util.UserDataBaseError
	} else {
		code = util.CodeSuccess
	}
	return
}

// 设置profile默认值
func (p *Profile) BeforeSave() (err error) {
	if p.UserID == 0 {
		err = errors.New("userID不能为空")
		return
	}
	if p.ID == 0 {
		p.ID, _ = util.GetID()
	}
	if p.Nickname == "" {
		p.Nickname = "user_" + strconv.Itoa(rand.Int())
	}
	if p.Desc == "" {
		p.Desc = "I am " + p.Nickname
	}
	if p.AvatarUrl == "" {
		p.AvatarUrl = "https://s1.ax1x.com/2018/04/04/C9c2GV.png"
	}
	return
}

// 创建profile
func (p *Profile) Create() util.MyCode {
	if err := dao.DB.Create(&p).Error; err != nil {
		return util.UserDataBaseError
	}
	return util.CodeSuccess
}

// 更新profile
func (p *Profile) UpdateProfile() util.MyCode{
	var maps = make(map[string]interface{})
	maps["nickname"] = p.Nickname
	maps["gender"] = p.Gender
	maps["desc"] = p.Desc
	maps["avatar_url"] = p.AvatarUrl


	err := dao.DB.Model(&p).Where("user_id = ?", p.UserID).Update(maps).Error
	if err != nil {
		return util.ProfileSaveFail
	}
	return util.CodeSuccess
}