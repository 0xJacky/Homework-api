package routers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		c.Next()
	}
}

func InitRouter() *gin.Engine {
	r = gin.New()
	r.Use(gin.Logger())

	r.Use(recovery())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "page not found",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	globalRoute()

	studentRoute()

	teacherRoute()

	adminRoute()

	return r
}
