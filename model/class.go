package model

import "log"

type Class struct {
	Model
	Name   string `json:"name" binding:"required" gorm:"unique"`
	UserID uint   `json:"user_id"`
	User   *User  `json:"user,omitempty"`
}

func GetClass(id string) (class Class, err error) {
	err = db.First(&class, id).Error
	return
}

func (class *Class) Insert() error {
	err := db.Create(&class).Error
	return err
}

func (class *Class) Update(n *Class) error {
	err := db.Model(&class).Updates(n).Error
	return err
}
