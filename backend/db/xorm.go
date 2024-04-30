package db

import (
	log "backend/logger"
	"backend/settings"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"os"
	"time"
	"xorm.io/xorm"
)

var Orm *xorm.Engine
var Rpool *redis.Pool

func InitDb() {

	var err error
	Orm, err = xorm.NewEngine(settings.Cfg.Database.MysqlSettings.DriverName, settings.Cfg.Database.MysqlSettings.Url)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	//显示sql
	if settings.Cfg.ReleaseMode {
		Orm.ShowSQL(false)
	} else {
		Orm.ShowSQL(true)
	}
	//设置时区
	Orm.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	/*
			err = Orm.Sync2(new(ZjMenuEntity))
		if err != nil {
			log.Error(err.Error())
		}

	*/

	/// Orm.SetDefaultCacher(NewRedisCacher(settings.Cfg.RedisSettings.Addr, settings.Cfg.RedisSettings.Password, settings.Cfg.RedisSettings.DB, DEFAULT_EXPIRATION, nil))
	//Orm.MapCacher(&CasbinRule{}, nil)
	// Orm.MapCacher(&SysUserStats{}, nil)
	/**/
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

func CloseDb() {
	Orm.Close()
	Rpool.Close()

}
