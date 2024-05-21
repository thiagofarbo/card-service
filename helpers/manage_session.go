package helpers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"jwt-authentication/models"
)

func SetSession(c *gin.Context, user *models.User) {
	session := sessions.Default(c)
	session.Set("userID", user.ID)
	err := session.Save()
	if err != nil {
		return
	}
}
