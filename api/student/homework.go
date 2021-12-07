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
	var assign *model.Assign
	assign, err = model.InitAssign(&model.Assign{
		UserId:     api.CurrentUser(c).ID,
		HomeworkId: h.ID,
	})
	c.JSON(http.StatusOK, gin.H{
		"homework": h,
		"score":    assign.Score,
	})
}

func GetHomeworks(c *gin.Context) {
	user := api.CurrentUser(c)
	data := model.GetHomeworkList(c, c.Param("id"),
		user.ID,
		c.Query("name"))

	c.JSON(http.StatusOK, data)
}

func AssignHomework(c *gin.Context) {

}
