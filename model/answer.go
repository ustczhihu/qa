package model

import (
	"log"
	"qa/dao"
	"qa/util"
)

// Answer回答
type Answer struct {
	GORMBase
	Content          string      `json:"content" binding:"required" gorm:"type:varchar(4000);not null" label:"内容"`
	QuestionID       uint64      `json:"questionId,string" gorm:"not null" validate:"required" label:"对应问题ID"`
	Question         Question    `json:"-" gorm:"ForeignKey:QuestionID;associationForeignKey:ID"`
	AnswerProfileID	 uint64      `json:"userId,string" binding:"required" gorm:"not null" validate:"required" label:"回答者ID"`
	AnswerProfile    Profile     `json:"creator" gorm:"ForeignKey:AnswerProfileID;associationForeignKey:UserID"`
	//Type             int         `json:"type"`
	//Comments         []Comment   `json:"-"`
	//CommentsCounts   int         `json:"commentsCounts" label:"评论数"`
	//Supporters       []Supporter `json:"-"`
	//SupportersCounts int         `json:"supportersCounts" label:"赞同数"`
	//Supported        bool        `json:"supported" gorm:"-"`
}

// Create func 创建回答
func (a *Answer) Create() util.MyCode {
	if err := dao.DB.Create(&a).Error; err != nil {
		return util.AnswerDataBaseError
	}
	return util.CodeSuccess
}

// Delete func 删除回答
func (a *Answer) Delete() util.MyCode {
	if err := dao.DB.Delete(&a).Error; err != nil {
		return util.AnswerDataBaseError
	}
	return util.CodeSuccess
}

// Update func 更新回答
func (a *Answer) Update() util.MyCode {
	if err := dao.DB.Model(&a).Updates(&a).Error; err != nil {
		return util.AnswerDataBaseError
	}
	return util.CodeSuccess
}

// Get func 根据answerid查询特定的回答
func (a *Answer) Get() (code util.MyCode) {

	if err := dao.DB.Where(&a).Preload("Question").Preload("AnswerProfile").First(&a).Error; err != nil {
		code = util.AnswerDataBaseError
		return
	}
	code = util.CodeSuccess
	return
}


// GetList func 根据问题和用户查询回答列表
func (a *Answer) GetList(pageSize int, pageNum int) (answers []Answer, total int64, code util.MyCode) {

	err := dao.DB.Preload("Question").Preload("AnswerProfile").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&answers).Count(&total).Error;
	if err != nil {
		code = util.AnswerDataBaseError
		return
	}
	code = util.CodeSuccess
	return
}

// GetOrderList func 根据问题和用户查询回答列表(按照指定的order排列)
func (a *Answer) GetOrderList(limit int, offset int, order string) (answers []Answer, total int64, code util.MyCode) {

	if err := dao.DB.Offset(offset).Limit(limit).Preload("Question").Preload("AnswerProfile").Order(order).Find(&answers, a).Error; err != nil {
		log.Print(err)
	}

	return
}

