package controller

import (
	"backos/apps/accounts"
	"backos/apps/accounts/model"
	log "backos/logger"
	"backos/relations"
	"backos/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary 用户登录
// @Description
// @Tags Account
// @Accept json
// @Produce json
// @Param req body enlist.AccountLoginReq true "user login"
// @Success 200 {object} gin.H
// @Router /login [post]
func Login(c *gin.Context) {
	req := new(model.AccountLoginReq)
	if err := c.ShouldBindJSON(req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 4002, "success": false, "msg": relations.CUS_ERR_4002})
		return
	}

	if req.Username == "" || req.Password == "" || req.Pin == "" || req.Sncode == "" || req.LoginType == "" {
		fmt.Println(req)
		c.JSON(http.StatusOK, gin.H{"code": 4002, "success": false, "msg": relations.CUS_ERR_4002})
		return
	}

	req.Ip = utils.RemoteIp(c.Request)

	res, err := accounts.DoLogin(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 99999, "success": false, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    relations.WEB_STATUS_BACK,
		"success": true,
		"data":    res,
		"message": "",
	})
	return
}

// @Summary 生成验证码
// @Description 需要用户名
// @Tags Account
// @Accept json
// @Produce json
// @Param username path string true "get pin code"
// @Success 200 {object} gin.H
// @Router /pin/{username} [get]
func GenPin(c *gin.Context) {
	uname := c.Param("username")
	if uname == "" {
		c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "msg": relations.CUS_ERR_4002, "success": false})
		return
	}
	sncode, pin, err := accounts.GeneratePinCode(uname)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": false, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    relations.WEB_STATUS_BACK,
		"success": true,
		"sncode":  sncode,
		"pin":     pin,
	})
	return
}
