package teacher

import (
	"github.com/0xJacky/Homework-api/api"
	"github.com/0xJacky/Homework-api/model"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"net/http"
	"time"
)

func AddHomework(c *gin.Context) {
	var json struct {
		Name        string         `json:"name" binding:"required"`
		Description string         `json:"description"`
		Deadline    time.Time      `json:"deadline" binding:"required"`
		ClassId     uint           `json:"class_id" binding:"required"`
		Template    datatypes.JSON `json:"template"`
	}

	if !api.BindAndValid(c, &json) {
		return
	}

	homework := model.Homework{
		Name:        json.Name,
		Description: json.Description,
		Deadline:    json.Deadline,
		ClassId:     json.ClassId,
		Template:    json.Template,
	}

	err := homework.Insert()
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"mesg": "发布作业成功",
		"data": homework,
	})
}

func EditHomework(c *gin.Context) {
	id := c.Param("id")
	homework, err := model.GetHomework(id)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	var newHomework model.Homework
	if !api.BindAndValid(c, &newHomework) {
		return
	}
	err = homework.Update(&newHomework)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"mesg": "修改作业成功",
		"data": homework,
	})
}

func DeleteHomework(c *gin.Context) {
	id := c.Param("id")
	homework, err := model.GetHomework(id)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	err = homework.Delete()
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"mesg": "删除作业成功",
	})
}

func GetHomework(c *gin.Context) {
	id := c.Param("id")
	homework, err := model.GetHomework(id)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": homework,
	})
}

func GetHomeworks(c *gin.Context) {
	user := api.CurrentUser(c)
	data := model.TeacherGetHomeworkList(c, user.ID,
		c.Param("id"), c.Query("name"))

	c.JSON(http.StatusOK, data)
}
