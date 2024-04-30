package apps

import (
	accsrv "backend/apps/account/service"
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
	res.Tokens = "SCUI.Administrator.Auth"
	res.UserInfo = &pb.UserInfo{
		UserId:    1,
		Username:  "Administrator",
		Showname:  "管理员",
		Dashboard: "0",
		Avatar:    "",
		Roles: []string{"SA",
			"admin",
			"Auditor"},
	}
	return res, nil
}
