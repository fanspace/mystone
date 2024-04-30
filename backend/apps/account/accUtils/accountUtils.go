package accUtils

import (
	"backend/apps/account/model"
	"backend/commons"
	log "backend/logger"
	"backend/relations"
	"backend/settings"
	"backend/utils"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	PWD_SECRET_STRING_FORE = "z2kxMGptMzBiMXJtbTxubXA0P2tleTg"
	PWD_SECRET_STRING_BACK = "NYTY4NGUwODgxMzMxNGE0YjYxZTM5ZY" // DefaultPageSize is
	BAN_TYPE_USERNAME      = "userStat"
	BAN_TYPE_IPADDRESS     = "ipStat"
)

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
	redpin, err := commons.GetPinCode(keyname)
	if err != nil {
		return err
	}
	if pin != redpin {
		return errors.New(relations.CUS_ERR_2027)
	}
	return nil
}

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
	err := commons.SetPinCode(keyname, pin)
	if err != nil {
		return "", "", err
	}
	return sncode, pin, nil
}

// 验证用户  登录使用
func AuthUser(account string, password string) (*model.Account, error) {
	user, err := model.GetAccountByUsername(account)
	if err != nil {
		return nil, err
	}
	if user.Status == settings.VOptions.GetInt32("UserStatus.USER_STATUS_BAND") {
		return nil, errors.New(relations.CUS_ERR_1011)
	}
	et, err := commons.LoginErr(user.Username, false)
	if err != nil {
		return nil, errors.New(relations.CUS_ERR_1012)
	}
	if et >= settings.VOptions.GetInt("NumVal.MAX_LOGIN_ERR_TIMES") {
		return nil, errors.New(fmt.Sprintf("用户名或密码输入错误次数过多，请%d分钟后再试", settings.VOptions.GetInt("NumVal.LOGIN_ERR_REFRESH_TIME")))
	} else {
		if isok := CheckUserAndPwd(user, password); isok {
			commons.RedisDelKey(fmt.Sprintf("%s_%s_%s", relations.STR_PREFIX, relations.STR_LOGIN_ERR, user.Username))
			return user, nil
		} else {
			et, err := commons.LoginErr(user.Username, true)
			if err != nil {
				return nil, errors.New(relations.CUS_ERR_1012)
			}
			if et >= settings.VOptions.GetInt("NumVal.MAX_LOGIN_ERR_TIMES") {
				return nil, errors.New(fmt.Sprintf("用户名或密码输错%d次，该账号%d分钟内无法登录", settings.VOptions.GetInt("NumVal.MAX_LOGIN_ERR_TIMES"), settings.VOptions.GetInt("NumVal.LOGIN_ERR_REFRESH_TIME")))
			}
			return nil, errors.New(fmt.Sprintf("用户名或密码已输错 %d / 5 次", et))
		}
	}
}

// 验证 密码
func CheckUserAndPwd(user *model.Account, inputpassword string) bool {
	if user == nil || user.Password == "" || user.PasswordHash == "" || inputpassword == "" {
		log.Error("用户出现关键字段为空的情况,传参出现异常：" + user.Username)
		return false
	}

	shastr1 := GenAccountPwd(inputpassword, user.PasswordHash, user.UserType)
	if shastr1 == user.Password {
		return true
	}
	return false
}

func GenAccountPwd(password string, passwordhash string, usertype int32) string {
	sec := ""
	if usertype < settings.VOptions.GetInt32("UserType.TEACHER") {
		sec = PWD_SECRET_STRING_FORE
	} else if usertype >= settings.VOptions.GetInt32("UserType.TEACHER") && usertype < settings.VOptions.GetInt32("UserType.MANAGER") {
		//sec = consts.PWD_SECRET_STRING_COM
		sec = PWD_SECRET_STRING_BACK
	} else {
		sec = PWD_SECRET_STRING_BACK
	}
	ori := password + passwordhash
	pwd := []byte(ori)
	key := []byte(sec)
	m := hmac.New(sha256.New, key)
	m.Write(pwd)
	signature := strings.ToUpper(hex.EncodeToString(m.Sum(nil)))
	//fmt.Println(signature)
	return signature
}
