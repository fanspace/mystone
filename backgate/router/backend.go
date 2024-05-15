package router

import (
	"backgate/controller"
	"backgate/middleware"
	"backgate/rbac"
	"backgate/relations"
	"backgate/settings"
	"github.com/gin-gonic/gin"
)

func initBackendRouter(router *gin.Engine) {
	prefix := relations.APISITE_PREFIX

	if !settings.Cfg.ReleaseMode {
		router.GET(prefix+"/apis", controller.ListApis)
	}

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
		menuMgrGroup.GET("/list", controller.QueryAllMenus)
		menuMgrGroup.POST("/query", controller.FetchMenu)
		menuMgrGroup.POST("/add", controller.AddMenu)
		menuMgrGroup.POST("/update", controller.UpdateMenu)
		menuMgrGroup.POST("/del", controller.DelMenu)
	}

	apiMgrGroup := router.Group(prefix + "/apiMgr")
	apiMgrGroup.Use(middleware.MustLogin(), middleware.MustAuthorizer(rbac.Casbin))
	{
		apiMgrGroup.POST("/list", controller.ListApis)
		apiMgrGroup.GET("/list/:pid", controller.ListApiByPid)
		apiMgrGroup.POST("/add", controller.AddApi)
		apiMgrGroup.POST("/update", controller.UpdateApi)
		apiMgrGroup.POST("/del", controller.DelApi)
	}

	dictMgrGroup := router.Group(prefix + "/dictMgr")
	dictMgrGroup.Use(middleware.MustLogin(), middleware.MustAuthorizer(rbac.Casbin))
	{
		//dictMgrGroup.POST("/list", controller.ListDicts)
		dictMgrGroup.POST("/add", controller.AddDict)
		dictMgrGroup.POST("/update", controller.UpdateDict)
		dictMgrGroup.POST("/del", controller.DelDict)
	}

}
