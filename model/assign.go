package model

type Assign struct {
	Model
	UserId   	uint    	`json:"user_id"`
	User     	*User   	`json:"user"`
	Uploads   	[]Upload 	`json:"upload"`
	Grade    	uint    	`json:"grade"`
	HomeworkId 	uint 		`json:"homework_id"`
}
