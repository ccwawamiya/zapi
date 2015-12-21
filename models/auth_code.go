package models

import (
	"github.com/astaxie/beego/orm"
)

type Auth_code struct {
	Id      int `orm:"pk"`
	Puid    int
	Code    string
	Game_id string
	Times   int
	Ctime   string
	Mtime   string
}

func init() {
	orm.RegisterModelWithPrefix("platform_", new(Auth_code))
}
