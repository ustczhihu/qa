package model

import (
	"github.com/jinzhu/gorm"
	"qa/dao"
	util "qa/util"
)

// Question 问题
type Question struct {
	GORMBase
	Title          string  `json:"title" gorm:"type:varchar(500);not null" validate:"required,min=6,max=30,endswith=?" label:"问题题目"`
	Content        string  `json:"content" gorm:"type:longtext" validate:"max=200" label:"问题描述"`
	AnswerCount    int     `json:"answerCount" gorm:"type:int;DEFAULT:0;"`
	ViewCount      int     `json:"viewCount" gorm:"type:int;DEFAULT:0;"`
	UserID         uint64  `json:"userId,string" gorm:"not null" validate:"required" label:"提问者ID"`
	CreatorProfile Profile `json:"creator" gorm:"foreignKey:UserID;associationForeignKey:UserID"`
}

//查询问题是否存在
func CheckQuestion(title string) util.MyCode {
	var question Question
	dao.DB.Select("id").Where("title=?", title).First(&question)
	if question.ID > 0 {
		return util.QuestionExist
	} else {
		return util.QuestionNotExist
	}
}

//创建问题
func (q *Question) Create() util.MyCode {
	if err := dao.DB.Create(&q).Error; err != nil {
		return util.QuestionDataBaseError
	}
	return util.CodeSuccess
}

//删除问题
func (q *Question) Delete() util.MyCode {
	if err := dao.DB.Delete(&q).Error; err != nil {
		return util.QuestionDataBaseError
	}
	return util.CodeSuccess
}

//更新问题
func (q *Question) Update() util.MyCode {
	if err := dao.DB.Update(&q).Error; err != nil {
		return util.QuestionDataBaseError
	}
	return util.CodeSuccess
}

//查询问题
func (q *Question) Get() (question Question, code util.MyCode) {
	if err := dao.DB.Where(&q).Preload("CreatorProfile").
		First(&question).Error; err != nil {
		code = util.QuestionDataBaseError
		return
	}
	code = util.CodeSuccess
	return
}

//查询所有问题id
func GetAllQuestionId() (questionList []Question, code util.MyCode) {
	err := dao.DB.Select("id,view_count").
		Find(&questionList).
		Error
	if err != nil {
		code = util.QuestionDataBaseError
		return
	}
	code = util.CodeSuccess
	return
}

//查询所有问题
func GetAllQuestion(pageSize int, pageNum int, order string) (questionList []Question, total int64, code util.MyCode) {

	err := dao.DB.Preload("CreatorProfile").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Order(order).
		Find(&questionList).
		Count(&total).
		Error
	if err != nil {
		code = util.QuestionDataBaseError
		return
	}
	code = util.CodeSuccess
	return
}

//查询所有问题ByUserId
func GetAllQuestionByUserId(pageSize int, pageNum int, order string, userId uint64) (questionList []Question, total int64, code util.MyCode) {

	err := dao.DB.Preload("CreatorProfile").
		Limit(pageSize).Offset((pageNum-1)*pageSize).
		Order(order).
		Where("user_id=?", userId).
		Find(&questionList).
		Count(&total).
		Error
	if err != nil {
		code = util.QuestionDataBaseError
		return
	}
	code = util.CodeSuccess
	return
}

//查询所有问题ByTitle
func GetAllQuestionByTitle(pageSize int, pageNum int, order string, title string) (questionList []Question, total int64, code util.MyCode) {

	err := dao.DB.Preload("CreatorProfile").
		Limit(pageSize).Offset((pageNum-1)*pageSize).
		Order(order).
		Where("title LIKE ?", "%"+title+"%").
		Find(&questionList).
		Count(&total).
		Error
	if err != nil {
		code = util.QuestionDataBaseError
		return
	}
	code = util.CodeSuccess
	return
}

//浏览数+1
func (q *Question) IncrView() (err error) {

	if err := dao.DB.Model(&q).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
		util.Log.Error(err)
	}
	return
}