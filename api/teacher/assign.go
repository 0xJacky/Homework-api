package teacher

import (
	"github.com/0xJacky/Homework-api/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAssignList(c *gin.Context) {
	data := model.TeacherGetAssignList(c, c.Param("id"), c.Query("user.name"))

	c.JSON(http.StatusOK, data)
}
