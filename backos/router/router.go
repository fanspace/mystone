package router

import (
	acccon "backos/apps/accounts/controller"
	. "backos/apps/api/controller"
	"backos/middleware"
	"backos/relations"
	"backos/settings"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	const (
		prefix = relations.APISITE_PREFIX
	)

	config := cors.DefaultConfig()
	config.AllowOrigins = settings.Cfg.CorsSettings.Allows
	config.AllowMethods = []string{"OPTIONS", "POST", "GET"}
	config.AddAllowHeaders("Authorization", "Logintype")
	router.Use(cors.New(config))
	router.Use(middleware.FilterBan())

	if !settings.Cfg.ReleaseMode {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	}
	// ------------------     公共方法     --------------
	router.GET("/", CheckStat)
	router.GET(prefix+"/ping", Ping)

	router.POST(prefix+"/login", acccon.Login)
	router.GET(prefix+"/pin/:username", acccon.GenPin)

	return router
}
