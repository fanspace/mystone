package controller

import (
	log "backgate/logger"
	"backgate/rbac"
	"backgate/relations"
	"backgate/service"
	pb "backgate/training"
	"backgate/utils"
	"backgate/viewmodel"
	"github.com/duke-git/lancet/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// @Summary 菜单
// @Description 查询我的菜单
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]any
// @Router /auth/menu/mine [get]
// @Date   9/17/2020 10:58 AM
func QueryMyMenu(c *gin.Context) {
	req := new(pb.MenuQueryReq)
	req.Type = 0
	req.Domain = "backend"
	mc := c.MustGet("mc").(*utils.MyClaim)
	isRoot, _ := rbac.Casbin.HasRoleForUser(mc.Username, "root")
	if isRoot {
		req.Type = 0
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
	datas := new(viewmodel.MenuResList)
	datas.Menu = res.(*pb.MenuListRes).Menus
	datas.DashboardGrid = []string{"welcome", "ver", "time", "progress", "echarts", "about"}
	datas.Permissions = []string{"list.add", "list.edit", "list.delete", "user.add", "user.edit", "user.delete"}
	c.JSON(http.StatusOK, gin.H{"code": 200, "success": true, "msg": "", "data": datas, "datetime": res.(*pb.MenuListRes).Total})

}

// @Summary menu:list
// @Description 菜单管理|查询菜单列表
// @Tags MenuMgr
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param req body training.MenuQueryReq true "query menu List"
// @Success 200 {object} map[string]any
// @Router /menuMgr/list [post]
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
	datas := new(viewmodel.MenuResList)
	datas.Menu = res.(*pb.MenuListRes).Menus
	datas.DashboardGrid = []string{"welcome", "ver", "time", "progress", "echarts", "about"}
	datas.Permissions = []string{"list.add", "list.edit", "list.delete", "user.add", "user.edit", "user.delete"}
	c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": true, "msg": "", "data": datas, "datetime": res.(*pb.MenuListRes).Total})
}

// @Summary 菜单管理
// @Description 菜单管理|查询菜单
// @Tags MenuMgr
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param req body training.MenuQueryReq true "query menu By Id"
// @Success 200 {object} map[string]any
// @Router /menuMgr/query [post]
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

// @Summary 菜单管理
// @Description 菜单管理|新增菜单
// @Tags MenuMgr
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param req body training.MenuQueryReq true "query menu By Id"
// @Success 200 {object} map[string]any
// @Router /menuMgr/add [post]
// @Date   9/21/2020 10:58 AM
func AddMenu(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": true, "data": ""})
	return
}

// @Summary 菜单管理
// @Description 菜单管理|编辑菜单
// @Tags MenuMgr
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param req body training.MenuQueryReq true "query menu By Id"
// @Success 200 {object} map[string]any
// @Router /menuMgr/update [post]
// @Date   9/21/2020 10:58 AM
func UpdateMenu(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": true, "data": ""})
	return
}

// @Summary 菜单管理
// @Description 菜单管理|删除菜单
// @Tags MenuMgr
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param req body training.MenuQueryReq true "query menu By Id"
// @Success 200 {object} map[string]any
// @Router /menuMgr/del [post]
// @Date   9/21/2020 10:58 AM
func DelMenu(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": relations.WEB_STATUS_BACK, "success": true, "data": ""})
	return
}
