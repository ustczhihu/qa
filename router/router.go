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
		question.GET("/queryAll", controller.GetAllQuestion)
	}
	return
}
