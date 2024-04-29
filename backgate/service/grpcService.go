package service

import (
	log "backgate/logger"
	"backgate/relations"
	"backgate/settings"
	"errors"
	"fmt"
	"google.golang.org/grpc"
)

/**
* @Description
* @Author  lf
* @Date   5/11/2020 2:36 PM
* @Param
* @Return
* @Exception
*
 */
func dialgrpc(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return conn, err
}

var GrpcPools map[string]*grpc.ClientConn

func InitGrpcs() {
	if settings.Cfg.GrpcSettings != nil && len(settings.Cfg.GrpcSettings) > 0 {
		GrpcPools = make(map[string]*grpc.ClientConn)
		for grpcname, gprcaddr := range settings.Cfg.GrpcSettings {
			fmt.Println(grpcname, gprcaddr)
			conn, err := dialgrpc(gprcaddr)
			if err != nil {
				log.Error(err.Error())
			}
			GrpcPools[grpcname] = conn
		}
	}
}

func GetGrpcPool(grpcname string) (*grpc.ClientConn, error) {
	/*var conn *grpc.ClientConn
	if !settings.Cfg.DaprMode {
		conn = GrpcPools[grpcname]
	} else {
		conn = GrpcPools["daprrpc"]
	}
	*/
	if settings.Cfg.DaprMode {
		grpcname = "daprrpc"
	}

	//if conn == nil {
	if addr, ok := settings.Cfg.GrpcSettings[grpcname]; ok {
		return dialgrpc(addr)
	} else {
		return nil, errors.New(relations.CUS_ERR_4004)
	}

	//}
	//return conn, nil
}
