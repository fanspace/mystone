package router

import (
	"backgate/controller"
	"backgate/relations"
	"github.com/gin-gonic/gin"
)

func initBackendRouter(router *gin.Engine) {
	prefix := relations.APISITE_PREFIX
	router.POST(prefix+"/login", controller.Login)
	router.GET(prefix+"/pin/:username", controller.GenPin)
}
