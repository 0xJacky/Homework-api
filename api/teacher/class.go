package teacher

import (
	"github.com/0xJacky/Homework-api/api"
	"github.com/0xJacky/Homework-api/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddClass(c *gin.Context) {
	var json struct {
		Name string `json:"name" binding:"required"`
	}
	if !api.BindAndValid(c, &json) {
		return
	}
	user := api.CurrentUser(c)
	class := &model.Class{
		Name:   json.Name,
		UserID: user.ID,
	}

	err := class.Insert()

	if err != nil {
		api.ErrHandler(c, err)
		return
	}

	userClass := model.UserClass{
		UserID:  user.ID,
		ClassID: class.ID,
	}

	err = userClass.Save()

	if err != nil {
		api.ErrHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, class)
}

func EditClass(c *gin.Context) {
	var code int
	var mesg string
	id := c.Param("id")
	var json struct {
		Name string `json:"name"`
	}
	if !api.BindAndValid(c, &json) {
		return
	}
	class, err := model.GetClass(id)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}

	n := model.Class{
		Name: json.Name,
	}
	err = class.Update(&n)
	if err != nil {
		code = http.StatusForbidden
		mesg = "更新班级信息失败: 已存在同名班级"
	} else {
		code = http.StatusOK
		mesg = "更新班级信息成功"
	}
	c.JSON(code, gin.H{
		"code": code,
		"mesg": mesg,
		"data": class,
	})
}

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

	data := model.TeacherGetClasses(c, user.ID, c.Query("name"))

	c.JSON(http.StatusOK, data)
}

func DeleteClass(c *gin.Context) {

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
	user := api.CurrentUser(c)
	id := c.Param("id")
	class, err := model.GetClass(id)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	// 自己创建的班级不能退出
	if class.UserID == user.ID {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "自己创建的班级不能退出",
		})
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
