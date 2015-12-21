package models

import (
	"github.com/astaxie/beego/orm"
)

type Game struct {
	Game_id         int `orm:"pk"`
	Uid             int
	Game_key        string
	Name            string
	Update_time     int
	Security_key    string
	Notify_url      string
	Version_control int
	On_line_status  int
	Show_status     int
}

func init() {
	orm.RegisterModelWithPrefix("g5_", new(Game))
}
