package model

type Assign struct {
	Model
	UserId   uint    `json:"user_id"`
	User     *User   `json:"user"`
	UploadId uint    `json:"upload_id"`
	Upload   *Upload `json:"upload"`
	Grade    uint    `json:"grade"`
}
