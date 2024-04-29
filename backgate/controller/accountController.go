package controller

import (
	log "backgate/logger"
	"backgate/relations"
	"backgate/service"
	pb "backgate/training"
	"backgate/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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
	req := new(pb.AccountLoginReq)
	if err := c.ShouldBindJSON(req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 4002, "success": false, "msg": relations.CUS_ERR_4002})
		return
	}
	fmt.Println(req)
	if req.Username == "" || req.Password == "" || req.Pin == "" || req.Sncode == "" || req.LoginType == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4002, "success": false, "msg": relations.CUS_ERR_4002})
		return
	}
	req.Ip = utils.RemoteIp(c.Request)

	res, err := service.DoLogin(req)
	if err != nil {
		if strings.Index(err.Error(), "desc =") > 0 {
			msg := strings.Split(err.Error(), "desc = ")[1]
			c.JSON(http.StatusOK, gin.H{"code": 99999, "success": false, "msg": msg})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"code": 99999, "success": false, "msg": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    relations.WEB_STATUS_BACK,
		"success": true,
		"data":    res,
		"message": "",
	})
	return
}
