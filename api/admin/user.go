package admin

import (
	"github.com/0xJacky/Homework-api/api"
	"github.com/0xJacky/Homework-api/model"
	"github.com/0xJacky/Homework-api/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func AddUser(c *gin.Context) {
	var json struct {
		Name     string `json:"name" binding:"required"`
		SchoolID string `json:"school_id" binding:"required"`
		// 隐藏密码
		Password    string `json:"password" binding:"required,min=6"`
		Power       int    `json:"power" binding:"min=1,max=2"`
		SuperUser   int    `json:"super_user" binding:"min=-1,max=1"`
		Gender      int    `json:"gender" binding:"min=0,max=2"`
		Phone       string `json:"phone"`
		Email       string `json:"email" binding:"omitempty,email"`
		Description string `json:"description"`
	}
	if !api.BindAndValid(c, &json) {
		return
	}

	if json.Password != "" {
		json.Password = pkg.PasswordHash(json.Password)
	}

	user := model.User{
		Name:        json.Name,
		SchoolID:    json.SchoolID,
		Password:    json.Password,
		Power:       json.Power,
		SuperUser:   json.SuperUser,
		Gender:      json.Gender,
		Phone:       json.Phone,
		Email:       json.Email,
		Description: json.Description,
	}

	if user.IsConflicted() {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": gin.H{
				"school_id": "登录名已被占用",
			},
			"message": "参数校验失败",
		})
		return
	}

	user.Insert()

	c.JSON(http.StatusOK, user)
}

func EditUser(c *gin.Context) {
	id := c.Param("id")
	user := model.NewUser(id)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}

	var json struct {
		Name     string `json:"name"`
		SchoolID string `json:"school_id"`
		// 隐藏密码
		Password    string `json:"password" binding:"omitempty,min=6"`
		Power       int    `json:"power" binding:"omitempty,min=1,max=2"`
		SuperUser   int    `json:"super_user" binding:"omitempty,min=-1,max=1"`
		Gender      int    `json:"gender" binding:"omitempty,min=0,max=2"`
		Phone       string `json:"phone"`
		Email       string `json:"email" binding:"omitempty,email"`
		Description string `json:"description"`

		LastActive *time.Time `json:"last_active"`
	}
	if !api.BindAndValid(c, &json) {
		return
	}

	if json.Password != "" {
		json.Password = pkg.PasswordHash(json.Password)
	}

	n := model.User{
		Name:        json.Name,
		SchoolID:    json.SchoolID,
		Password:    json.Password,
		Power:       json.Power,
		SuperUser:   json.SuperUser,
		Gender:      json.Gender,
		Phone:       json.Phone,
		Email:       json.Email,
		Description: json.Description,
	}

	user.Updates(&n)

	c.JSON(http.StatusOK, user)
}

func GetUser(c *gin.Context) {
	user := model.NewUser(c.Param("id"))
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetUserList(c *gin.Context) {
	data := model.GetUserList(c, c.Query("name"),
		c.Query("school_id"), c.Query("class_id"))

	c.JSON(http.StatusOK, data)
}

func DeleteUser(c *gin.Context) {
	user := model.NewUser(c.Param("id"))
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}
	user.Delete()
	c.JSON(http.StatusNoContent, nil)
}
