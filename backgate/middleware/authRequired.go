package middleware

import (
	"backgate/service"
	"backgate/settings"
	"backgate/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func MustLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		logintype := c.Request.Header.Get("LoginType")
		jwttokens := c.Request.Header.Get("Authorization")
		if logintype == "" {
			//c.Set("mc", new(utils.MyClaim))
			//c.JSON(http.StatusUnauthorized, gin.H{"code": 403, "msg": "非认证的请求来源!"})
			//c.Abort()
			//return
		}
		if strings.Contains(jwttokens, "Bearer ") {
			jwttokens = strings.Replace(jwttokens, "Bearer ", "", 1)
			isauth, mc := utils.ParseJwt(jwttokens, logintype)
			if mc == nil || !isauth {
				c.Set("mc", new(utils.MyClaim))
				c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未登录或登录已过期!"})
				c.Abort()
				return
			} else {
				isBaned := service.IsUserBaned(mc.Username)
				//fitType := validUserTypeAndLoginType(logintype, mc.UserType)
				if isBaned {
					c.Set("mc", new(utils.MyClaim))
					c.JSON(http.StatusUnauthorized, gin.H{"code": 403, "msg": "账号已被锁定!"})
					c.Abort()
					return
				} else {
					c.Set("mc", mc)
					c.Next()
				}
				/*else if !fitType {
					c.Set("mc", new(utils.MyClaim))
					c.JSON(http.StatusUnauthorized, gin.H{"code": 403, "msg": "禁止该用户登录!"})
					c.Abort()
					return
				}*/

			}
		} else {
			fmt.Println("jwt token is not Bearer")
			c.Set("mc", new(utils.MyClaim))
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未登录或登录已过期!!!"})
			c.Abort()
			return
		}
	}
}

func validUserTypeAndLoginType(logintype string, usertype int32) bool {
	switch logintype {
	case "fore":
		if usertype == settings.VOptions.GetInt32("UserType.NORMAL") || usertype == settings.VOptions.GetInt32("UserType.INACT") {
			return true
		}
	case "com":
		if usertype == settings.VOptions.GetInt32("UserType.AGENCY") {
			return true
		}
	case "back":
		if usertype == settings.VOptions.GetInt32("UserType.MANAGER") || usertype == settings.VOptions.GetInt32("UserType.ROOT") {
			return true
		}
	default:
		return false
	}
	return false
}
