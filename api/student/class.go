package student

import (
	"github.com/0xJacky/Homework-api/api"
	"github.com/0xJacky/Homework-api/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetClass(c *gin.Context)  {
	id := c.Param("id")
	class, err := model.GetClass(id)
	if err != nil{
		api.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": class,
	})
}

func GetClasses(c *gin.Context)  {
	user := api.CurrentUser(c)

	classes, err := user.GetUserClasses()
	if err != nil {
		api.ErrHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": classes,
	})
}

func JoinClass(c *gin.Context)  {
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
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"mesg": "退出班级成功",
	})
}
