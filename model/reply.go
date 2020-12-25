package model


// Reply 回复
type Reply struct {
	GORMBase
	Content            string  `json:"content" label:"对应回答ID"`
	Type               int     `json:"type"`
	Comment            Comment `json:"-" gorm:"ForeignKey:CommentID"`
	CommentID          int     `json:"commentID" label:"对应评论ID"`
	ReplyFromProfile   Profile `json:"replyFromProfile" gorm:"ForeignKey:ReplyFromProfileID"`
	ReplyFromProfileID int     `json:"replyFromProfileID" label:"回复来源者ID"`
	ReplyToProfile     Profile `json:"replyToProfile" gorm:"ForeignKey:ReplyToProfileID"`
	ReplyToProfileID   int     `json:"replyToProfileID" label:"回复对方ID"`
}

