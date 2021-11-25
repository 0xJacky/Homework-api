package auth

import (
	"github.com/0xJacky/Homework-api/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout(c *gin.Context) {
	currentToken, _ := model.CurrentToken(c)
	_ = model.DeleteJwt(currentToken)
	c.JSON(http.StatusNoContent, gin.H{})
}
