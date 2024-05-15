package service

import (
	log "backgate/logger"
	pb "backgate/training"
)

func HandleMenus(req *pb.MenuHandleReq) (*pb.CommonResponse, error) {
	// 获取菜单
	// 获取菜单权限
	// 获取菜单按钮权限
	res, err := DealGrpcCall(req, "HandleMenu", "backendrpc")
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return res.(*pb.CommonResponse), nil
}
