package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"log"
	"net/http"
	"strings"
)

type ValidError struct {
	Key     string
	Message string
}

func bindAndValid(c *gin.Context, v interface{}) bool {
	errs := make(map[string]string)
	err := c.ShouldBindJSON(v)
	if err != nil {
		log.Println(err)
		uni := ut.New(zh.New())
		trans, _ := uni.GetTranslator("zh")
		v, ok := binding.Validator.Engine().(*val.Validate)

		if ok {
			_ = zhTranslations.RegisterDefaultTranslations(v, trans)
		}

		verrs, ok := err.(val.ValidationErrors)
		if !ok {
			log.Println(verrs)
			c.JSON(http.StatusNotAcceptable, gin.H{
				"message": "请求参数错误",
				"code":    http.StatusNotAcceptable,
			})
			return false
		}

		for key, value := range verrs.Translate(trans) {
			errs[key[strings.Index(key, ".")+1:]] = value
		}

		c.JSON(http.StatusNotAcceptable, gin.H{
			"errors":  JsonSnakeCase{errs},
			"message": "请求参数错误",
			"code":    http.StatusNotAcceptable,
		})

		return false
	}

	return true
}
