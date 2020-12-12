package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qa/model"
	"qa/util"
	"strconv"
)

// 查看用户的profile
func GetProfile(c *gin.Context)  {
	var p model.Profile
	var u model.User

	var id = c.Param("id")
	u.ID, _ = strconv.ParseUint(id, 10 ,64)

	p,code := model.GetByUserID(u.ID)

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"message" : code.Msg(),
		"data" : p,
	})
}

// 更新用户的profile
func UpdateProfile(c *gin.Context) {
	var p model.Profile

	if err := c.ShouldBind(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    util.UserInvalidParams,
			"message": util.UserInvalidParams.Msg(),
		})
		return
	}

	code := p.UpdateProfile()

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"message" : code.Msg(),
	})
}