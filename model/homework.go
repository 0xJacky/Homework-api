package model

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Homework struct {
	Model
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	ClassId     uint      `json:"class_id"`
	Class       *Class    `json:"class"`
	Uploads     []Upload  `json:"upload,omitempty"`
	Assigns     []Assign  `json:"assign_id,omitempty"`
}

func GetHomework(conds ...interface{}) (homework Homework, err error) {
	err = db.First(&homework, conds...).Error
	return
}

func GetHomeworkList(c *gin.Context, userId, classId, name interface{}) (data *DataList) {
	var h []struct {
		Homework
		Score uint `json:"score"`
	}
	var count int64
	result := db.Select("homeworks.*, score").
		Model(&Homework{}).Joins("join assigns on homeworks.id=assigns.homework_id").
		Where("class_id", classId).
		Where("assigns.user_id", userId)

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
	return
}

func (h *Homework) Update(n *Homework) (err error) {
	err = db.Model(h).Updates(n).Error
	return
}
