package global

import (
	"github.com/0xJacky/Homework-api/api"
	"github.com/0xJacky/Homework-api/live"
	"github.com/0xJacky/Homework-api/model"
	"github.com/0xJacky/Homework-api/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendTestMessage(c *gin.Context) {
	var json struct {
		UserId  uint   `json:"user_id" binding:"required"`
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Times   int    `json:"times" binding:"min=1,max=10"`
	}

	if !api.BindAndValid(c, &json) {
		return
	}

	for i := 0; i < json.Times; i++ {
		live.SetMessage(json.UserId, gin.H{
			"data": gin.H{
				"title":   json.Title,
				"content": json.Content,
			},
		})
	}
	c.JSON(http.StatusOK, json)
}

func GetMessageList(c *gin.Context) {

	user := api.CurrentUser(c)

	pageOffset, currentPage := pkg.GetPage(c)

	messages, total := model.GetMessageList(user.ID, pageOffset)

	c.JSON(http.StatusOK, gin.H{
		"data":       messages,
		"pagination": pkg.GetPagination(c, total, currentPage),
	})
}

func ReadMessage(c *gin.Context) {
	user := api.CurrentUser(c)
	m := model.ReadMessage(user.ID, c.Param("id"))
	c.JSON(http.StatusOK, m)
}

func DeleteMessage(c *gin.Context) {
	user := api.CurrentUser(c)
	model.DeleteMessage(user.ID, c.Param("id"))
	c.JSON(http.StatusNoContent, gin.H{})
}

func DeleteAllMessage(c *gin.Context) {
	user := api.CurrentUser(c)
	model.DeleteAllMessages(user.ID)
	c.JSON(http.StatusNoContent, gin.H{})
}
