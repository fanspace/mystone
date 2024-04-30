package service

import (
	"backend/apps/account/accUtils"
	"backend/relations"
	"backend/settings"
	pb "backend/training"
	"errors"
)

func UserLogin(req *pb.AccountLoginReq) (string, int64, error) {
	acc, err := accUtils.AuthUser(req.Username, req.Password)
	if err != nil {
		// 这里开始锁定
		go BanErrIp(req.Ip)
		return "", 0, err
	}
	if acc == nil {
		return "", 0, errors.New(relations.CUS_ERR_1010)
	}

	if acc.Status == settings.VOptions.GetInt32("UserStatus.USER_STATUS_BAND") {
		return "", 0, errors.New(relations.CUS_ERR_1011)
	}

	if acc.UserType <= relations.ZERO {
		return "", 0, errors.New(relations.CUS_ERR_1015)
	}
	if acc.UserType < settings.VOptions.GetInt32("USER_TYPE_TEACHER") {
		if req.LoginType == "backend" {
			return "", 0, errors.New(relations.CUS_ERR_1109)
		}
	}
	return "", 0, nil
}
