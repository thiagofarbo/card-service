package controllers

import (
	"github.com/gin-gonic/gin"
	"jwt-authentication/models"
)

func Signup(c *gin.Context) {
	var dataForm models.UserRequest
	c.Bind(&dataForm)
}
