package middleware

import (
	"backos/commons"
	"backos/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FilterBan() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := utils.RemoteIp(c.Request)
		ipband := commons.IsIpBand(ip)
		if ipband {
			c.Set("mc", new(utils.MyClaim))
			c.JSON(http.StatusUnauthorized, gin.H{"code": 403, "msg": "检测到异常操作，该ip已被封禁60分钟"})
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
