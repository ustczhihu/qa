package model

import (
	"qa/dao"
	"qa/util"
)

type AnswerVote struct {
	GORMBase
	QuestionID uint64 `json:"questionId,string" binding:"required"`
	AnswerID  uint64  `json:"answerId,string" binding:"required"`
	UserID    uint64  `json:"userId,string"`
	Direction float64 `json:"direction binding:"oneof=1 0 -1"` //1-赞，0-没有，-1-踩
}

//查询所有vote id
func GetAllVoteId() (voteList []AnswerVote, code util.MyCode) {
	err := dao.DB.Select("question_id,answer_id,user_id,direction").
		Find(&voteList).
		Error
	if err != nil {
		code = util.VoteDataBaseError
		return
	}
	code = util.CodeSuccess
	return
}

//更新vote
func (v *AnswerVote) Update() (err error) {
	if err := dao.DB.Model(&v).Update(&v).Error; err != nil {
		util.Log.Error(err)
	}
	return
}