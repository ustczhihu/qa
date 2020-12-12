package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qa/model"
	util "qa/util"
	"strconv"
)

// 创建问题
func AddQuestion(c *gin.Context) {
	var q model.Question

	if err := c.ShouldBind(&q); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    util.QuestionInvalidParams,
			"message": util.QuestionInvalidParams.Msg(),
		})
		return
	}

	msg, code := util.Validate(&q)
	if code != util.CodeSuccess {
		c.JSON(
			http.StatusOK, gin.H{
				"code":    code,
				"message": msg,
			},
		)
		return
	}

	//验证token中的userid和传过来的userid是否一致
	userID := c.MustGet("userID").(uint64)

	if userID != q.UserID {
		code = util.QuestionUserIdNotMatch
		c.JSON(
			http.StatusOK, gin.H{
				"code":    code,
				"message": code.Msg(),
			},
		)
		return
	}

	code = model.CheckQuestion(q.Title)
	if code != util.QuestionNotExist {
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": code.Msg(),
		})
		return
	}

	code = q.Create()
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": code.Msg(),
	})
}

// 查询所有问题
func GetAllQuestion(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	var qlist []model.Question
	var code util.MyCode
	var total int64
	qlist, total, code = model.GetAllQuestion(pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": code.Msg(),
		"data": gin.H{
			"questionList": qlist,
			"total":        total,
		},
	})
}
