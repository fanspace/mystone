package utils

import (
	log "backos/logger"
	"backos/relations"
	"backos/settings"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

type MyCustomClaims struct {
	Domain   string `json:"domain"`
	Usid     int64  `json:"usid"`
	AgencyId int64  `json:"agency_id"`
	Username string `json:"username"`
	UserType int32  `json:"userType"`
	Status   int32  `json:"status"`
	Device   int32  `json:"device"`
	jwt.StandardClaims
}
type MyClaim struct {
	Domain   string `json:"domain"`
	Usid     int64  `json:"usid"`
	AgencyId int64  `json:"agency_id"`
	Username string `json:"username"`
	UserType int32  `json:"userType"`
	Status   int32  `json:"status"`
	Device   int32  `json:"device"`
}

// 前端解析jwt
func ParseJwt(tokenstr string, logintype string) (bool, *MyClaim) {
	token, err := jwt.ParseWithClaims(strings.TrimSpace(tokenstr), &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if settings.Cfg.ReleaseMode {
			switch logintype {
			case "frontend":
				return []byte(relations.JWT_SECRET_STRING_PROD_NOR), nil
			case "com":
				return []byte(relations.JWT_SECRET_STRING_PROD_COM), nil
			case "backend":
				return []byte(relations.JWT_SECRET_STRING_PROD_MAN), nil
			default:
				return nil, errors.New(relations.CUS_ERR_4002)
			}

		} else {
			return []byte(relations.JWT_SECRET_STRING_DEV), nil
		}
	})
	if err != nil {
		log.Error(err.Error())
		return false, nil
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		mc := new(MyClaim)
		mc.Domain = claims.Domain
		mc.Usid = claims.Usid
		mc.Username = claims.Username
		mc.UserType = claims.UserType
		mc.Status = claims.Status
		mc.AgencyId = claims.AgencyId
		mc.Device = claims.Device
		return true, mc
	} else {
		log.Info(fmt.Sprintf("%v %v", claims.Usid, claims.StandardClaims.ExpiresAt))
		return false, nil
	}
}
