package model

// Supporter
type Supporter struct {
	GORMBase
	Answer    Answer  `json:"-" gorm:"ForeignKey:AnswerID"`
	AnswerID  int     `json:"answerID"`
	Profile   Profile `json:"profile" gorm:"ForeignKey:ProfileID"`
	ProfileID int     `json:"profileID"`
}