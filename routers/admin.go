package routers

import (
	"github.com/0xJacky/Homework-api/api/admin"
	"github.com/0xJacky/Homework-api/api/global"
)

func adminRoute() {
	g := r.Group("/admin", AuthRequired(), SuperUser())
	{

		// 服务器状态统计
		// g.GET("analytic", admin.Analytic)
		// 消息测试
		g.POST("/live", global.SendTestMessage)
		// 用户
		g.GET("users", admin.GetUserList)
		g.GET("user/:id", admin.GetUser)
		g.POST("user/:id", admin.EditUser)
		g.POST("user", admin.AddUser)
		g.DELETE("user/:id", admin.DeleteUser)
	}
}
