package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"translate/logic"
)

type TranslateController struct{}

// 结构体必须大写 否则找不到
type Translate struct {
	Proxy              string `json:"proxy,omitempty"`
	LocalDeeplxService string `json:"LocalDeeplxService,omitempty"`
}
type TranslateResponseBody struct {
	Proxy string `json:"proxy,omitempty"`
	Msg   string `json:"msg"`
}

/*
curl --location --request POST 'http://127.0.0.1:8193/api/v1/telegram/download' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \

	--data-raw '{
	    "urls": [
	        "string"
	    ],
	    "proxy": "string"
	}'
*/
func (t TranslateController) DownloadAll(ctx *gin.Context) {
	//fmt.Printf("url = %s \nproxy = %s\n", req.URLs, req.Proxy)
	rep := TranslateResponseBody{
		Msg: "已经开始翻译字幕",
	}
	log.Println("开始转换")

	go logic.Start()
	ctx.JSON(200, rep)
}
