package model

import "time"

type Homework struct {
	Model
	Name        string    	`json:"name" binding:"required"`
	Description string    	`json:"description"`
	Deadline    time.Time 	`json:"deadline"`
	ClassId		uint 		`json:"class_id"`
	Uploads     []Upload    `json:"upload"`
	Assigns 	[]Assign	`json:"assign_id"`
}

func GetHomework(id string) (homework Homework, err error) {
	err = db.First(&homework, id).Error
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
