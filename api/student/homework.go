package student

import (
	"github.com/0xJacky/Homework-api/api"
	"github.com/0xJacky/Homework-api/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetHomework(c *gin.Context) {
	h, err := model.GetHomework(c.Param("id"))
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	assign := model.Assign{
		UserId:     api.CurrentUser(c).ID,
		HomeworkId: h.ID,
	}
	err = model.InitAssign(&assign)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"homework": h,
		"score":    assign.Score,
		"assign":   assign,
	})
}

func GetHomeworks(c *gin.Context) {
	user := api.CurrentUser(c)
	data := model.GetHomeworkList(c, c.Param("id"),
		user.ID,
		c.Query("name"))

	c.JSON(http.StatusOK, data)
}
