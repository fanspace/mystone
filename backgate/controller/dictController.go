package controller

import (
	"backgate/relations"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary dict:add
// @Description 字典管理|新增字典
// @Tags dictMgr
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param req body training.MenuQueryReq true "query menu By Id"
// @Success 200 {object} map[string]any
// @Router /dictMgr/add [post]
// @Date   9/21/2020 10:58 AM
func AddDict(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": true, "data": ""})
	return
}

// @Summary dict:update
// @Description 字典管理|编辑字典
// @Tags dictMgr
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param req body training.MenuQueryReq true "query menu By Id"
// @Success 200 {object} map[string]any
// @Router /dictMgr/update [post]
// @Date   9/21/2020 10:58 AM
func UpdateDict(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": true, "data": ""})
	return
}

// @Summary dict:add
// @Description 字典管理|删除字典
// @Tags dictMgr
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param req body training.MenuQueryReq true "query menu By Id"
// @Success 200 {object} map[string]any
// @Router /dictMgr/del [post]
// @Date   9/21/2020 10:58 AM
func DelDict(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": true, "data": ""})
	return
}
