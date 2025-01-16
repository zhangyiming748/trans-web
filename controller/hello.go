package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type HelloController struct{}

/*
curl --location --request GET 'http://127.0.0.1:8192/api/v1/s1/gethello?user=<user>' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)'
*/
func (h HelloController) GetHello(ctx *gin.Context) {
	user := ctx.Query("user")
	ctx.String(200, fmt.Sprintf("Hello, %s,this is translater!", user))
}

// 结构体必须大写 否则找不到
type RequestBody struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type ResponseBody struct {
	Name string `json:"name"`
}

/*
 */
func (h HelloController) PostHello(ctx *gin.Context) {
	fmt.Println("get")
	var requestBody RequestBody
	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	} else {
		fmt.Println(requestBody)
	}
	fmt.Println(requestBody.Name, requestBody.Age)
	var rep ResponseBody
	rep.Name = fmt.Sprintf("我已经%d年没见过%s了", requestBody.Age, requestBody.Name)
	ctx.JSON(200, rep)
}
