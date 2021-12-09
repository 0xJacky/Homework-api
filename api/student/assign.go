package student

import (
	"github.com/0xJacky/Homework-api/api"
	"github.com/0xJacky/Homework-api/model"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/datatypes"
	"net/http"
	"time"
)

func AssignHomework(c *gin.Context) {
	var json struct {
		Answer datatypes.JSON `json:"answer" binding:"required"`
	}
	if !api.BindAndValid(c, &json) {
		return
	}
	user := api.CurrentUser(c)
	assign := model.Assign{
		UserId:     user.ID,
		HomeworkId: cast.ToUint(c.Param("id")),
	}
	err := model.InitAssign(&assign)

	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	t := time.Now()
	err = assign.Update(&model.Assign{
		Answer:   json.Answer,
		AssignAt: &t,
	})
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	var score uint
	score, err = assign.Judge()
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	err = assign.Update(&model.Assign{
		ObjectiveScore: score,
	})
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, assign)
}
