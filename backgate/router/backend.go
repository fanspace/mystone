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

	// 只验证登录,不验证权限
	authGroup := router.Group(prefix + "/auth")
	authGroup.Use(middleware.MustLogin())
	{
		authGroup.GET("/menu/mine", controller.QueryMyMenu)
	}

	menuMgrGroup := router.Group(prefix + "/menuMgr")
	menuMgrGroup.Use(middleware.MustLogin(), middleware.MustAuthorizer(rbac.Casbin))
	{
		menuMgrGroup.POST("/list", controller.QueryAllMenus)
		menuMgrGroup.POST("/query", controller.FetchMenu)
		menuMgrGroup.POST("/add", controller.AddMenu)
		menuMgrGroup.POST("/update", controller.UpdateMenu)
		menuMgrGroup.POST("/del", controller.DelMenu)
	}
}
