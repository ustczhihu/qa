package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"qa/model"
	util "qa/util"
	"strconv"
)

type Creator struct {
	ID        uint64 `json:"userId"`
	Nickname  string `json:"nickname"`
	AvatarUrl string `json:"avatarUrl"`
}

type QuestionVo struct {
	ID        uint64  `json:"id"`
	CreatedAt string  `json:"createAt"`
	UpdatedAt string  `json:"updateAt"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	Creator   Creator `json:"creator"`
}

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

	code := q.Create()
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

	var qvolist []QuestionVo
	var p model.Profile
	for _, v := range qlist {
		var qvo QuestionVo
		qvo.ID=v.ID
		qvo.Title=v.Title
		qvo.Content=v.Content
		qvo.CreatedAt=v.CreatedAt
		qvo.UpdatedAt=v.UpdatedAt

		p, _ = model.GetByUserID(v.UserID)

		fmt.Printf("%#v",p)
		qvo.Creator.ID=v.UserID
		qvo.Creator.Nickname=p.Nickname
		qvo.Creator.AvatarUrl=p.AvatarUrl
		qvolist=append(qvolist, qvo)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": code.Msg(),
		"data": gin.H{
			"questionList": qvolist,
			"total":        total,
		},
	})
}
