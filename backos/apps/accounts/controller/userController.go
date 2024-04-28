package controller

import (
	"backos/relations"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary 退出登录
// @Description 用户退出登录
// @Tags Account
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} gin.H
// @Router /logout [get]
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": true})
	return
}
