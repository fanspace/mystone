package service

import (
	"backgate/db"
	log "backgate/logger"
	"backgate/relations"
	"backgate/settings"
	"errors"
	red "github.com/gomodule/redigo/redis"
)

func GetAllBanUsers() (int64, []string, error) {
	redis := db.Rpool.Get()
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

func SetPinCode(sncode string, pin string) error {
	redis := db.Rpool.Get()
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
	redis := db.Rpool.Get()
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

// expireTime 0表示不过期, 单位为秒
func SetSimpleKey(keyname string, expireTime int, value interface{}) error {
	redis := db.Rpool.Get()
	defer redis.Close()
	_, err := redis.Do("SET", keyname, value)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if expireTime > 0 {
		_, err = redis.Do("EXPIRE", keyname, expireTime)
		if err != nil {
			log.Error("redis set EXPIRE  :" + keyname + " 失败")
			return err
		}
	}

	return nil
}

func GetSimpleKeyWithoutTTL(keyname string) (interface{}, error) {
	redis := db.Rpool.Get()
	defer redis.Close()
	res, err := redis.Do("GET", keyname)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return res, nil
}

func GetSimpleKey(keyname string) (interface{}, error) {

	redis := db.Rpool.Get()
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
