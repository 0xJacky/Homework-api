package global

import (
	"github.com/0xJacky/Homework-api/api"
	"github.com/0xJacky/Homework-api/model"
	"github.com/0xJacky/Homework-api/pkg"
	"github.com/0xJacky/Homework-api/settings"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path"
)

func UserInfo(c *gin.Context) {
	user := api.CurrentUser(c).GetFullUserInfo()
	c.JSON(http.StatusOK, user)
}

func EditUserInfo(c *gin.Context) {
	user := api.CurrentUser(c)

	var json struct {
		Name      string `json:"name" binding:"required"`
		SchoolID  string `json:"school_id" binding:"required,min=2"`
		Gender    int    `json:"gender" binding:"min=1,max=2"`
		Phone     string `json:"phone"`
		Email     string `json:"email" binding:"email"`
	}

	if !api.BindAndValid(c, &json) {
		return
	}

	if user.Power == settings.Student.Int() {
		user.Updates(&model.User{
			Name: json.Name,
			Gender: json.Gender,
			Phone: json.Phone,
			Email: json.Email,
		})
	} else {
		user.Updates(&model.User{
			Name: json.Name,
			SchoolID: json.SchoolID,
			Gender: json.Gender,
			Phone: json.Phone,
			Email: json.Email,
		})
	}

	c.JSON(http.StatusOK, user)
}

func UploadAvatar(c *gin.Context) {
	user := api.CurrentUser(c)
	name := uuid.New().String()

	dst, err := api.UploadSingleFile(c, path.Join("avatar", name))

	if err != nil {
		api.ErrHandler(c, err)
		return
	}

	// 删除旧照片
	if user.Avatar != "" {
		if _, err = os.Stat(user.Avatar); err == nil {
			_ = os.Remove(user.Avatar)
		}
		upload := model.Upload{
			Path: user.Avatar,
		}
		upload.DeleteByPath()
	}
	
	user.Updates(&model.User{
		Avatar: dst,
	})
	
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"url":     dst,
	})
}

func ResetPassword(c *gin.Context) {
	user := api.CurrentUser(c)
	// 避免缓存影响，缓存不缓 password
	u := model.NewUser(user.ID)
	user = &u

	var json struct {
		OrigPassword string `json:"orig_password" binding:"required,min=6"`
		NewPassword  string `json:"new_password" binding:"required,min=6"`
	}

	if !api.BindAndValid(c, &json) {
		return
	}

	// 验证原来的密码
	err := pkg.PasswordVerify(user.Password, json.OrigPassword)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "原密码错误",
		})
		return
	}

	// 改新密码
	user.Updates(&model.User{
		Password:  pkg.PasswordHash(json.NewPassword),
	})
	
	if err != nil {
		api.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
