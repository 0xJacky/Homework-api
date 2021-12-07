package model

type Class struct {
	Model
	Name   string `json:"name" binding:"required" gorm:"unique"`
	UserID uint   `json:"user_id"`
	User   *User  `json:"user,omitempty"`
}

func GetClass(id string) (class Class, err error) {
	err = db.Joins("User").First(&class, id).Error
	return
}

func (class *Class) Insert() error {
	err := db.Create(&class).Error
	db.First(class, class.ID)
	return err
}

func (class *Class) Update(n *Class) error {
	err := db.Model(&class).Updates(n).Error
	db.First(class, class.ID)
	return err
}

func (class *Class) IsJoined(u *User) bool {
	var tmp UserClass
	db.Where("user_id", u.ID).Where("class_id", class.ID).First(&tmp)
	// 有关联数据，ID != 0
	return tmp.ID != 0
}
