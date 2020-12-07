package router

import (
	"github.com/gin-gonic/gin"
	"qa/controller"
	"qa/middleware"
)

func Init() (r *gin.Engine) {
	gin.DisableConsoleColor()
	r = gin.Default()

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
