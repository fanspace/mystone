package service

import (
	log "backgate/logger"
	"backgate/relations"
	pb "backgate/training"
	"errors"
)

func DoLogin(req *pb.AccountLoginReq) (*pb.AccountLoginRes, error) {
	isbaned := IsUserBaned(req.Username)
	if isbaned {
		return nil, errors.New(relations.CUS_ERR_1011)
	}
	err := ValidPinCode(req.Username, req.Sncode, req.Pin)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	res, err := DealGrpcCall(req, "HandleAccountLogin", "backendrpc")
	return res.(*pb.AccountLoginRes), err
}
