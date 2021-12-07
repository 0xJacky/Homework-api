package model

import "time"

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

func GetHomework(id string) (homework Homework, err error) {
	err = db.First(&homework, id).Error
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
