package apps

import (
	accsrv "backend/apps/account/service"
	"backend/apps/res/service"
	log "backend/logger"
	"backend/utils"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)
import pb "backend/training"

type BackendGrpcService struct{}

// 登录
func (ags BackendGrpcService) HandleAccountLogin(ctx context.Context, in *pb.AccountLoginReq) (*pb.AccountLoginRes, error) {
	res := new(pb.AccountLoginRes)
	token, usid, err := accsrv.UserLogin(in)
	if err != nil {
		return nil, err
	}
	fmt.Println(token, usid)
	res.Token = token
	res.UserInfo = &pb.UserInfo{
		UserId:    usid,
		Username:  "Administrator",
		Showname:  "管理员",
		Dashboard: "0",
		Avatar:    "",
		Roles: []string{"SA",
			"admin",
			"Auditor"},
	}
	fmt.Println(res.Token, res.UserInfo)
	return res, nil
}

// 获得所有的菜单
func (ags BackendGrpcService) QueryAllMenus(ctx context.Context, in *pb.MenuQueryReq) (*pb.MenuListRes, error) {
	return service.QueryMenuList(in)
}

// 获得单个菜单
func (ags BackendGrpcService) FetchMenu(ctx context.Context, in *pb.MenuQueryReq) (*pb.MenuRes, error) {
	return service.FetchMenu(in)
}

// 初始化接口
func (ags BackendGrpcService) InitRes(ctx context.Context, in *pb.InitResourceReq) (*pb.ResourcesRes, error) {
	return service.InitRes(in)

}

func TokenInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		err := authToken(ctx, info.FullMethod)
		if err != nil {
			log.Error("---------> the unaryServerInterceptor: " + err.Error())
			return nil, err
		}
		return handler(ctx, req)
	}
}

// 简单实现，用常量替代，这里如果不使用tls,还是jwt或自己实现加密更好些
func authToken(ctx context.Context, fullmethod string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("客户端校验失败")
	}
	var (
		appid  string
		appkey string
	)
	if val, ok := md["appid"]; ok {
		appid = val[0]
	}

	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}

	nouns := utils.EncryptGrpcCredentials(appid)
	if nouns != appkey {
		return errors.New("backend grpc token authenticate failed")
	}
	return nil
}
