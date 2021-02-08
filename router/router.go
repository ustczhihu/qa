package router

import (
	"github.com/gin-gonic/gin"
	"qa/config"
	"qa/controller"
	"qa/middleware"
)

func Init() (r *gin.Engine) {

	gin.SetMode(config.Conf.Mode)
	r = gin.New()
	r.Use(middleware.Log())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	user := r.Group("/user")
	{
		user.GET("/validate", controller.RegisterValidate)
		user.POST("/register", controller.Register)
		user.POST("/login", controller.Login)
	}

	question := r.Group("/question")
	{
		question.POST("/add", middleware.JwtToken(), controller.AddQuestion)
		question.POST("/update", middleware.JwtToken(), controller.UpdateQuestion)
		question.POST("/delete", middleware.JwtToken(), controller.DeleteQuestion)
		question.GET("/queryAllByUserId", middleware.JwtToken(), controller.GetAllQuestionByUserId)
		question.GET("/queryAllByTitle", middleware.JwtToken(), controller.GetAllQuestionByTitle)
		question.GET("/queryHotList",controller.GetQuestionHotList)
		question.GET("/queryAll", controller.GetAllQuestion)
		question.GET("/get", controller.GetQuestion)
		question.GET("/queryAnswerListByScore", controller.GetAnswerListByScore)
		question.GET("/queryVoteInfo",  middleware.JwtToken(),controller.GetVoteInfo)
		question.GET("/queryAnswerListByUserId",  middleware.JwtToken(),controller.GetAnswerListByUserId)
	}

	profile := r.Group("/profile")
	profile.Use(middleware.JwtToken())
	{
		profile.GET("/getByUserID/:id", controller.GetProfile)
		profile.POST("/updateProfile", controller.UpdateProfile)
		profile.POST("/uploadAvatarUrl", controller.UpLoad)
	}

	answer := r.Group("/answer")
	{
		answer.POST("/add", middleware.JwtToken(), controller.AddAnswer)
		answer.GET("/", controller.GetAnswers)
		answer.GET("/:ans_id", controller.GetAnswer)
		answer.PUT("/:ans_id", middleware.JwtToken(), controller.UpdateAnswer)
		answer.DELETE("/:ans_id", middleware.JwtToken(), controller.DeleteAnswer)
	}

	vote := r.Group("/vote")
	{
		vote.POST("/update", middleware.JwtToken(), controller.VoteForAnswer)
	}
	return
}
