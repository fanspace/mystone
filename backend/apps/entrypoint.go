package apps

import (
	accsrv "backend/apps/account/service"
	"backend/apps/res/service"
	"context"
	"fmt"
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
