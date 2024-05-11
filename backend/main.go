package main

import (
	"backend/apps"
	"backend/db"
	log "backend/logger"
	"backend/settings"
	pb "backend/training"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	settings.InitCfg()
	db.InitDb()
	defer db.CloseDb()
	log.InitLogger()
	var BackendGrpcService = apps.BackendGrpcService{}
	listen, err := net.Listen("tcp", "0.0.0.0"+settings.Cfg.ServiceSettings.ListenAddress)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	//s := grpc.NewServer()
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(apps.TokenInterceptor()),
	}
	s := grpc.NewServer(serverOptions...)
	pb.RegisterBackendGrpcServer(s, BackendGrpcService)
	log.Info("Account Grpc Server Starting on " + settings.Cfg.ServiceSettings.ListenAddress)

	s.Serve(listen)
}
