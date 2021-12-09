package routers

import (
	"github.com/0xJacky/Homework-api/api/admin"
	"github.com/0xJacky/Homework-api/api/global"
)

func adminRoute() {
	g := r.Group("/admin", AuthRequired(), SuperUser())
	{

		// 服务器状态统计
		g.GET("analytic", admin.Analytic)
		// 消息测试
		g.POST("/live", global.SendTestMessage)
		// 用户
		g.GET("users", admin.GetUserList)
		g.GET("user/:id", admin.GetUser)
		g.POST("user/:id", admin.EditUser)
		g.POST("user", admin.AddUser)
		g.DELETE("user/:id", admin.DeleteUser)

		// 数据导入
		d := g.Group("data_import")
		{
			d.StaticFile("/student_template", "./static/学生导入模板.xlsx")
			d.StaticFile("/teacher_template", "./static/教师导入模板.xlsx")
			d.POST("/parse_student_excel", admin.DataImportParseStudentExcel)
			d.POST("/import_excel_student", admin.ImportExcelStudent)
			d.POST("/parse_teacher_excel", admin.DataImportParseTeacherExcel)
			d.POST("/import_excel_teacher", admin.ImportExcelTeacher)
		}
	}
}
