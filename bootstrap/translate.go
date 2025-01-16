package bootstrap

import (
	"translate/controller"

	"github.com/gin-gonic/gin"
)

func InitTranslate(engine *gin.Engine) {
	routeGroup := engine.Group("/api")
	{
		c := new(controller.TranslateController)
		//routeGroup.GET("/v1/s1/gethello", c.GetHello)
		routeGroup.POST("/v1/translate", c.DownloadAll)
	}
}
