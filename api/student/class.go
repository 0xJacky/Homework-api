package student

import (
	"github.com/0xJacky/Homework-api/api"
	"github.com/0xJacky/Homework-api/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetClass(c *gin.Context) {
	user := api.CurrentUser(c)
	id := c.Param("id")
	class, err := model.GetClass(id)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": class,
		"join": class.IsJoined(user),
	})
}

func GetClasses(c *gin.Context) {
	user := api.CurrentUser(c)

	data := user.GetUserClasses(c)

	c.JSON(http.StatusOK, data)
}

func JoinClass(c *gin.Context) {
	user := api.CurrentUser(c)
	id := c.Param("id")
	class, err := model.GetClass(id)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	err = user.JoinClass(class)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"mesg": "加入班级成功",
	})
}

func ExitClass(c *gin.Context) {
	var code int
	var mesg string
	user := api.CurrentUser(c)
	id := c.Param("id")
	class, err := model.GetClass(id)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	err = user.ExitClass(class)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	code = http.StatusOK
	mesg = "退出班级成功"
	c.JSON(code, gin.H{
		"code": code,
		"mesg": mesg,
	})
}
