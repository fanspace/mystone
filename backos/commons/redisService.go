package commons

import (
	. "backos/db"
	log "backos/logger"
	"backos/relations"
	"backos/settings"
	"errors"
	"fmt"
	red "github.com/gomodule/redigo/redis"
	"strings"
)

func SetPinCode(sncode string, pin string) error {
	redis := Rpool.Get()
	defer redis.Close()
	_, err := redis.Do("SET", sncode, pin)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	_, err = redis.Do("EXPIRE", sncode, settings.VOptions.GetInt32("NumVal.PIN_CODE_INVALID_MINUTE")*60) //5分钟失效
	if err != nil {
		log.Error("redis set EXPIRE  :" + sncode + " 失败")
		return err
	}
	return nil
}

func GetPinCode(sncode string) (string, error) {
	redis := Rpool.Get()
	defer redis.Close()
	tmpttl, err := red.Int64(redis.Do("Ttl", sncode))
	if err != nil {
		log.Error(err.Error())
		return "", errors.New(relations.CUS_ERR_1009)
	}
	if tmpttl == -2 {
		return "", errors.New(relations.CUS_ERR_2026)
	} else if tmpttl == -1 {
		return "", errors.New(relations.CUS_ERR_2025)
	} else {
		pin, err := red.String(redis.Do("GET", sncode))
		if err != nil {
			log.Error(err.Error())
			return "", err
		}
		if pin == "" {
			return "", errors.New(relations.CUS_ERR_4004)
		}
		_, err = redis.Do("Del", sncode)
		if err != nil {
			log.Error("redis Delete  :" + sncode + " 失败")
			return "", err
		}
		return pin, nil
	}
}

func GetAllBanUsers() (int64, []string, error) {
	redis := Rpool.Get()
	defer redis.Close()
	res := make([]string, 0)
	keyname := "enlist_banuser"
	total, err := red.Int64(redis.Do("SCARD", keyname))
	if err != nil {
		log.Error(err.Error())
		return 0, res, err
	}
	if total > 0 {
		res, err = red.Strings(redis.Do("SMEMBERS", keyname))
		if err != nil {
			log.Error(err.Error())
			return 0, res, err
		}
	}
	return total, res, nil
}

//  *********************** 微信相关
// qr sence

func SetSenceId(senceid string, usid int64) error {
	if senceid == "" || usid == 0 {
		return errors.New(relations.CUS_ERR_4002)
	}
	redis := Rpool.Get()
	defer redis.Close()
	_, err = redis.Do("Set", senceid, fmt.Sprintf("%d", usid))
	if err != nil {
		log.Error("wx qrcode sence set failed")
		return err
	}

	_, err = redis.Do("EXPIRE", senceid, settings.VOptions.GetInt32("NumVal.PIN_CODE_INVALID_MINUTE")*60) //5分钟失效
	if err != nil {
		log.Error("wx qrcode sence set EXPIRE  failed")
		return err
	}
	return nil
}

func ClearQrCodeSenceId(keyname string) error {
	if keyname == "" {
		return errors.New(relations.CUS_ERR_4002)
	}
	redis := Rpool.Get()
	defer redis.Close()
	var err error
	_, err = redis.Do("Del", keyname)
	if err != nil {
		log.Error("redis Delete  :" + keyname + " 失败")
		return err
	}
	return nil
}

func CheckWxLogin(sencestr string) (int64, string) {
	redis := Rpool.Get()
	defer redis.Close()
	usid, err := red.Int64(redis.Do("HGet", sencestr, "usid"))
	if err != nil {
		if strings.Index(err.Error(), "nil returned") < 0 {
			log.Error(err.Error())
		}
		return 0, ""
	}
	token, err := red.String(redis.Do("HGet", sencestr, "token"))
	if err != nil {
		if strings.Index(err.Error(), "nil returned") < 0 {
			log.Error(err.Error())
		}
		return 0, ""
	}

	return usid, token
}

/*    uniform    */

func SetSimpleKey(keyname string, value interface{}) error {
	redis := Rpool.Get()
	defer redis.Close()
	_, err := redis.Do("SET", keyname, value)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	_, err = redis.Do("EXPIRE", keyname, 60*60) //60分钟失效
	if err != nil {
		log.Error("redis set EXPIRE  :" + keyname + " 失败")
		return err
	}
	return nil

}

func GetSimpleKey(keyname string) (interface{}, error) {
	redis := Rpool.Get()
	defer redis.Close()
	tmpttl, err := red.Int64(redis.Do("Ttl", keyname))
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New(relations.CUS_ERR_1009)
	}
	if tmpttl == -2 {
		return nil, nil
	} else if tmpttl == -1 {
		_, err := redis.Do("DEL", keyname)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		return nil, nil
	} else {
		res, err := redis.Do("GET", keyname)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		return res, nil
	}
}

