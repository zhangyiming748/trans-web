package sql

import (
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func SetMysql(mysql string) {
	var err error
	sqlService := os.Getenv("SQLSERVICE")
	if sqlService == "" {
		log.Fatalf("SQLSERVICE environment variable is not set")
	}
	source := strings.Join([]string{"root:163453@tcp(", sqlService, ")/Translate?charset=utf8"}, "")
	engine, err = xorm.NewEngine("mysql", source)
	if err != nil {
		panic(err)
	}
	engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
}
func GetMysql() *xorm.Engine {
	return engine
}
