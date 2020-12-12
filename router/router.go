package router

import (
	"github.com/gin-gonic/gin"
	"qa/config"
	"qa/controller"
	"qa/middleware"
)

func Init() (r *gin.Engine) {

	gin.SetMode(config.Conf.Mode)
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

	profile := r.Group("/profile")
	profile.Use(middleware.JwtToken())
	{
		profile.GET("/getByUserID/:id", controller.GetProfile)
		profile.POST("/updateProfile", controller.UpdateProfile)
		profile.POST("/uploadAvatarUrl", controller.UpLoad)
	}
	return
}
