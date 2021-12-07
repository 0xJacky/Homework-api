package model

type UserClass struct {
	Model
	UserID  uint   `json:"user_id"`
	User    *User  `json:"user,omitempty"`
	ClassID uint   `json:"class_id"`
	Class   *Class `json:"class,omitempty"`
}
