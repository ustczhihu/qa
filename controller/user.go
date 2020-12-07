package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qa/middleware"
	"qa/model"
	"qa/util"
)

//用户注册检验
func RegisterValidate(c *gin.Context) {
	var u model.User

	if err := c.ShouldBind(&u); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    util.UserInvalidParams,
			"message": util.UserInvalidParams.Msg(),
		})
		return
	}

	_, code := u.Get()
	if code == util.UserNotExist {
		c.JSON(http.StatusOK, gin.H{
			"code":    util.CodeSuccess,
			"message": util.CodeSuccess.Msg(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    util.UserExist,
			"message": util.UserExist.Msg(),
		})
	}
}

//用户注册
func Register(c *gin.Context) {
	var u model.User
	var p model.Profile

	if err := c.ShouldBind(&u); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    util.UserInvalidParams,
			"message": util.UserInvalidParams.Msg(),
		})
		return
	}

	code := u.Create()
	if code != util.CodeSuccess{
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": code.Msg(),
		})
		return
	}

	p.UserID=u.ID
	code = p.Create()
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": code.Msg(),
	})

}

//用户登录
func Login(c *gin.Context) {
	var u model.User

	if err := c.ShouldBind(&u); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    util.UserInvalidParams,
			"message": util.UserInvalidParams.Msg(),
		})
		return
	}

	var token string
	var code util.MyCode
	var user model.User
	user, code = u.CheckLogin()

	if code == util.CodeSuccess {
		token, code = middleware.SetToken(u.Username)
	}else{
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": code.Msg(),
		})
		return
	}

	var p model.Profile
	p,code=model.GetByUserID(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": code.Msg(),
		"token":   token,
		"data": p,
	})
}
