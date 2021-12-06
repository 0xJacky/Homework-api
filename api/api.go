package api

import (
	"bytes"
	"encoding/json"
	"github.com/0xJacky/Homework-api/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"regexp"
)

type JsonSnakeCase struct {
	Value interface{}
}

func (c JsonSnakeCase) MarshalJSON() ([]byte, error) {
	// Regexp definitions
	var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
	var wordBarrierRegex = regexp.MustCompile(`(\w)([A-Z])`)
	marshalled, err := json.Marshal(c.Value)
	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			return bytes.ToLower(wordBarrierRegex.ReplaceAll(
				match,
				[]byte(`${1}_${2}`),
			))
		},
	)
	return converted, err
}

func ErrHandler(c *gin.Context, err error) {
	log.Println(err.Error())
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": err.Error(),
	})
}

func currentUser(c *gin.Context) *model.User {
	return c.MustGet("user").(*model.User)
}

func CurrentUser(c *gin.Context) *model.User {
	return currentUser(c)
}

func BindAndValid(c *gin.Context, v interface{}) bool {
	return bindAndValid(c, v)
}
