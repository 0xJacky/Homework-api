package model

import "github.com/gin-gonic/gin"

type Assign struct {
	Model
	UserId     uint     `json:"user_id"`
	User       *User    `json:"user,omitempty"`
	Uploads    []Upload `json:"upload,omitempty"`
	Score      uint     `json:"score"`
	HomeworkId uint     `json:"homework_id"`
}

func InitAssign(n *Assign) error {

	err := db.FirstOrCreate(n).Error

	return err
}

func (a *Assign) Update(n *Assign) (err error) {
	err = db.Model(a).Updates(n).Error
	db.Preload("Uploads").First(a, a.ID)
	return
}

func TeacherGetAssignList(c *gin.Context, homeworkId, studentName interface{}) (data *DataList) {
	var assigns []struct {
		Assign
	}
	var count int64
	result := db.Model(&UserClass{}).Select("assigns.*").Joins("User").
		Joins("left join assigns on assigns.user_id=user_classes.user_id").
		Joins("join homeworks on homeworks.id=homework_id").
		Where("user_classes.class_id = homeworks.class_id").
		Where("homework_id", homeworkId)

	if studentName != "" {
		result = result.Where("users.name LIKE ?", "%"+studentName.(string)+"%")
	}

	result.Scopes(orderAndPaginate(c)).Find(&assigns)
	data = GetListWithPagination(&assigns, c, count)
	return
}
