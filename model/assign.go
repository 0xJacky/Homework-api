package model

import (
	"github.com/0xJacky/Homework-api/settings"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"time"
)

type Assign struct {
	Model
	UserId         uint           `json:"user_id"`
	User           *User          `json:"user,omitempty"`
	Uploads        []Upload       `json:"uploads,omitempty"`
	ObjectiveScore uint           `json:"objective_score"`
	Score          uint           `json:"score"`
	HomeworkId     uint           `json:"homework_id"`
	Homework       *Homework      `json:"homework,omitempty"`
	Answer         datatypes.JSON `json:"answer"`
	AssignAt       *time.Time     `json:"assign_at"`
}

func InitAssign(n *Assign) error {

	err := db.Where("user_id = ? AND homework_id = ?", n.UserId, n.HomeworkId).Preload("Uploads").FirstOrCreate(n).Error

	return err
}

func (a *Assign) Update(n *Assign) (err error) {
	err = db.Model(a).Updates(n).Error
	db.Preload("Uploads").First(a, a.ID)
	return
}

func FirstAssign(conds ...interface{}) (a Assign, err error) {
	err = db.Joins("User").Preload("Uploads").
		Joins("Homework").
		First(&a, conds...).Error
	return
}

func TeacherGetAssignList(c *gin.Context, homeworkId, schoolId, studentName interface{}) (data *DataList) {
	var assigns []struct {
		Assign
		Deadline time.Time `json:"deadline"`
	}
	var count int64
	result := db.Model(&UserClass{}).Select("assigns.*, deadline").Joins("User").
		Joins("left join assigns on assigns.user_id=user_classes.user_id "+
			"and assigns.homework_id = ?", homeworkId).
		Joins("join homeworks on homeworks.id=?", homeworkId).
		Where("user_classes.class_id = homeworks.class_id").
		Where("power", settings.Student)

	if schoolId != "" {
		result = result.Where("users.school_id LIKE ?", "%"+schoolId.(string)+"%")
	}

	if studentName != "" {
		result = result.Where("users.name LIKE ?", "%"+studentName.(string)+"%")
	}

	result.Scopes(orderAndPaginate(c)).Find(&assigns)
	data = GetListWithPagination(&assigns, c, count)

	return
}
