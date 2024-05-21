package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"jwt-authentication/db"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionID := session.Get("userId")
		if session == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		userId := sessionID.(uint64)

		user := db.UserFind(userId)

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("userID", user)

		c.Next()
	}
}
