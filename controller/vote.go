package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qa/logic"
	"qa/model"
	"qa/util"
)



//回答投票
func VoteForAnswer(c *gin.Context) {
	var v model.AnswerVote
	if err:=c.ShouldBindJSON(v);err!=nil{
		c.JSON(http.StatusOK, gin.H{
			"code":    util.VoteInvalidParams,
			"message": util.VoteInvalidParams.Msg(),
		})
		return
	}
	userID := c.MustGet("userID").(uint64)
	v.UserID=userID
	logic.UpdateVoteChan<-v

	c.JSON(http.StatusOK, gin.H{
		"code":    util.CodeSuccess,
		"message": util.CodeSuccess.Msg(),
	})


}