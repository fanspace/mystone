package rabbitmq

import (
	log "backos/logger"
	"backos/relations"
	"backos/settings"
	"encoding/json"
	"fmt"
)

/*type ZjRabbitMsg struct {
	Type     string
	Act      string
	Username string
}*/

func ProduceUserBanOrRelease(cate string, act string, username string) error {

	zm := new(RabbitMsg)
	zm.Type = cate
	zm.Act = act
	zm.Username = username
	msgdata, err := json.Marshal(zm)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	queueExchange := QueueExchange{
		fmt.Sprintf("%s%s", relations.APP_NAME, "RabbitUser"),
		relations.RABBIT_USER_ROUTING_KEY,
		fmt.Sprintf("%s_user_exc", relations.APP_NAME),
		"fanout",
		fmt.Sprintf("amqp://%s:%s@%s:%d/", settings.Cfg.RabbitMqSettings.User, settings.Cfg.RabbitMqSettings.Pwd, settings.Cfg.RabbitMqSettings.Url, settings.Cfg.RabbitMqSettings.Port),
	}

	Send(queueExchange, string(msgdata))
	return nil
}
