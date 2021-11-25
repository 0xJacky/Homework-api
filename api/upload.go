// Package api 对应开发文档6.1.12上传文件
package api

import (
	"github.com/0xJacky/Homework-api/model"
	"github.com/0xJacky/Homework-api/pkg"
	"github.com/gin-gonic/gin"
	"log"
	"path"
	"path/filepath"
)

type UploadForm struct {
	UserId   string `json:"user_id" valid:"Required;"`
	Path     string `json:"path" valid:"Required;"`
	Url      string `json:"url" valid:"Required;"`
	NoticeId int    `json:"notice_id" valid:"Required;"`
}

// UploadSingleFile 通用文件上传函数，请勿直接在路由中调用
func UploadSingleFile(c *gin.Context, filePath string) (dst string, err error) {
	user := currentUser(c)
	file, err := c.FormFile("file")
	if err != nil {
		ErrHandler(c, err)
		return
	}
	// file.Filename SHOULD NOT be trusted.
	log.Println(file.Filename)
	ext := filepath.Ext(file.Filename)
	dst = path.Join("upload", filePath+ext)
	pkg.ExistsOrCreate(path.Dir(dst))
	// save file to specific dst
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		ErrHandler(c, err)
		return
	}
	// 数据库中记录数据
	upload := model.Upload{
		UserId: user.ID,
		Path:   dst,
		Size:   uint(file.Size),
	}

	err = model.CreateOrUpdate(&model.Upload{},
		map[string]string{
			"path": dst,
		}, &upload)

	if err != nil {
		ErrHandler(c, err)
		return
	}

	return
}
