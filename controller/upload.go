package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qa/model"
)

func UpLoad(c *gin.Context) {
	file, fileHeader, _ := c.Request.FormFile("file")

	fileSize := fileHeader.Size

	url, code := model.UpLoadFile(file, fileSize)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": code.Msg(),
		"url":     url,
	})

}