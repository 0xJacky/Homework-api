package auth

import (
	"errors"
	"github.com/0xJacky/Homework-api/api"
	"github.com/0xJacky/Homework-api/model"
	"github.com/0xJacky/Homework-api/pkg"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Login(c *gin.Context) {
	var json struct {
		SchoolID         string `json:"school_id" binding:"min=2"`
		Password         string `json:"password" binding:"min=6"`
		GeetestChallenge string `json:"geetest_challenge"`
		GeetestValidate  string `json:"geetest_validate"`
		GeetestSeccode   string `json:"geetest_seccode"`
	}

	if !api.BindAndValid(c, &json) {
		return
	}

	user, err := model.FindUser(json.SchoolID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "用户不存在或密码错误",
				"code":    4031,
			})
			return
		}
		api.ErrHandler(c, err)
		return
	}


	if err = pkg.PasswordVerify(user.Password, json.Password); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "用户不存在或密码错误",
			"code":    4031,
		})
		return
	}

	log.Println("---> user login success: ", user.Name)
	if err = model.UpdateUserTime(&user); err != nil {
		api.ErrHandler(c, err)
		return
	}
	var token string
	// generate jwt, and save jwt to database
	token, err = model.GenerateJWT(user.ID, user.SchoolID, user.Power)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	// sava jwt to redis
	err = user.SaveJwt(token)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	// api return
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"token":   token,
	})

}
