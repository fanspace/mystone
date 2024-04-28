package main

import (
	"backos/db"
	"backos/rabbitmq"

	log "backos/logger"
	"backos/rbac"
	"backos/router"
	"backos/server"
	"backos/settings"
	"os"
)

// @title Enlist Apply Api Gate
// @version 1.0
// @description 考试报名 Apply Api.

// @contact.email lf6128@163.com

// @license.name GPL v3
// @license.url http://www.gnu.org/licenses/quick-guide-gplv3.html

// @host 192.168.0.36:10001
// @BasePath /mgmt

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// mysentinel.InitSentinel()
	settings.InitCfg()
	db.InitDb()
	defer db.CloseDb()
	rbac.InitCasbin()
	router := router.InitRouter()
	srv := server.NewServer()
	defer srv.ShutDown()
	go rabbitmq.InitRabConsumer()
	srv.Router = router

	err := srv.Start()
	if err != nil {
		log.Error(err.Error())
		log.Info("Server has stopped")
		os.Exit(0)
	}
}
