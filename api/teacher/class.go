package teacher

import (
	"github.com/0xJacky/Homework-api/api"
	"github.com/0xJacky/Homework-api/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddClass(c *gin.Context)  {
	var code int
	var mesg string
	var class model.Class
	if !api.BindAndValid(c, &class){
		return
	}
	err := class.Insert()
	if err != nil{
		code = http.StatusForbidden
		mesg = "创建班级失败: 已存在同名班级"
	}else{
		code = http.StatusOK
		mesg = "创建班级成功"
	}
	c.JSON(code, gin.H{
		"code": code,
		"mesg": mesg,
		"data": class,
	})
}

func EditClass(c *gin.Context)  {
	var code int
	var mesg string
	id := c.Param("id")
	var json struct{
		Name string `json:"name"`
	}
	if !api.BindAndValid(c, &json){
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
	if err != nil{
		code = http.StatusForbidden
		mesg = "更新班级信息失败: 已存在同名班级"
	}else{
		code = http.StatusOK
		mesg = "更新班级信息成功"
	}
	c.JSON(code, gin.H{
		"code": code,
		"mesg": mesg,
		"data": class,
	})
}

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

func DeleteClass(c *gin.Context)  {

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
