package service

import (
	pb "backend/training"
	"errors"
)

// 初始化接口
func InitRes(req *pb.InitResourceReq) (*pb.ResourcesRes, error) {
	res := new(pb.ResourcesRes)
	if len(req.ResCate) == 0 || len(req.ResItem) == 0 {
		return res, errors.New("length of list is zero")
	}
	//fmt.Println(req.ResCate)
	//fmt.Println(req.ResItem)

	return res, nil

}
