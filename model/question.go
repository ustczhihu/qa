package model

import (
	"qa/dao"
	util "qa/util"
)

// Question 问题
type Question struct {
	GORMBase
	Title   string `json:"title" gorm:"type:varchar(500)"`
	Content string `json:"content" gorm:"type:varchar(2000)"`
	UserID  uint64 `json:"userId,string"`
}

// 创建问题
func (q *Question) Create() util.MyCode {
	if err := dao.DB.Create(&q).Error; err != nil {
		return util.QuestionDataBaseError
	}
	return util.CodeSuccess
}

// 查询所有问题
func GetAllQuestion(pageSize int, pageNum int) (questionList []Question, total int64, code util.MyCode) {

	err := dao.DB.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&questionList).Count(&total).Error
	if err != nil {
		code = util.QuestionDataBaseError
		return
	}
	code = util.CodeSuccess
	return
}
