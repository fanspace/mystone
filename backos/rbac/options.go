package rbac

import (
	"backos/relations"
	"fmt"
	rds "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type WatcherOptions struct {
	rds.Options
	Channel    string
	IgnoreSelf bool
	LocalID    string
	Password   string
}

func initConfig(option *WatcherOptions) {
	if option.LocalID == "" {
		option.LocalID = uuid.New().String()
	}
	if option.Channel == "" {
		option.Channel = fmt.Sprintf("/%s_rbac", relations.RED_CASBIN_CHANNEL_NAME)
	}
}

// casbin update msg
type CasbinUpdateMsg struct {
	Method string
	ID     string
	Params string
}
