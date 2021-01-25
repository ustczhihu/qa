package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qa/logic"
	"qa/model"
	util "qa/util"
	"strconv"
)



// 创建回答
func AddAnswer(c *gin.Context) {
	var a model.Answer
	var code util.MyCode
	if err := c.ShouldBind(&a); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    util.AnswerInvalidParams,
			"message": util.AnswerInvalidParams.Msg(),
		})
		return
	}
	// 验证token中的userid和传过来的userid是否一致
	userID := c.MustGet("userID").(uint64)
	// 不一致则返回 AnswerUserIdNotMatch 错误
	if userID != a.AnswerProfileID {
		code = util.AnswerUserIdNotMatch
		c.JSON(
			http.StatusOK, gin.H{
				"code":    code,
				"message": code.Msg(),
			},
		)
		return
	}
	// 根据userId获取回答者的profile用于返回
	var answerProfile model.Profile
	answerProfile, code  = model.GetByUserID(userID)
	a.AnswerProfile = answerProfile

	code = a.Create()
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": code.Msg(),
		"data": a,
	})

	//向question的channel发送消息以增加回答数
	logic.UpdateQuestionAnswerCountChan<-a.QuestionID
}

// UpdateAnswer func 更新回答
func UpdateAnswer(c *gin.Context) {

	var a model.Answer
	var code util.MyCode

	var id = c.Param("ans_id")
	a.ID, _ = strconv.ParseUint(id,10,64)
	if err := c.ShouldBind(&a); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    util.AnswerInvalidParams,
			"message": util.AnswerInvalidParams.Msg(),
		})
		return
	}

	// 验证token中的userid和传过来的userid是否一致
	userID := c.MustGet("userID").(uint64)
	// 不一致则返回 AnswerUserIdNotMatch 错误
	if userID != a.AnswerProfileID {
		code = util.AnswerUserIdNotMatch
		c.JSON(
			http.StatusOK, gin.H{
				"code":    code,
				"message": code.Msg(),
			},
		)
		return
	}
	// 更新
	code = a.Update()
	// 更新完成后查询answer并返回
	var updatedAnswer model.Answer
	updatedAnswer.ID = a.ID
	updatedAnswer.Get()

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": code.Msg(),
		"data": updatedAnswer,
	})
}


// DeleteAnswer func
func DeleteAnswer(c *gin.Context) {

	var a model.Answer
	var code util.MyCode

	var id = c.Param("ans_id")
	a.ID, _ = strconv.ParseUint(id,10,64)

	println("request answerId: ",a.ID)

	if err := c.ShouldBind(&a); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    util.AnswerInvalidParams,
			"message": util.AnswerInvalidParams.Msg(),
		})
		return
	}
	// 验证token中的userid和传过来的userid是否一致
	userID := c.MustGet("userID").(uint64)
	// 不一致则返回 AnswerUserIdNotMatch 错误
	if userID != a.AnswerProfileID {
		code = util.AnswerUserIdNotMatch
		c.JSON(
			http.StatusOK, gin.H{
				"code":    code,
				"message": code.Msg(),
			},
		)
		return
	}

	code = a.Delete();

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": code.Msg(),
	})
}


// 根据回答id获取回答详情
func GetAnswer(c *gin.Context) {

	var a model.Answer
	var code util.MyCode
	var err error

	var id = c.Param("ans_id")
	a.ID, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		code = util.AnswerInvalidParams
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": code.Msg(),
		})
		return
	}

	code = a.Get()
	if code != util.CodeSuccess{
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": code.Msg(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": code.Msg(),
		"data": a,
		})
}


// 查询问题下所有回答
func GetAnswers(c *gin.Context) {
	var a model.Answer

	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	questionID, _ := strconv.ParseUint(c.Query("question_id"),10,64)
	a.QuestionID = questionID

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	var alist []model.Answer
	var total int64
	var code util.MyCode

	alist, total, code = a.GetList(pageSize, pageNum)
	if code != util.CodeSuccess{
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": code.Msg(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": code.Msg(),
		"data": gin.H{
			"answerList": alist,
			"total":      total,
		},
	})
}
