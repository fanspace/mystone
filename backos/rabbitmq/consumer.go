package rabbitmq

import (
	"backos/commons"
	log "backos/logger"
	"backos/relations"
	"backos/settings"
	"encoding/json"
	"fmt"
)

type RecvProImp struct {
}

//// 实现消费者 消费消息失败 自动进入延时尝试  尝试3次之后入库db
/*
返回值 error 为nil  则表示该消息消费成功
否则消息会进入ttl延时队列  重复尝试消费3次
3次后消息如果还是失败 消息就执行失败  进入告警 FailAction
*/

type RabbitMsg struct {
	Type     string
	Act      string
	Username string
}

func (t *RecvProImp) Consumer(dataByte []byte) error {
	log.Info("********************    对用户进行账号锁定、解除操作   ********************")
	log.Info(string(dataByte))
	rmsg := new(RabbitMsg)
	err := json.Unmarshal(dataByte, rmsg)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	if rmsg.Type == "userStat" {
		if rmsg.Act == "ban" {
			err = commons.BanUserCache(rmsg.Username)
			if err != nil {
				log.Error(err.Error())
				return nil
			}
		} else {
			err = commons.ReleaseUserCache(rmsg.Username)
			if err != nil {
				log.Error(err.Error())
				return nil
			}
		}
	} else if rmsg.Type == "ipStat" {
		if rmsg.Act == "ban" {
			err = commons.BanIpCache(rmsg.Username)
			if err != nil {
				log.Error(err.Error())
				return nil
			}
		} else {
			err = commons.ReleaseIpCache(rmsg.Username)
			if err != nil {
				log.Error(err.Error())
				return nil
			}
		}
	} else {
		log.Error(relations.CUS_ERR_4002)
		return nil
	}

	return nil
}

//消息已经消费3次 失败了 请进行处理
/*
如果消息 消费3次后 仍然失败  此处可以根据情况 对消息进行告警提醒 或者 补偿  入库db  钉钉告警等等
*/
func (t *RecvProImp) FailAction(dataByte []byte) error {
	fmt.Println(string(dataByte))
	fmt.Println("任务处理失败了，我要进入db日志库了")
	return nil
}

func InitRabConsumer() {
	t := &RecvProImp{}
	Recv(QueueExchange{
		fmt.Sprintf("%s%s", relations.APP_NAME, "RabbitUser"),
		relations.RABBIT_USER_ROUTING_KEY,
		fmt.Sprintf("%s_user_exc", relations.APP_NAME),
		"fanout",
		fmt.Sprintf("amqp://%s:%s@%s:%d/", settings.Cfg.RabbitMqSettings.User, settings.Cfg.RabbitMqSettings.Pwd, settings.Cfg.RabbitMqSettings.Url, settings.Cfg.RabbitMqSettings.Port),
	}, t, 1)

}
