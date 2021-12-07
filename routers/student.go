package routers

import (
	"github.com/0xJacky/Homework-api/api/student"
	"github.com/0xJacky/Homework-api/settings"
)

func studentRoute() {
	g := r.Group("/student", AuthRequired(), Can(settings.Student))
	{
		// 班级详情
		g.GET("/class/:id", student.GetClass)
		// 班级列表
		g.GET("/classes", student.GetClasses)
		// 加入班级
		g.POST("/class/:id/join", student.JoinClass)
		// 退出班级
		g.POST("/class/:id/exit", student.ExitClass)

		// 获取作业列表
		g.GET("/class/:id/homeworks", student.GetHomeworks)
		// 作业详情
		g.GET("/homework/:id", student.GetHomework)
	}
}
