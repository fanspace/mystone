package service

import (
	"backgate/relations"
	"backgate/utils"
	"errors"
	"fmt"
	"time"
)

/**
 *  @author: lf
 *  @Date: 12/15/2020 2:56 PM
 *  @Description: 生成验证码
 */
func GeneratePinCode(username string) (string, string, error) {
	pre := utils.RandomStr(3)
	now := time.Now().Unix()
	ext := utils.RandomStr(3)
	sncode := fmt.Sprintf("%s%d%s", pre, now, ext)
	keyname := fmt.Sprintf("%s_%s_%s", relations.APP_NAME, username, sncode)
	pin := utils.RandomInt(4)
	err := SetPinCode(keyname, pin)
	if err != nil {
		return "", "", err
	}
	return sncode, pin, nil
}

/**
 *  @author: lf
 *  @Date: 12/15/2020 2:57 PM
 *  @Description: 核对验证码
 */
func ValidPinCode(username string, sncode string, pin string) error {
	if username == "" || sncode == "" || pin == "" {
		return errors.New(relations.CUS_ERR_4002)
	}
	keyname := fmt.Sprintf("%s_%s_%s", relations.APP_NAME, username, sncode)
	redpin, err := GetPinCode(keyname)
	if err != nil {
		return err
	}
	if pin != redpin {
		return errors.New(relations.CUS_ERR_2027)
	}
	return nil
}

/**
 *  @Classname AccountService
 *  @author: lf6128@163.com
 *  @Date: 2022/5/30 10:05
 *  @Description:
 */
func ValidRegPinCode(mobile string, pin string) error {
	if mobile == "" || pin == "" {
		return errors.New(relations.CUS_ERR_4002)
	}
	keyname := fmt.Sprintf("%s_%s_%s", relations.APP_NAME, mobile, pin)
	redpin, err := GetPinCode(keyname)
	if err != nil {
		return err
	}
	if pin != redpin {
		return errors.New(relations.CUS_ERR_2027)
	}
	return nil
}
