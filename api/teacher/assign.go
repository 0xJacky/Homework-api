package teacher

import (
	"github.com/0xJacky/Homework-api/api"
	"github.com/0xJacky/Homework-api/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAssignList(c *gin.Context) {
	data := model.TeacherGetAssignList(c, c.Param("id"),
		c.Query("user.school_id"), c.Query("user.name"))

	c.JSON(http.StatusOK, data)
}

func GetAssign(c *gin.Context) {
	a, err := model.FirstAssign(c.Param("id"))
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, a)
}

func EditAssign(c *gin.Context) {
	var json struct {
		Score uint `json:"score" binding:"min=1,max=100"`
	}
	if !api.BindAndValid(c, &json) {
		return
	}
	a, err := model.FirstAssign(c.Param("id"))
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	err = a.Update(&model.Assign{
		Score: json.Score,
	})
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, a)
}
