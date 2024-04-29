package db

import (
	log "backgate/logger"
	"backgate/settings"
	"github.com/gomodule/redigo/redis"
	"os"
	"time"
)

var Rpool *redis.Pool

func InitRed() {
	Rpool = &redis.Pool{
		MaxIdle:     settings.Cfg.Database.RedisSettings.MaxIdle,
		MaxActive:   settings.Cfg.Database.RedisSettings.MaxActive,
		IdleTimeout: time.Duration(settings.Cfg.Database.RedisSettings.IdleTimeout) * time.Second,
		// Other pool configuration not shown in this example.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", settings.Cfg.Database.RedisSettings.Addr)
			if err != nil {
				//return nil, err
				log.Error(err.Error())
				os.Exit(1)
			}
			if _, err := c.Do("AUTH", settings.Cfg.Database.RedisSettings.Password); err != nil {
				c.Close()
				//return nil, err
				log.Error(err.Error())
				os.Exit(1)
			}
			if _, err := c.Do("SELECT", settings.Cfg.Database.RedisSettings.DB); err != nil {
				c.Close()
				//return nil, err
				log.Error(err.Error())
				os.Exit(1)
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func CloseRed() {
	Rpool.Close()
}
