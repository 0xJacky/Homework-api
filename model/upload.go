// Package model 对应开发文档 6.1.12 上传文件
package model

type Upload struct {
	Model

	UserId 		uint   	`json:"user_id"`
	Path   		string 	`json:"path"`
	Size   		uint   	`json:"size"`

	HomeworkId 	uint 	`json:"homework_id" gorm:"default:NULL"`
	AssignId 	uint 	`json:"assign_id" gorm:"default:NULL"`
}

func FirstUpload(conds ...interface{}) (u Upload, err error) {
	err = db.First(&u, conds...).Error
	return
}

// UpdateUploadPath 修改表中存储的相对路径
func UpdateUploadPath(oldPath string, newPath string) (err error) {
	err = db.Model(&Upload{}).Where("path = ?", oldPath).Update("path", newPath).Error
	return
}

func (u *Upload) Save() (err error) {
	err = db.Create(u).Error
	db.First(u, u.ID)
	return
}

func (u *Upload) Updates(n *Upload) {
	db.Model(&Upload{}).Where("id", u.ID).Updates(n)

	db.First(u, u.ID)
}

func (u *Upload) DeleteByPath() {
	db.Unscoped().Model(&Upload{}).Where("path", u.Path).Delete(u)
}
