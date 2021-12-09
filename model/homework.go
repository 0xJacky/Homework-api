package model

import (
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"time"
)

type Homework struct {
	Model
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Deadline    time.Time      `json:"deadline"`
	ClassId     uint           `json:"class_id"`
	Class       *Class         `json:"class,omitempty"`
	Uploads     []Upload       `json:"upload,omitempty"`
	Assigns     []Assign       `json:"assign_id,omitempty"`
	Template    datatypes.JSON `json:"template"`
	Answer      datatypes.JSON `json:"-"`
}

type HomeworkWithAnswer struct {
	Model
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Deadline    time.Time      `json:"deadline"`
	ClassId     uint           `json:"class_id"`
	Class       *Class         `json:"class,omitempty"`
	Uploads     []Upload       `json:"upload,omitempty"`
	Assigns     []Assign       `json:"assign_id,omitempty"`
	Template    datatypes.JSON `json:"template"`
	Answer      datatypes.JSON `json:"answer"`
}

func GetHomework(id string) (homework Homework, err error) {
	err = db.Joins("Class").First(&homework, id).Error
	return
}

func TeacherGetHomework(id string) (homework HomeworkWithAnswer, err error) {
	err = db.Model(&Homework{}).
		Joins("Class").First(&homework, id).Error
	return
}

func GetHomeworkList(c *gin.Context, userId, classId, name interface{}) (data *DataList) {
	var h []struct {
		Homework
		Score    uint       `json:"score"`
		AssignAt *time.Time `json:"assign_at"`
	}
	var count int64
	result := db.Select("homeworks.*, score, assigns.assign_at").
		Model(&Homework{}).Joins("left join assigns on homeworks.id=assigns.homework_id"+
		" and assigns.user_id = ?", userId).
		Where("class_id", classId)

	if name != "" {
		result = result.Where("name LIKE ?", "%"+name.(string)+"%")
	}
	result.Count(&count)

	result.Scopes(orderAndPaginate(c)).Find(&h)

	data = GetListWithPagination(&h, c, count)

	return
}

func TeacherGetHomeworkList(c *gin.Context, userId, classId, name interface{}) (data *DataList) {
	var h []struct {
		Homework
		Score uint `json:"score"`
	}
	var count int64
	result := db.
		Model(&Homework{}).Joins("Class").
		Where("class_id", classId).
		Where("user_id", userId)

	if name != "" {
		result = result.Where("name LIKE ?", "%"+name.(string)+"%")
	}
	result.Count(&count)

	result.Scopes(orderAndPaginate(c)).Find(&h)

	data = GetListWithPagination(&h, c, count)

	return
}

func (h *Homework) Insert() (err error) {
	err = db.Create(&h).Error
	db.First(h, h.ID)
	return
}

func (h *Homework) Update(n *Homework) (err error) {
	err = db.Model(h).Updates(n).Error
	db.First(h, h.ID)
	return
}

func (h *Homework) Delete() (err error) {
	err = db.Delete(h).Error
	return
}
