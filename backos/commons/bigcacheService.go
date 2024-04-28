package commons

import (
	log "backos/logger"
	"fmt"
	"github.com/allegro/bigcache/v2"
	"time"
)

var BCache *bigcache.BigCache

// minio文件专用
var FCache *bigcache.BigCache
var err error

func InitCache() {
	config := bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 32,

		// time after which entry can be evicted
		LifeWindow: 60 * time.Minute,

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		CleanWindow: 5 * time.Minute,

		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 1000 * 10 * 60,

		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 500,

		// prints information about additional memory allocation
		Verbose: true,

		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 8192,

		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		OnRemove: nil,

		// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A constant representing the reason will be passed through.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// Ignored if OnRemove is specified.
		OnRemoveWithReason: nil,
	}
	configFile := bigcache.Config{
		Shards:             32,
		LifeWindow:         12 * 60 * time.Minute,
		CleanWindow:        5 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		Verbose:            true,
		HardMaxCacheSize:   8192,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}
	BCache, err = bigcache.NewBigCache(config)
	if err != nil {
		log.Error(err.Error())
	}
	FCache, err = bigcache.NewBigCache(configFile)
	if err != nil {
		log.Error(err.Error())
	}
	err = LoadFromRedis()
	if err != nil {
		log.Error(err.Error())
	}

}

func CloseCache() {
	if BCache != nil {
		BCache.Close()
	}
}

func LoadFromRedis() error {
	total, blacks, err := GetAllBanUsers()
	if err != nil {
		log.Error(err.Error())
		return err
	} else {
		if total > 0 {
			for _, v := range blacks {
				BanUserCache(v)
			}
		}
	}
	return nil
}

func BanUserCache(username string) error {
	BCache.Set(fmt.Sprintf("ban_%s", username), []byte("ban"))
	log.Info("username： " + username + "已被加入黑名单")
	return nil
}

func ReleaseUserCache(username string) error {
	BCache.Delete(fmt.Sprintf("ban_%s", username))
	log.Info("username： " + username + "已从黑名单中删除")
	return nil
}

func IsUserBaned(username string) bool {
	entry, _ := BCache.Get(fmt.Sprintf("ban_%s", username))
	if string(entry) != "" {
		return true
	}
	return false
}

func BanIpCache(ip string) error {
	BCache.Set(fmt.Sprintf("ban_%s", ip), []byte("ban"))
	log.Info("ip： " + ip + "已被加入黑名单，持续1小时")
	return nil
}

func ReleaseIpCache(ip string) error {
	BCache.Delete(fmt.Sprintf("ban_%s", ip))
	log.Info("ip： " + ip + "已从黑名单中删除")
	return nil
}
func IsIpBand(ip string) bool {
	entry, _ := BCache.Get(fmt.Sprintf("ban_%s", ip))
	if string(entry) != "" {
		return true
	}
	return false
}

func SetBigCache(key string, val string) error {
	BCache.Set(key, []byte(val))
	return nil
}
func GetBigCache(key string) (string, error) {
	r, err := BCache.Get(key)
	return string(r), err
}

func SetFileCache(key string, val string) error {
	FCache.Set(key, []byte(val))
	return nil
}
func GetFileCache(key string) (string, error) {
	r, err := FCache.Get(key)
	return string(r), err
}
