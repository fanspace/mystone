package accounts

import (
	"backos/apps/accounts/model"
	"backos/apps/accounts/service"
	"encoding/base64"
	"fmt"
	"net/url"
)

// DoLogin /*
// @title    DoLogin
// @description   User Sign In
// @auth      Lf             时间（2024/4/26 15:19）
// @param     req        model.AccountLoginReq         ""
// @return    res, err         *model.AccountLoginRes              ""

func DoLogin(req *model.AccountLoginReq) (*model.AccountLoginRes, error) {
	enEscapeUrl, _ := url.QueryUnescape(req.Password)
	sDec, err := base64.StdEncoding.DecodeString(enEscapeUrl)
	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
		return nil, err
	}
	req.Password = string(sDec)
	return service.Login(req)
}
