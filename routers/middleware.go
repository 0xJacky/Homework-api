package routers

import (
	"net/http"

	"github.com/0xJacky/Homework-api/model"
	"github.com/0xJacky/Homework-api/settings"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := model.CurrentUser(c)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Next()
	}
}

func Can(power ...settings.UserType) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*model.User)
		match := false

		for k := range power {
			if power[k].Int() == user.Power {
				match = true
				break
			}
		}
		if !match {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "forbidden",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func SuperUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*model.User)
		if user.SuperUser != 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "forbidden",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
