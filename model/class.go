package model

type Class struct {
	Model
	Name 		string 			`json:"name" binding:"required" gorm:"unique"`
	Users		[]User 			`json:"users" gorm:"many2many:user_classes;"`
	Homeworks	[]Homework 		`json:"homeworks"`
}

func GetClass(id string) (class Class, err error){
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
