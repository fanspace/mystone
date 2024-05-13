package controller

import (
	"backgate/relations"
	"backgate/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckStat(c *gin.Context) {
	c.String(200, "ok")
}

func Ping(c *gin.Context) {
	c.String(200, "pong")
}

func ListApis(c *gin.Context) {
	_, err := service.ListAllApis()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": false, "data": ""})
		return
	}
	c.String(200, "pong")
}

// @Summary api:add
// @Description 接口管理|新增接口
// @Tags apiMgr
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param req body training.MenuQueryReq true "query menu By Id"
// @Success 200 {object} map[string]any
// @Router /apiMgr/add [post]
// @Date   9/21/2020 10:58 AM
func AddApi(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": true, "data": ""})
	return
}

// @Summary api:update
// @Description 接口管理|编辑接口
// @Tags apiMgr
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param req body training.MenuQueryReq true "query menu By Id"
// @Success 200 {object} map[string]any
// @Router /apiMgr/update [post]
// @Date   9/21/2020 10:58 AM
func UpdateApi(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": true, "data": ""})
	return
}

// @Summary api:del
// @Description 接口管理|删除接口
// @Tags apiMgr
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param req body training.MenuQueryReq true "query menu By Id"
// @Success 200 {object} map[string]any
// @Router /apiMgr/del [post]
// @Date   9/21/2020 10:58 AM
func DelApi(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": true, "data": ""})
	return
}
