package main

import (
	"backgate/db"
	"backgate/rabbitmq"
	"backgate/rbac"
	"backgate/router"
	"backgate/server"
	"backgate/settings"
	"github.com/prometheus/common/log"
	"os"
)

func main() {
	settings.InitCfg()
	db.InitRed()
	defer db.CloseRed()
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
