package sql

import (
	"log"
	"testing"
	"time"
	"xorm.io/xorm"
)

// go test -v -run TestCreate

func TestCreate(t *testing.T) {
	session := "root:163453@tcp(192.168.1.9:3306)/Translate?charset=utf8"
	engine, err := xorm.NewEngine("mysql", session)
	if err != nil {
		panic(err)
	}
	engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	err = engine.Ping()
	if err != nil {
		log.Fatalf("ping mysql err: %v", err)
	}
}
