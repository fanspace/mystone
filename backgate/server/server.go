package server

import (
	log "backgate/logger"
	"backgate/service"
	"backgate/settings"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type GracefulServer struct {
	Server *http.Server
	Router *gin.Engine
	//SqlSupplier *SqlSupplier
	Log *log.Logger
}

var Srv *GracefulServer

func LoggerConfigFromLoggerConfig() *log.LoggerConfiguration {
	return &log.LoggerConfiguration{
		EnableConsole: settings.Cfg.LogSettings.EnableConsole,
		ConsoleJson:   *settings.Cfg.LogSettings.ConsoleJson,
		ConsoleLevel:  strings.ToLower(settings.Cfg.LogSettings.ConsoleLevel),
		EnableFile:    settings.Cfg.LogSettings.EnableFile,
		FileJson:      *settings.Cfg.LogSettings.FileJson,
		FileLevel:     strings.ToLower(settings.Cfg.LogSettings.FileLevel),
		FileLocation:  settings.GetLogFileLocation(settings.Cfg.LogSettings.FileLocation),
	}
}
func NewServer() *GracefulServer {
	log.Info("project professional And tecnical personnel api gate server")
	Srv = &GracefulServer{}
	Srv.Log = log.NewLogger(LoggerConfigFromLoggerConfig())
	log.RedirectStdLog(Srv.Log)

	// 使用server logger 作为全局的logger
	log.InitGlobalLogger(Srv.Log)
	// 开启bigcache
	service.InitCache()
	service.InitGrpcs()
	service.InitApis()
	return Srv
}

func (gracefulServer *GracefulServer) Start() error {
	log.Info("server start")
	gracefulServer.Server = &http.Server{
		Addr:    settings.Cfg.ServiceSettings.ListenAddress,
		Handler: gracefulServer.Router,

		//ReadTimeout:  time.Duration(Cfg.ServiceSettings.ReadTimeout) * time.Second,
		//	WriteTimeout: time.Duration(Cfg.ServiceSettings.WriteTimeout) * time.Second,
	}
	log.Info("Start Listening and serving HTTP on " + settings.Cfg.ServiceSettings.ListenAddress)
	//log.Info(fmt.Sprintf("Current version is %v (%v/%v)", CurrentVersion, BuildDate, BuildHash))

	return nil
}

func (gracefulServer *GracefulServer) ShutDown() {

	go func() {
		// service connections
		if err := Srv.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Info("listen: %s\n")
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutdown Server ...")
	service.CloseCache()
	log.Info("Destroy BigCache ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := Srv.Server.Shutdown(ctx); err != nil {
		log.Error("Server Shutdown err:" + err.Error())
	}
	log.Info("Server exiting")
}
