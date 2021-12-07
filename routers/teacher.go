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
		// 退出班级
		g.POST("/class/:id/exit", teacher.ExitClass)

		// 发布作业
		g.POST("/homework", teacher.AddHomework)
		// 修改作业
		g.PUT("/homework", teacher.EditHomework)
		// 删除作业
		g.DELETE("/homework/:id", teacher.DeleteHomework)
		// 作业详情
		g.GET("/homework/:id", teacher.GetHomework)
		// 作业列表
		g.GET("/homework", teacher.GetHomeworks)
	}
}