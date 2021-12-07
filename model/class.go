package model

import "log"

type Class struct {
	Model
	Name 		string `json:"name" binding:"required" gorm:"unique"`
	Users		[]User `json:"users" gorm:"many2many:user_classes;"`
}

func GetClass(id string) (class Class, err error){
	result := db.First(&class, id)
	err = result.Error
	if err != nil {
		log.Println("找不到该班级")
	}
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
