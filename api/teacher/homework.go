package teacher

import (
	"github.com/0xJacky/Homework-api/api"
	"github.com/0xJacky/Homework-api/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddHomework(c *gin.Context)  {
	var homework model.Homework
	if !api.BindAndValid(c, &homework) {
		return
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

func EditHomework(c *gin.Context)  {
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

func DeleteHomework(c *gin.Context)  {
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

func GetHomework(c *gin.Context)  {
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

func GetHomeworks(c *gin.Context)  {

}
