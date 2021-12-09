package admin

import (
	"encoding/json"
	"fmt"
	"github.com/0xJacky/Homework-api/api"
	"github.com/0xJacky/Homework-api/model"
	"github.com/0xJacky/Homework-api/pkg"
	"github.com/0xJacky/Homework-api/settings"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
)

type DataImportTool struct {
	collegeMap     map[string]uint
	majorMap       map[string]uint
	classMap       map[uint]map[string]uint
	UnknownExcel   bool
	UnknownCollege bool
	UnknownMajor   bool
	Students       []Student
	Teachers       []Teacher
}

// Student 学号,姓名,性别
type Student struct {
	UserId   uint   `json:"user_id"`
	SchoolId string `json:"school_id"`
	Name     string `json:"name"`
	Gender   int    `json:"gender"`
}

// Teacher 工号,姓名,学院
type Teacher struct {
	UserId   uint   `json:"user_id"`
	SchoolId string `json:"school_id"`
	Name     string `json:"name"`
}

func NewDataImportTool() *DataImportTool {
	return &DataImportTool{
		collegeMap: make(map[string]uint),
		majorMap:   make(map[string]uint),
		classMap:   make(map[uint]map[string]uint),
	}
}

func (t *DataImportTool) TestStudentExcel(file io.Reader) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		fmt.Println(err)
		t.UnknownExcel = true
		return
	}
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		t.UnknownExcel = true
		return
	}
	genderT := map[string]int{
		"男": 2,
		"女": 1,
	}

	// 检查格式是否正确
	if len(rows) < 2 {
		t.UnknownExcel = true
		return
	}

	rule := []string{"学号", "姓名", "性别"}
	row := rows[0]
	if len(row) != len(rule) {
		t.UnknownExcel = true
		return
	}
	for i := range rule {
		if row[i] != rule[i] {
			t.UnknownExcel = true
			return
		}
	}

	rows = rows[1:]

	for _, row := range rows {
		user, _ := model.FindUser(row[0])
		t.Students = append(t.Students, Student{
			UserId:   user.ID,
			SchoolId: row[0],
			Name:     row[1],
			Gender:   genderT[row[2]],
		})

		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}

	data, _ := json.MarshalIndent(t.Students, "", "\t")

	fmt.Println(string(data))
}

func DataImportParseStudentExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	tool := NewDataImportTool()
	excel, err := file.Open()
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	tool.TestStudentExcel(excel)
	c.JSON(http.StatusOK, gin.H{
		"data":            tool.Students,
		"unknown_college": tool.UnknownCollege,
		"unknown_major":   tool.UnknownMajor,
		"unknown_excel":   tool.UnknownExcel,
	})
}

func ImportExcelStudent(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	tool := NewDataImportTool()
	excel, err := file.Open()
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	tool.TestStudentExcel(excel)
	if tool.UnknownExcel {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "表格不合法",
		})
		return
	}
	if tool.UnknownMajor || tool.UnknownCollege {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "存在未创建学院或专业",
		})
		return
	}

	for i := range tool.Students {
		s := tool.Students[i]
		u, err := model.FindUser(s.SchoolId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				schoolId := s.SchoolId

				user := model.User{
					Name:      s.Name,
					SchoolID:  schoolId,
					Password:  pkg.PasswordHash(schoolId),
					SuperUser: 0,
					Power:     settings.Student.Int(),
					Gender:    s.Gender,
				}
				user.Insert()
			} else {
				log.Println(err)
			}
		} else {
			u.Updates(&model.User{
				Name:      s.Name,
				SuperUser: 0,
				Power:     settings.Student.Int(),
				Gender:    s.Gender,
			})
		}
	}

	excel, _ = file.Open()
	tool = NewDataImportTool()
	tool.TestStudentExcel(excel)

	c.JSON(http.StatusOK, gin.H{
		"data":            tool.Students,
		"unknown_college": tool.UnknownCollege,
		"unknown_major":   tool.UnknownMajor,
		"unknown_excel":   tool.UnknownExcel,
	})
}

func (t *DataImportTool) TestTeacherExcel(file io.Reader) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		fmt.Println(err)
		t.UnknownExcel = true
		return
	}
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		t.UnknownExcel = true
		return
	}

	// 检查格式是否正确
	if len(rows) < 2 {
		t.UnknownExcel = true
		return
	}

	rule := []string{"工号", "姓名"}
	row := rows[0]
	if len(row) != len(rule) {
		t.UnknownExcel = true
		return
	}
	for i := range rule {
		if row[i] != rule[i] {
			t.UnknownExcel = true
			return
		}
	}

	rows = rows[1:]

	for _, row := range rows {
		user, _ := model.FindUser(row[0])
		t.Teachers = append(t.Teachers, Teacher{
			UserId:   user.ID,
			SchoolId: row[0],
			Name:     row[1],
		})

		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}

	data, _ := json.MarshalIndent(t.Teachers, "", "\t")

	fmt.Println(string(data))
}

func DataImportParseTeacherExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	tool := NewDataImportTool()
	excel, err := file.Open()
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	tool.TestTeacherExcel(excel)
	c.JSON(http.StatusOK, gin.H{
		"data":            tool.Teachers,
		"unknown_college": tool.UnknownCollege,
		"unknown_major":   tool.UnknownMajor,
		"unknown_excel":   tool.UnknownExcel,
	})
}

func ImportExcelTeacher(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	tool := NewDataImportTool()
	excel, err := file.Open()
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	tool.TestTeacherExcel(excel)
	if tool.UnknownExcel {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "表格不合法",
		})
		return
	}
	if tool.UnknownMajor {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "存在未创建学院",
		})
		return
	}

	for i := range tool.Teachers {
		s := tool.Teachers[i]
		_, err = model.FindUser(s.SchoolId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				schoolId := s.SchoolId
				user := model.User{
					Name:      schoolId,
					SchoolID:  schoolId,
					Password:  pkg.PasswordHash(schoolId),
					SuperUser: 0,
					Power:     settings.Teacher.Int(),
				}
				user.Insert()
			} else {
				log.Println(err)
			}
		}
	}

	excel, _ = file.Open()
	tool = NewDataImportTool()
	tool.TestTeacherExcel(excel)

	c.JSON(http.StatusOK, gin.H{
		"data":            tool.Teachers,
		"unknown_college": tool.UnknownCollege,
		"unknown_major":   tool.UnknownMajor,
		"unknown_excel":   tool.UnknownExcel,
	})
}
