package models

import (
	"github.com/astaxie/beego/orm"
)

type Users struct {
	Puid          int `orm:"pk"`
	Platform_code string
	Guid          string
	Original      string
	Wlan          string
	Game_id       int
	User_id       int
	Ctime         string
	Mtime         string
}

func init() {
	orm.RegisterModelWithPrefix("platform_", new(Users))
}
