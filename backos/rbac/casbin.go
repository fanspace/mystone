package rbac

import (
	log "backos/logger"
	"backos/relations"
	"backos/settings"
	"encoding/json"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/xorm-adapter/v2"
	"github.com/go-redis/redis/v8"
	"os"
	"strings"
	"time"
)

var (
	Casbin *casbin.Enforcer
)

func InitCasbin() {
	a, err := xormadapter.NewAdapter("mysql", settings.Cfg.Database.MysqlSettings.Url, true)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	Casbin, err = casbin.NewEnforcer("conf/authz_model.conf", a)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	if !settings.Cfg.ReleaseMode {
		Casbin.EnableLog(true)
	} else {
		Casbin.EnableLog(false)
	}
	wc := WatcherOptions{
		Options: redis.Options{
			Network:  "tcp",
			Password: settings.Cfg.Database.RedisSettings.Password,
		},
		Channel:    fmt.Sprintf("/%s_rbac", relations.RED_CASBIN_CHANNEL_NAME),
		LocalID:    settings.Cfg.Smark,
		IgnoreSelf: true,
	}
	w, err := NewWatcher(settings.Cfg.Database.RedisSettings.Addr, wc)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	Casbin.SetWatcher(w)
	w.SetUpdateCallback(updateCallback)

}
func updateCallback(msg string) {
	timea := time.Now()
	begin := timea.UnixNano()
	sm := settings.Cfg.Smark
	msgcon := strings.Replace(msg, `\`, "", -1)
	msgs := new(CasbinUpdateMsg)
	err := json.Unmarshal([]byte(msgcon), msgs)
	if err != nil {
		log.Error(err.Error())
		err := Casbin.LoadPolicy()
		if err != nil {
			log.Error(err.Error())
		}
		end := time.Now().UnixNano()
		log.Info(fmt.Sprintf("********更新casbin规则：sender is %s - reciever is %s  共耗时 ********   %d ms", msg, sm, (end-begin)/1000/1000))
	} else if msgs.ID == "" {
		err := Casbin.LoadPolicy()
		if err != nil {
			log.Error(err.Error())
		}
		end := time.Now().UnixNano()
		log.Info(fmt.Sprintf("********更新casbin规则：sender is %s - reciever is %s  共耗时 ********   %d ms", msg, sm, (end-begin)/1000/1000))
	} else {
		if sm == msgs.ID {
			sm = "self"
			end := time.Now().UnixNano()
			log.Info(fmt.Sprintf("********更新casbin规则：sender is %s - %s  共耗时 ********   %d ms", msgs.ID, sm, (end-begin)/1000/1000))
		} else {
			err := Casbin.LoadPolicy()
			if err != nil {
				log.Error(err.Error())
			}
			end := time.Now().UnixNano()
			log.Info(fmt.Sprintf("********更新casbin规则：sender is %s - reciever is %s  共耗时 ********   %d ms", msgs.ID, sm, (end-begin)/1000/1000))
		}
	}

}
