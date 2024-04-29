package router

import (
	. "backgate/controller"
	"backgate/middleware"
	"backgate/relations"
	"backgate/settings"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var Router *gin.Engine

func InitRouter() *gin.Engine {
	Router = gin.Default()
	const (
		prefix = relations.APISITE_PREFIX
	)

	config := cors.DefaultConfig()
	config.AllowOrigins = settings.Cfg.CorsSettings.Allows
	config.AllowMethods = []string{"OPTIONS", "POST", "GET"}
	config.AddAllowHeaders("Authorization", "Logintype")
	Router.Use(cors.New(config))
	Router.Use(middleware.FilterBan())

	if !settings.Cfg.ReleaseMode {
		Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	// ------------------     公共方法     --------------
	Router.GET("/", CheckStat)
	Router.GET(prefix+"/ping", Ping)

	initBackendRouter(Router)
	return Router
}
