package settings

import (
	"gopkg.in/ini.v1"
	"log"
	"strconv"
)

type UserType int

const (
	None UserType = iota
	Student
	Teacher
)

func (u UserType) String() string {
	return strconv.Itoa(int(u))
}

func (u UserType) Int() int {
	return int(u)
}

var Conf *ini.File

type App struct {
	PageSize  int
	JwtSecret string
	Salt1     string
	Salt2     string
}

var AppSettings = &App{}

type Server struct {
	HttpPort string
	RunMode  string
}

var ServerSettings = &Server{}

type DataBase struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

var DataBaseSettings = &DataBase{}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Prefix   string
}

var RedisSettings = &RedisConfig{}

type Geetest struct {
	ID              string
	Key             string
	BypassUrl       string
	CycleTime       int
	BypassStatusKey string
}

var GeetestSettings = &Geetest{}

func init() {
	var err error
	Conf, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.init, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSettings)
	mapTo("server", ServerSettings)
	mapTo("database", DataBaseSettings)
	mapTo("redis", RedisSettings)
	mapTo("geetest", GeetestSettings)
}

func mapTo(section string, v interface{}) {
	err := Conf.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("setting.mapTo %s err: %v", section, err)
	}
}
