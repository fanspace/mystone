package router

import (
	"backgate/controller"
	"backgate/middleware"
	"backgate/rbac"
	"backgate/relations"
	"github.com/gin-gonic/gin"
)

func initBackendRouter(router *gin.Engine) {
	prefix := relations.APISITE_PREFIX
	router.POST(prefix+"/login", controller.Login)
	router.GET(prefix+"/pin/:username", controller.GenPin)

	// 只验证登录
	authGroup := router.Group(prefix + "/auth")
	authGroup.Use(middleware.MustLogin())
	{
		authGroup.GET("/menu/mine", controller.QueryMyMenu)
	}

	resGroup := router.Group(prefix + "/res")
	resGroup.Use(middleware.MustLogin(), middleware.MustAuthorizer(rbac.Casbin))
	{
		resGroup.POST("/menu/list", controller.QueryAllMenus)
		resGroup.POST("/menu/query", controller.FetchMenu)
	}
}
