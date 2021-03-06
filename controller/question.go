package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qa/dao"
	"qa/logic"
	"qa/model"
	util "qa/util"
	"strconv"
)

type QuestionVo struct {
	model.GORMBase
	Title          string        `json:"title"`
	Content        string        `json:"content"`
	AnswerCount    int           `json:"answerCount"`
	ViewCount      int           `json:"viewCount"`
	UserID         uint64        `json:"userId,string"`
	CreatorProfile model.Profile `json:"creator"`
	BestAnswer     model.Answer  `json:"bestAnswer"`
}

type QuestionHotListVo struct {
	model.GORMBase
	Title          string        `json:"title"`
	Score 			float64 	`json:"score"`
}

//创建问题
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
	logic.InitQuestionScore(q)
	logic.CreateQuestionViewCountChan <- q.ID
	logic.CreateQuestionAnswerCountChan <- q.ID
}

//更新问题
func UpdateQuestion(c *gin.Context) {
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

	code = q.Update()
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": code.Msg(),
	})
}

//删除问题
func DeleteQuestion(c *gin.Context) {
	var q model.Question
	var code util.MyCode

	if err := c.ShouldBind(&q); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    util.QuestionInvalidParams,
			"message": util.QuestionInvalidParams.Msg(),
		})
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
	code = q.Delete()
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": code.Msg(),
	})
}

//按用户id获取问题列表
func GetAllQuestionByUserId(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	order := c.Query("order")
	//验证token中的userid和传过来的userid是否一致
	userID := c.MustGet("userID").(uint64)

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 5
	}

	if pageNum == 0 {
		pageNum = 1
	}

	switch order {
	case "answercount":
		order = "answer_count desc"
	case "viewcount":
		order = "view_count desc"
	default:
		order = "updated_at desc"
	}

	var qlist []model.Question
	var code util.MyCode
	var total int64

	qlist, total, code = model.GetAllQuestionByUserId(pageSize, pageNum, order, userID)

	var qvolist []QuestionVo
	for _, q := range qlist {
		//a := model.Answer{QuestionID: q.ID}
		//a.Get()
		a:=logic.GetBestAnswer(strconv.FormatUint(q.ID,10))
		var qvo QuestionVo
		util.Copy(&qvo, &q)
		qvo.BestAnswer = a
		qvolist = append(qvolist, qvo)
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

//按问题title模糊搜索问题列表
func GetAllQuestionByTitle(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	order := c.Query("order")
	title := c.Query("title")

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 5
	}

	if pageNum == 0 {
		pageNum = 1
	}

	switch order {
	case "answercount":
		order = "answer_count desc"
	case "viewcount":
		order = "view_count desc"
	default:
		order = "updated_at desc"
	}

	var qlist []model.Question
	var code util.MyCode
	var total int64

	qlist, total, code = model.GetAllQuestionByTitle(pageSize, pageNum, order, title)

	var qvolist []QuestionVo
	for _, q := range qlist {
		//a := model.Answer{QuestionID: q.ID}
		//a.Get()
		a:=logic.GetBestAnswer(strconv.FormatUint(q.ID,10))
		var qvo QuestionVo
		util.Copy(&qvo, &q)
		qvo.BestAnswer = a
		qvolist = append(qvolist, qvo)
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

//查询所有问题
func GetAllQuestion(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	order := c.Query("order")

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 5
	}

	if pageNum == 0 {
		pageNum = 1
	}

	switch order {
	case "answercount":
		order = "answer_count desc"
	case "viewcount":
		order = "view_count desc"
	default:
		order = "updated_at desc"
	}

	var qlist []model.Question
	var code util.MyCode
	var total int64

	qlist, total, code = model.GetAllQuestion(pageSize, pageNum, order)

	var qvolist []QuestionVo
	for _, q := range qlist {
		//a := model.Answer{QuestionID: q.ID}
		//a.Get()
		a:=logic.GetBestAnswer(strconv.FormatUint(q.ID,10))
		var qvo QuestionVo
		util.Copy(&qvo, &q)
		qvo.BestAnswer = a
		qvolist = append(qvolist, qvo)
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

//查询单个问题详情
func GetQuestion(c *gin.Context) {
	qid, _ := strconv.ParseUint(c.Query("id"), 10, 64)

	var q model.Question
	q.ID = qid
	code := q.Get()
	c.JSON(
		http.StatusOK, gin.H{
			"code":    code,
			"message": code.Msg(),
			"data":    q,
		},
	)
	if code == util.CodeSuccess {
		logic.UpdateQuestionViewCountChan <- q.ID
	}
}

//查询问题热榜
func GetQuestionHotList(c *gin.Context) {
	qlist := logic.QuestionHotList()

	if qlist != nil {
		var qvolist []QuestionHotListVo
		for _, q := range qlist {
			var qvo QuestionHotListVo
			util.Copy(&qvo, &q)
			qvo.Score=dao.RDB.ZScore(logic.ZSetKey,strconv.FormatUint(q.ID,10)).Val()
			qvolist = append(qvolist, qvo)
		}

		c.JSON(
			http.StatusOK, gin.H{
				"code":    util.CodeSuccess,
				"message": util.CodeSuccess.Msg(),
				"data":    qvolist,
			},
		)
	} else {
		c.JSON(
			http.StatusOK, gin.H{
				"code":    util.CodeError,
				"message": util.CodeError.Msg(),
				"data":    qlist,
			},
		)
	}
}

func GetAnswerListByScore(c *gin.Context){
	pageSize, _ := strconv.ParseInt(c.Query("pagesize"),10,64)
	pageNum, _ := strconv.ParseInt(c.Query("pagenum"),10,64)
	questionID:= c.Query("question_id")

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

	alist= logic.GetAnswerListByScore(questionID,pageNum, pageSize)
	c.JSON(http.StatusOK, gin.H{
		"code":    util.CodeSuccess,
		"message": util.CodeSuccess.Msg(),
		"data": gin.H{
			"answerList": alist,
		},
	})
}

func GetVoteInfo(c *gin.Context){
	aid:= c.Query("answerId")
	userID := c.MustGet("userID").(uint64)
	direction:=logic.GetVoteInfo(aid,userID)

	c.JSON(http.StatusOK, gin.H{
		"code":    util.CodeSuccess,
		"message": util.CodeSuccess.Msg(),
		"data":direction,
	})
}

func GetAnswerListByUserId(c *gin.Context){
	pageSize, _ := strconv.ParseInt(c.Query("pagesize"),10,64)
	pageNum, _ := strconv.ParseInt(c.Query("pagenum"),10,64)
	userID := c.MustGet("userID").(uint64)

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	alist,total,code := model.GetAnswerListByUserId(userID,pageSize,pageNum)
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": code.Msg(),
		"total":total,
		"data": gin.H{
			"answerList": alist,
		},
	})
}