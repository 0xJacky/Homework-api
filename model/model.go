package model

import (
	"fmt"
	"github.com/0xJacky/Homework-api/pkg"
	"github.com/0xJacky/Homework-api/rds"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"log"
	"time"

	"errors"
	"github.com/0xJacky/Homework-api/settings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

type Model struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" sql:"index"`
}

func init() {
	dbs := settings.DataBaseSettings

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbs.User, dbs.Password, dbs.Host, dbs.Port, dbs.Name)

	log.Println(dsn)

	var err error

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	// 自动迁移
	AutoMigrate(&User{})
	AutoMigrate(&Upload{})
	AutoMigrate(&Message{})

	// 测试 Redis
	rds.Test()
	// Add First User
	var father User
	err = db.First(&father, "id", 1).Error

	if err == gorm.ErrRecordNotFound {
		log.Println("正在创建初始用户 admin")
		db.Create(&User{
			SchoolID:  "admin",
			Password:  pkg.PasswordHash("123456"),
			SuperUser: 101,
		})
	}

}

func AutoMigrate(obj interface{}) {
	if db.AutoMigrate(obj) != nil {
		panic("auto migrate fail")
	}
}

func orderAndPaginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		sort := c.DefaultQuery("sort", "desc")
		order := c.DefaultQuery("order_by", "id") +
			" " + sort

		page := cast.ToInt(c.Query("page"))
		if page == 0 {
			page = 1
		}
		pageSize := settings.AppSettings.PageSize
		offset := (page - 1) * settings.AppSettings.PageSize

		return db.Order(order).Offset(offset).Limit(pageSize)
	}
}

type Models []Model

type Pagination struct {
	Total       int64 `json:"total"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	TotalPages  int64 `json:"total_pages"`
}

type DataList struct {
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func GetListWithPagination(models interface{},
	c *gin.Context, totalRecords int64) (result *DataList) {

	page := cast.ToInt(c.Query("page"))
	if page == 0 {
		page = 1
	}

	result = &DataList{}

	result.Data = models
	result.Pagination = &Pagination{
		Total:       totalRecords,
		PerPage:     settings.AppSettings.PageSize,
		CurrentPage: page,
		TotalPages:  pkg.TotalPage(totalRecords),
	}

	return
}

// CreateOrUpdate 不存在则创建，存在则更新
func CreateOrUpdate(model interface{}, search map[string]string, obj interface{}) error {
	result := db.Model(model)
	for key, value := range search {
		result = result.Where(key+" = ?", value)
	}
	var tmp interface{}
	result = result.First(tmp)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err := db.Model(model).Create(obj).Error
			if err != nil {
				return err
			}
			return nil
		}
		return result.Error
	}
	err := result.Updates(obj).Error
	if err != nil {
		return err
	}
	return nil
}
