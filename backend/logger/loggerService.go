package logger

import (
	"backend/settings"
	"strings"
)

//var log_Rpc_Address = settings.Cfg.RpcSettings.LogRpc
//var grpcLogPool *grpc.ClientConn

func LoggerConfigFromLoggerConfig() *LoggerConfiguration {
	return &LoggerConfiguration{
		EnableConsole: settings.Cfg.LogSettings.EnableConsole,
		ConsoleJson:   *settings.Cfg.LogSettings.ConsoleJson,
		ConsoleLevel:  strings.ToLower(settings.Cfg.LogSettings.ConsoleLevel),
		EnableFile:    settings.Cfg.LogSettings.EnableFile,
		FileJson:      *settings.Cfg.LogSettings.FileJson,
		FileLevel:     strings.ToLower(settings.Cfg.LogSettings.FileLevel),
		FileLocation:  settings.GetLogFileLocation(settings.Cfg.LogSettings.FileLocation),
	}
}

func InitLogger() {
	logconf := NewLogger(LoggerConfigFromLoggerConfig())
	RedirectStdLog(logconf)

	// 使用server logger 作为全局的logger
	InitGlobalLogger(logconf)
}

/*

func SendLogger(logtime string, module string, infotype string,  info string) error  {
	req := new(pb.LoggerInfo)
	if logtime == "" {
		logtime = time.Now().Format("2006-01-02T15:04:05Z")
	}
	if module == "" {
		module = settings.Cfg.AppName
	}
	req.Info = []byte(info)
	req.InfoType = infotype
	req.LoggerTime = logtime
	req.Module = module
	pool, err := getGrpcLogPool()
	if err != nil {
		log.Println(err.Error())
		return err
	} else {
		c := pb.NewLoggerGrpcClient(pool)
		res, err := c.SendLogger(context.Background(), req)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		if !res.Success {
			log.Println(res.Msg)

		}
	}
	return nil
}


func getGrpcLogPool() (*grpc.ClientConn, error) {
	var err error
	if grpcLogPool != nil {
		return grpcLogPool, nil
	}
	grpcLogPool, err = dialgrpc(log_Rpc_Address)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return grpcLogPool, nil
}
*/
/**
* @Description
* @Author  lf
* @Date   5/11/2020 2:36 PM
* @Param
* @Return
* @Exception
*
 */

/*
func dialgrpc(address string) (*grpc.ClientConn,error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return conn, err
}

*/