// 读取密码错误次数  inputerror == ture 登录后获取并修改，   inputerr = false  登录前获取
func LoginErr(username string, inputerror bool) (int, error) {
	redis := Rpool.Get()
	defer redis.Close()
	keyname := fmt.Sprintf("%s_%s_%s", relations.STR_PREFIX, relations.STR_LOGIN_ERR, username)
	tmpttl, err := red.Int64(redis.Do("Ttl", keyname))
	if err != nil {
		log.Error(err.Error())
		return 0, errors.New(relations.CUS_ERR_1009)
	}
	if tmpttl == -2 {
		if inputerror {
			_, err := redis.Do("SET", keyname, 1)
			if err != nil {
				log.Error("redis wrong at do set " + keyname)
				return 0, errors.New(relations.CUS_ERR_1008)
			}
			_, err = redis.Do("EXPIRE", keyname, settings.VOptions.GetInt("NumVal.LOGIN_ERR_REFRESH_TIME")*60)

			if err != nil {
				log.Error("redis set EXPIRE  :" + keyname + "  失败")
				return 0, err
			}
			return 1, nil
		} else {
			return 0, nil
		}

	} else {
		errtimes, err := red.Int(redis.Do("GET", keyname))
		if err != nil {
			log.Error(err.Error())
			return 0, err
		}
		if errtimes >= settings.VOptions.GetInt("NumVal.MAX_LOGIN_ERR_TIMES") {
			return errtimes, nil
		} else {
			if inputerror {
				nowtimes := errtimes + 1
				_, err := redis.Do("SET", keyname, nowtimes)
				if err != nil {
					log.Error("redis wrong at do set " + keyname)
					return 0, nil
				}
				interval := settings.VOptions.GetInt("NumVal.LOGIN_ERR_REFRESH_TIME") * 60
				_, err = redis.Do("EXPIRE", keyname, interval)

				if err != nil {
					log.Error("redis set EXPIRE  :" + keyname + "  失败")
					return 0, err
				}
				return nowtimes, nil
			}
			return errtimes, nil
		}
	}
}

// 删除登录错误次数
func RedisDelKey(keyname string) bool {
	redis := Rpool.Get()
	defer redis.Close()
	_, err := redis.Do("Del", keyname)
	if err != nil {
		log.Error("redis Delete  :" + keyname + " 失败")
		return false
	}
	return true
}

func GetIpErrTimes(ipname string) (int64, error) {
	redis := Rpool.Get()
	defer redis.Close()
	res, err := red.Int64(redis.Do("GET", ipname))
	if err != nil {
		log.Error(err.Error())
		return relations.ZERO, errors.New(relations.CUS_ERR_1009)
	}
	return res, nil
}

func SetIpErrTimes(ipname string) error {

	redis := Rpool.Get()
	defer redis.Close()
	_, err := redis.Do("INCR", ipname)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	tmpttl, err := red.Int64(redis.Do("TTL", ipname))
	if err != nil {
		log.Error(err.Error())
		return errors.New(relations.CUS_ERR_1009)
	}
	if tmpttl <= 0 {
		_, err = redis.Do("EXPIRE", ipname, settings.VOptions.GetInt32("NumVal.MAX_LOGIN_ERR_IP_MINUTE")*60) //2分钟失效
		if err != nil {
			log.Error("redis set EXPIRE  :" + ipname + " 失败")
			return err
		}
	}
	return nil
}

func RedisStoreSsm(sncode string, authcode int, usid int) bool {
	redis := Rpool.Get()
	defer redis.Close()
	redis.Send("MULTI")
	redis.Send("HSet", sncode, "authcode", authcode)
	redis.Send("HSet", sncode, "usid", usid)

	_, err := redis.Do("EXEC")
	if err != nil {
		log.Error("wrong at reset pwd for usid ")
		return false
	}
	_, err = redis.Do("EXPIRE", sncode, 1800)
	if err != nil {
		log.Error("redis set EXPIRE  :" + sncode + " 失败")
		return false
	}
	return true
}

func RedisReadSsm(sncode string, authcode string, usid int) error {
	redis := Rpool.Get()
	defer redis.Close()
	tmpttl, err := red.Int64(redis.Do("Ttl", sncode))
	if err != nil {
		log.Error(err.Error())
		return errors.New(relations.CUS_ERR_1009)
	}
	if tmpttl == -2 {
		return errors.New(relations.CUS_ERR_2026)
	} else if tmpttl == -1 {
		return errors.New(relations.CUS_ERR_2025)
	} else {
		tmpauth, err := red.String(redis.Do("HGet", sncode, "authcode"))
		if err != nil {
			log.Error(err.Error())
			return errors.New(relations.CUS_ERR_1009)
		}
		tmpuid, err := red.Int(redis.Do("HGet", sncode, "usid"))
		if err != nil {
			log.Error(err.Error())
			return errors.New(relations.CUS_ERR_1009)
		}

		if tmpuid != usid || tmpauth != authcode {
			return errors.New(relations.CUS_ERR_2026)
		}
		return nil
	}

}
