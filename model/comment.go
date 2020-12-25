package model

// Comment评论
type Comment struct {
	GORMBase
	Content          string  `json:"content" label:"内容"`
	Type             int     `json:"type"`
	Answer           Answer  `json:"-" gorm:"ForeignKey:AnswerID"`
	AnswerID         int     `json:"answerID" label:"对应回答ID"`
	CommentProfile   Profile `json:"commentProfile" gorm:"ForeignKey:CommentProfileID"`
	CommentProfileID int     `json:"commentProfileID" label:"评论者ID"`
	Replies          []Reply `json:"-"`
	RepliesCounts    int     `json:"repliesCounts" label:"回复数"`
}
