package sql

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func SetMysql(mysql string) {
	var err error
	//session := strings.Join([]string{"root:123456@tcp(", mysql, ")/Translate?charset=utf8"}, "")
	session := "root:163453@tcp(192.168.1.9:3306)/Translate?charset=utf8"
	engine, err = xorm.NewEngine("mysql", session)
	if err != nil {
		panic(err)
	}
	engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
}
func GetMysql() *xorm.Engine {
	return engine
}
