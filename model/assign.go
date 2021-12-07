package model

type Assign struct {
	Model
	UserId     uint     `json:"user_id"`
	User       *User    `json:"user"`
	Uploads    []Upload `json:"upload"`
	Score      uint     `json:"score"`
	HomeworkId uint     `json:"homework_id"`
}

func InitAssign(n *Assign) (*Assign, error) {

	err := db.FirstOrCreate(n).Error

	return n, err
}