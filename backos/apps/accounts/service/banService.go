package service

import (
	"backos/apps/accounts/consts"
	"backos/commons"
	log "backos/logger"
	"backos/rabbitmq"
	"backos/settings"
	"backos/utils"
	"fmt"
)

/**
 *  @author: lf
 *  @Date: 12/18/2020 12:47 PM
 *  @Description:  查询ip 是否在白名单中
 */
func IsWhiteIp(ip string) bool {
	whiteIps := settings.VOptions.GetStringSlice("WhiteIps")
	iswhite, err := utils.Contains(ip, whiteIps)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	return iswhite
}

/**
 *  @author: lf
 *  @Date: 12/18/2020 1:02 PM
 *  @Description: 处理ip,加黑
 */
func BanErrIp(ip string) error {
	iswhite := IsWhiteIp(ip)
	if iswhite {
		return nil
	}
	keyname := fmt.Sprintf("lerr_%s", ip)
	errtimes, err := commons.GetIpErrTimes(keyname)
	if err != nil {
		// 不管是否取回值，增加
	}
	if errtimes >= settings.VOptions.GetInt64("NumVal.MAX_LOGIN_ERR_IP_TIMES") {
		return rabbitmq.ProduceUserBanOrRelease(consts.BAN_TYPE_IPADDRESS, "ban", ip)
	}
	return commons.SetIpErrTimes(keyname)
}

/**
 *  @author: lf
 *  @Date: 12/18/2020 2:45 PM
 *  @Description: 解除
 */
// ，从cache中清除（使用广播）   从redis中清除，如果再达到20，则继续锁定，除非加入白名单
func ReleaseErrIp(ip string) error {
	commons.RedisDelKey(fmt.Sprintf("lerr_%s", ip))
	return rabbitmq.ProduceUserBanOrRelease(consts.BAN_TYPE_IPADDRESS, "release", ip)
}
