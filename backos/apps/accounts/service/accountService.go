package service

import (
	"backos/apps/accounts/model"
	"backos/commons"
	log "backos/logger"
	"backos/relations"
	"backos/settings"
	"errors"
	"fmt"
)

func Login(req *model.AccountLoginReq) (*model.AccountLoginRes, error) {
	fmt.Println(req)
	res := new(model.AccountLoginRes)
	isbaned := commons.IsUserBaned(req.Username)
	if isbaned {
		return nil, errors.New(relations.CUS_ERR_1011)
	}

	err := ValidPinCode(req.Username, req.Sncode, req.Pin)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	acc, err := AuthUser(req.Username, req.Password)
	if err != nil {
		// 这里开始锁定
		go BanErrIp(req.Ip)
		return res, err
	}
	if acc == nil {
		return res, errors.New(relations.CUS_ERR_1010)
	}

	if acc.Status == settings.VOptions.GetInt32("UserStatus.USER_STATUS_BAND") {
		return res, errors.New(relations.CUS_ERR_1011)
	}

	if acc.UserType <= relations.ZERO {
		return res, errors.New(relations.CUS_ERR_1015)
	}
	if acc.UserType < settings.VOptions.GetInt32("USER_TYPE_TEACHER") {
		if req.LoginType == "backend" {
			return res, errors.New(relations.CUS_ERR_1109)
		}
	}

	res = &model.AccountLoginRes{
		Token: "SCUI.Administrator.Auth",
		UserInfo: &model.UserInfo{
			UserId:    1,
			Username:  "Administrator",
			Showname:  "管理员",
			Dashboard: "0",
			Avatar:    "",
			Roles: []string{"SA",
				"admin",
				"Auditor"},
		},
	}
	return res, nil
}
