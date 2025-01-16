package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"translate/bootstrap"
	"translate/model"
	"translate/sql"
	"translate/util"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func testResponse(c *gin.Context) {
	c.JSON(http.StatusGatewayTimeout, gin.H{
		"code": http.StatusGatewayTimeout,
		"msg":  "timeout",
	})
}

func timeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(3000*time.Millisecond),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(testResponse),
	)
}
func init() {
	util.SetLog("gin.log")
	os.Setenv("MYSQL", "192.168.1.9:3306")
	if mysql := os.Getenv("MYSQL"); mysql != "" {
		sql.SetMysql(mysql)
	} else {
		log.Fatalf("must set env var: MYSQL")
	}
	err := sql.GetMysql().Ping()
	if err != nil {
		log.Fatalf("ping mysql err: %v", err)
	}
	err = sql.GetMysql().Sync2(new(model.TranslateCache))
	if err != nil {
		log.Printf("sync translate history err: %v", err)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
func main() {
	// gin服务
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	engine.Use(timeoutMiddleware())
	bootstrap.InitHello(engine)
	bootstrap.InitTranslate(engine)
	// 启动http服务
	engine.Run(":9004")
}
