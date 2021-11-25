package routers

import (
	"github.com/0xJacky/Homework-api/api/auth"
	"github.com/0xJacky/Homework-api/api/global"
	"github.com/0xJacky/Homework-api/live"
)

func globalRoute() {
	r.POST("/login", auth.Login)
	r.DELETE("/logout", AuthRequired(), auth.Logout)

	// 全局统一注册 ws
	r.GET("live", live.WsHandler)
	// 复用接口
	g := r.Group("/", AuthRequired())
	{
		// 消息列表
		g.GET("messages", global.GetMessageList)
		// 消息详情
		g.GET("message/:id", global.ReadMessage)
		// 删除消息
		g.DELETE("message/:id", global.DeleteMessage)
		// 删除全部消息
		g.DELETE("messages", global.DeleteAllMessage)
		// 获取当前用户
		g.GET("user", global.UserInfo)
		// 修改当前用户信息
		g.POST("user", global.EditUserInfo)
		// 重设当前用户的密码
		g.POST("reset_password", global.ResetPassword)
		// 上传头像
		g.POST("user/avatar", global.UploadAvatar)
	}
}
