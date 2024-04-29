package service

import (
	"backgate/db"
	log "backgate/logger"
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
