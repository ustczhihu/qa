package model


type AnswerVote struct {
	GORMBase
	AnswerID uint   `json:"answerId,string"`
	Answer   Answer `json:"-" gorm:"ForeignKey:AnswerID;associationForeignKey:ID"`
	UserID   uint   `json:"userId,string"`
	User     User   `json:"user" gorm:"ForeignKey:UserID;associationForeignKey:ID"`
	UpOrDown bool	`json:"upOrDown"`  // true 赞， false 踩
}