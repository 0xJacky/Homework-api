package pkg

import (
	"github.com/0xJacky/Homework-api/settings"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

func GetPage(c *gin.Context) (result int, page int) {
	result = 0
	page = com.StrTo(c.Query("page")).MustInt()
	if page > 1 {
		result = (page - 1) * settings.AppSettings.PageSize
	} else {
		page = 1
	}

	return
}

func TotalPage(total int64) int64 {
	n := total / int64(settings.AppSettings.PageSize)
	if total%int64(settings.AppSettings.PageSize) > 0 {
		n++
	}
	return n
}

func GetPagination(c *gin.Context, total int64, currentPage int) gin.H {
	return gin.H{
		"total":        total,
		"per_page":     settings.AppSettings.PageSize,
		"current_page": currentPage,
		"total_pages":  TotalPage(total),
	}
}
