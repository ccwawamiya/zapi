package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	mysqlsetting := beego.AppConfig.String("mysqluser") + ":" +
		beego.AppConfig.String("mysqlpass") + "@tcp(" +
		beego.AppConfig.String("mysqlurls") + ")/" +
		beego.AppConfig.String("mysqldb") +
		"?charset=utf8"
	orm.RegisterDataBase("default", "mysql", mysqlsetting)

}
