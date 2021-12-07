package model

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/0xJacky/Homework-api/rds"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cast"
	"log"
	"time"

	"github.com/0xJacky/Homework-api/settings"
	"gorm.io/gorm/clause"

	"github.com/gin-gonic/gin"
)

type User struct {
	Model

	Name     	string 		`json:"name"`
	SchoolID 	string 		`json:"school_id" gorm:"index:school_id,unique"`
	// 隐藏密码
	Password    string 		`json:"-"`
	Power       int    		`json:"power"`
	SuperUser   int    		`json:"super_user"`
	Gender      int    		`json:"gender"`
	Phone       string 		`json:"phone"`
	Email       string 		`json:"email"`
	Description string 		`json:"description" gorm:"type:longtext"`
	Avatar      string 		`json:"avatar"`
	Classes		[]Class 	`json:"classes" gorm:"many2many:user_classes;"`
	Assigns		[]Assign	`json:"assigns"`

	LastActive *time.Time `json:"last_active" gorm:"default:NULL"`
}

func (u *User) GetUserClasses() (classes []Class, err error) {
	err = db.Model(u).Association("Classes").Find(&classes)
	return
}

// JoinClass 加入班级
func (u *User) JoinClass(class Class) (err error) {
	err = db.Model(u).Association("Classes").Append(&class)
	return
}

// ExitClass 退出班级
func (u *User) ExitClass(class Class) (err error) {
	err = db.Model(u).Association("Classes").Delete(class)
	return
}

func (u *User) GetFullUserInfo() (user User) {

	db.First(&user, u.ID)
	return
}

func NewUser(conds ...interface{}) (user User) {
	db.Preload(clause.Associations).First(&user, conds)
	return
}

func GetUserList(c *gin.Context, name, schoolId, classId interface{}) (data *DataList) {
	var users []User
	var total int64

	result := db.Model(&User{})

	if name != "" {
		result = result.Where("name = ?", "%"+cast.ToString(name)+"%")
	}

	if schoolId != "" {
		result = result.Where("schoolId = ?", schoolId)
	}

	if classId != "" {
		result = result.Where("class_id = ?", classId)
	}

	result.Count(&total)

	result.Scopes(orderAndPaginate(c)).Find(&users)

	data = GetListWithPagination(&users, c, total)

	return
}

func (u *User) Insert() {
	db.Create(u)
	db.Preload(clause.Associations).First(u)
}

func (u *User) IsConflicted() bool {
	// 创建
	if u.ID == 0 {
		return db.Where("school_id", u.SchoolID).
			First(&User{}).RowsAffected > 0
	}
	// 修改
	var tmp User
	db.Where("school_id", u.SchoolID).
		First(&tmp)
	// 记录不存在
	if tmp.ID == 0 {
		return false
	}
	// 记录存在，但是查到的 ID 与请求的 ID 不同，说明冲突
	return tmp.ID != u.ID
}

func (u *User) Updates(user *User) {
	tmp := *u
	db.First(&tmp, "id", u.ID).Updates(*user)
	db.Preload(clause.Associations).First(u, u.ID)
	// 清缓存
	_ = rds.Del("user:" + cast.ToString(u.ID) + ":info")

	return
}

func (u *User) UpdatesWithoutPreload(n *User) {

	db.Model(&User{}).Where("id", u.ID).Updates(n)
	// 清缓存
	_ = rds.Del("user:" + cast.ToString(u.ID) + ":info")

	return
}

func (u *User) Save() {
	db.Where("id", u.ID).Updates(u)
	// 清缓存
	_ = rds.Del("user:" + cast.ToString(u.ID) + ":info")
}

func (u *User) Delete() {
	db.Delete(u)
}

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	SchoolID string `json:"school_id"`
	Power    int    `json:"power"`
	jwt.StandardClaims
}

func UpdateUserTime(user *User) error {
	return rds.Set("lastActive:user:"+cast.ToString(user.ID), time.Now(), 0)
}

func FindUser(SchoolID string) (user User, err error) {
	result := db.Where(&User{SchoolID: SchoolID}).First(&user)
	err = result.Error
	if err != nil {
		log.Println("用户不存在")
	}
	return
}

func CurrentToken(c *gin.Context) (token string, err error) {
	// 先验证 token 是否合法
	if len(c.Request.Header["Token"]) == 0 {
		if c.Query("token") == "" {
			return "", errors.New("no token")
		}
		tmp, _ := base64.StdEncoding.DecodeString(c.Query("token"))
		token = string(tmp)
	} else {
		token = c.Request.Header["Token"][0]
	}
	return
}

// CurrentUser 用户访问初始化
func CurrentUser(c *gin.Context) (user User, err error) {
	var token string
	token, err = CurrentToken(c)
	if err != nil {
		return user, err
	}

	var userId uint
	userId, err = ValidateJWT(token)
	if err != nil {
		return user, err
	}

	// 检查 redis 缓存
	var cache string
	key := "user:" + cast.ToString(userId) + ":info"
	cache, _ = rds.Get(key)

	// 拿不到缓存或无法反序列化，从数据库里取数据
	if cache == "" || json.Unmarshal([]byte(cache), &user) != nil {
		// 根据 ID 查找用户
		result := db.First(&user, userId)

		if result.Error != nil {
			return user, err
		}
		// 再次尝试缓存进 redis

		bytes, _ := json.Marshal(user)

		// 不捕获错误，有问题大不了下次继续查数据库，缓存1分钟
		_ = rds.Set(key, string(bytes), 1*time.Minute)
	}

	// 不管从哪里取数据，都要更新上次访问的时间
	err = UpdateUserTime(&user)

	if err != nil {
		return user, err
	}

	log.Println("[当前用户]", user.Name)

	return
}

func GenerateJWT(userID uint, schoolID string, userPower int) (signedToken string, err error) {
	claims := JWTClaims{
		UserID:   userID,
		SchoolID: schoolID,
		Power:    userPower,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * 24 * time.Hour).Unix(),
			Issuer:    "0xJacky",
		},
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = unsignedToken.SignedString([]byte(settings.AppSettings.JwtSecret))

	return
}

func (u *User) SaveJwt(jwt string) error {
	return rds.Set("jwt:"+jwt, u.ID, 15*24*time.Hour)
}

func DeleteJwt(jwt string) error {
	return rds.Del("jwt:" + jwt)
}

func ValidateJWT(jwt string) (userId uint, err error) {
	var uid string
	uid, err = rds.Get("jwt:" + jwt)
	userId = cast.ToUint(uid)

	return
}

func EditUsers(ids []int, user User) error {
	// 清缓存
	for i := range ids {
		_ = rds.Del("user:" + cast.ToString(ids[i]) + ":info")
	}
	return db.Where("id IN ?", ids).Updates(user).Error
}
