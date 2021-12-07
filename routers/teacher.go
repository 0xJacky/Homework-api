package routers

import (
	"github.com/0xJacky/Homework-api/api/teacher"
	"github.com/0xJacky/Homework-api/settings"
)

func teacherRoute()  {
	g := r.Group("/teacher", AuthRequired(), Can(settings.Teacher))
	{
		// 创建班级
		g.POST("/class", teacher.AddClass)
		// 班级详情
		g.GET("/class/:id", teacher.GetClass)
		// 修改班级信息
		g.PUT("/class/:id", teacher.EditClass)
		// 删除班级
		// g.DELETE("/class/:id", teacher.DeleteClass)
		// 班级列表
		g.GET("/classes", teacher.GetClasses)
		// 加入班级
		g.POST("/class/:id/join", teacher.JoinClass)
	}
}