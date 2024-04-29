package service

import (
	log "backgate/logger"
	"backgate/relations"
	"backgate/settings"
	pb "backgate/training"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"reflect"
	"time"
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
	conn := GrpcPools[grpcname]
	if conn == nil || conn.GetState().String() != "READY" {
		if settings.Cfg.DaprMode {
			grpcname = "daprrpc"
		}
		if addr, ok := settings.Cfg.GrpcSettings[grpcname]; ok {
			return dialgrpc(addr)
		} else {
			return nil, errors.New(relations.CUS_ERR_4007)
		}

	}
	return conn, nil
}
func DealGrpcCall[T any](req T, methodName string, grpcName string) (any, error) {
	var res []reflect.Value
	pool, err := GetGrpcPool(grpcName)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	//defer pool.Close()
	if pool == nil {
		log.Error(fmt.Sprintf("connect to %s failed", grpcName))
		return nil, errors.New(relations.CUS_ERR_4007)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var c any
	switch grpcName {
	case "backendrpc":
		c = pb.NewBackendGrpcClient(pool)

	default:
		return nil, errors.New(relations.CUS_ERR_4002)
	}
	value := reflect.ValueOf(c)
	f := value.MethodByName(methodName)
	var parms []reflect.Value
	parms = []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)}
	res = f.Call(parms)
	if len(res) != 2 {
		return nil, errors.New(relations.CUS_ERR_4007)
	}
	if res[1].Interface() != nil {
		err = res[1].Interface().(error)
		return nil, err
	}
	if res[0].Interface() != nil {
		return res[0].Interface(), nil
	}
	return nil, err
}
