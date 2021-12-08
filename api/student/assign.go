package student

import (
	"github.com/0xJacky/Homework-api/api"
	"github.com/0xJacky/Homework-api/model"
	"github.com/0xJacky/Homework-api/pkg"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"os"
	"path"
)

func UploadHomework(c *gin.Context) {
	user := api.CurrentUser(c)
	id := c.Param("id")

	assign := model.Assign{
		UserId:     user.ID,
		Score:      0,
		HomeworkId: pkg.StrToUInt(id),
	}

	err := model.InitAssign(&assign)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	files := form.File["file"]
	// 创建附件存储目录
	p := path.Join("upload", cast.ToString(user.ID), id)
	err = os.MkdirAll(p, 0766)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}

	//循环存文件到服务器本地
	for _, file := range files {
		var filepath string
		filepath = path.Join(p, file.Filename)
		var upload model.Upload
		upload = model.Upload{
			UserId:     user.ID,
			Path:       filepath,
			Size:       uint(file.Size),
			HomeworkId: pkg.StrToUInt(id),
			AssignId:   assign.ID,
		}
		// 将附件存到服务器本地
		err = c.SaveUploadedFile(file, filepath)
		if err != nil {
			api.ErrHandler(c, err)
			return
		}

		err = upload.Save()
		if err != nil {
			api.ErrHandler(c, err)
			return
		}

		assign.Uploads = append(assign.Uploads, upload)
	}
	err = assign.Update(&assign)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"data":   assign,
		"status": "done",
	})
}

func DeleteUpload(c *gin.Context) {
	user := api.CurrentUser(c)
	id := c.Param("upload_id")
	upload, err := model.FirstUpload(id)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}

	if user.ID != upload.UserId {
		c.JSON(http.StatusForbidden, gin.H{
			"code": http.StatusForbidden,
			"mesg": "无权删除该附件",
		})
		return
	}

	upload.DeleteByPath()
	err = os.Remove(upload.Path)
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	assign, err := model.FirstAssign(c.Param("assign_id"))

	if err != nil {
		api.ErrHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, assign.Uploads)
}
