package api

import (
	"bytes"
	"encoding/json"
	"github.com/0xJacky/Homework-api/model"
	"github.com/beego/beego/v2/adapter/validation"
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

func GetAndValidateForm(c *gin.Context, api_obj interface{}) {
	var err error
	var re bool
	valid := validation.Validation{}
	err = c.BindJSON(api_obj)
	if err != nil {
		ErrHandler(c, err)
		panic("json bind failed")
	}
	re, err = valid.Valid(api_obj)
	if err != nil {
		ErrHandler(c, err)
		panic("valid init failed")
	}
	if !re {
		validError := make(map[string]string)

		for _, v := range valid.Errors {
			validError[v.Field] = v.Message
			log.Println("validError", v.Field, v.Message)
		}

		c.JSON(http.StatusBadRequest, JsonSnakeCase{gin.H{
			"errors":  validError,
			"message": "请求参数错误",
			"code":    http.StatusBadRequest,
		}})
		panic("validate failed")
	}
}
