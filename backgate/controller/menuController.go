package controller

import (
	log "backgate/logger"
	"backgate/relations"
	"backgate/service"
	pb "backgate/training"
	"github.com/duke-git/lancet/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// @Summary 查询管理端菜单列表
// @Description 查询管理端菜单列表
// @Tags System/Res
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param req body enlist.MenuQueryReq true "query menu List"
// @Success 200 {object} gin.H
// @Router /res/menu/list [post]
// @Date   9/17/2020 10:58 AM
func QueryAllMenus(c *gin.Context) {
	req := new(pb.MenuQueryReq)
	if err := c.ShouldBindJSON(req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 4002, "success": false, "msg": relations.CUS_ERR_4002})
		return
	}
	if !slice.Contain(relations.DOMAINS_LIMITED, req.Domain) {
		c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": false, "msg": relations.CUS_ERR_4008})
		return
	}

	res, err := service.DealGrpcCall(req, "QueryAllMenus", "backendrpc")
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
	c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": true, "msg": "", "content": res.(*pb.MenuListRes).Menus, "datetime": res.(*pb.MenuListRes).Total})
}

// @Summary 查询管理端菜单
// @Description 查询管理端菜单
// @Tags System/Res
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param req body enlist.MenuQueryReq true "query menu By Id"
// @Success 200 {object} gin.H
// @Router /res/menu/query [post]
// @Date   9/21/2020 10:58 AM
func FetchMenu(c *gin.Context) {
	req := new(pb.MenuQueryReq)
	if err := c.ShouldBindJSON(req); err != nil {
		log.Info(err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 4002, "success": false, "msg": relations.CUS_ERR_4002})
		return
	}
	if !slice.Contain(relations.DOMAINS_LIMITED, req.Domain) {
		c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": false, "msg": relations.CUS_ERR_4008})
		return
	}
	res, err := service.DealGrpcCall(req, "FetchMenu", "backendrpc")
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
	c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": true, "data": res})
	return
}
