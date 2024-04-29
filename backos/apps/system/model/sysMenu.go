package model

import (
	"backos/db"
	log "backos/logger"
	"backos/relations"
	"errors"
	"fmt"
)

type MenuEntity struct {
	Id          int64  `json:"id" xorm:"not null pk autoincr comment('主键') BIGINT(20)"`
	Name        string `json:"name" xorm:"not null default '' comment('名称') VARCHAR(100)"`
	Path        string `json:"path" xorm:"not null default '' comment('路径') index VARCHAR(50)"`
	Component   string `json:"component" xorm:"not null default '' comment('组件') VARCHAR(100)"`
	Redirect    string `json:"redirect" xorm:"not null default '' comment('重定向') VARCHAR(200)"`
	Url         string `json:"url" xorm:"not null default '' comment('url') VARCHAR(200)"`
	MetaTitle   string `json:"meta_title" xorm:"not null default '' comment('meta标题') VARCHAR(50)"`
	MetaIcon    string `json:"meta_icon" xorm:"not null default '' comment('meta icon') VARCHAR(50)"`
	MetaNocache bool   `json:"meta_nocache" xorm:"not null default 0 comment('是否缓存（1:是 0:否）') TINYINT(1)"`
	Alwaysshow  bool   `json:"alwaysshow" xorm:"not null default 0 comment('是否总是显示（1:是0：否）') TINYINT(1)"`
	MetaAffix   int    `json:"meta_affix" xorm:"not null default 0 comment('是否加固（1:是0：否）') TINYINT(1)"`
	Type        int32  `json:"type" xorm:"not null default 2 comment('类型(1:固定,2:权限配置,3特殊)') TINYINT(4)"`
	Hidden      bool   `json:"hidden" xorm:"not null default 0 comment('是否隐藏（0否1是）') TINYINT(1)"`
	Pid         int64  `json:"pid" xorm:"not null default 0 comment('父ID') index(idx_list) INT(11)"`
	Sort        int32  `json:"sort" xorm:"not null default 0 comment('排序') index(idx_list) INT(11)"`
	Status      int32  `json:"status" xorm:"not null default 1 comment('状态（0禁止1启动）') index(idx_list) TINYINT(4)"`
	Level       int32  `json:"level" xorm:"not null default 0 comment('层级') TINYINT(4)"`
	IsLeaf      bool   `json:"isLeaf"  xorm:"not null default 1 comment('是否叶子节点（1:是 0:否）') TINYINT(1)"`
	Domain      string `json:"domain"  xorm:"VARCHAR(255)"`
	Group       string `json:"group" xorm:"VARCHAR(64)"`
	Remark      string `json:"remark" xorm:"VARCHAR(500)"`
	CreatedBy   int64  `json:"created_by" xorm:"BIGINT(20)"`
	CreatedAt   int64  `xorm:"created"`
	UpdatedBy   int64  `json:"updated_by" xorm:"BIGINT(20)"`
	UpdatedAt   int64  `xorm:"updated"`
	Version     int    `xorm:"version"`
}

func (sm *MenuEntity) InsertMenu() (int64, error) {
	lastmenu := new(MenuEntity)
	has, err := db.Orm.Where("`group` like ? ", sm.Group).And("domain like ?", sm.Domain).Desc("id").Limit(1).Get(lastmenu)
	if err != nil {
		return 0, err
	}
	if has && lastmenu != nil && lastmenu.Id > 0 {
		sm.Id = lastmenu.Id + 1
	} else {
		lastdata := new(MenuEntity)
		has, err := db.Orm.Where("id > ? ", 0).Desc("id").Limit(1).Get(lastdata)
		if err != nil {
			return 0, err
		}
		if !has {
			return 0, errors.New(relations.CUS_ERR_4004)
		}
		sm.Id = ((lastdata.Id/100)+1)*100 + 1
	}
	_, err = db.Orm.Insert(sm)
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}
	return sm.Id, nil
}

func (sm *MenuEntity) UpdateMenu() (bool, error) {
	_, err := db.Orm.ID(sm.Id).AllCols().Update(sm)
	if err != nil {
		log.Error(err.Error())
		return false, err
	}

	return true, nil
}

func GetMenuById(id int64) (*MenuEntity, error) {
	sm := new(MenuEntity)
	has, err := db.Orm.ID(id).Get(sm)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if !has {
		return nil, errors.New(relations.CUS_ERR_4004)
	}
	return sm, nil
}

func GetParentMenuByGrpName(grpname string, domain string) (*MenuEntity, error) {
	sm := new(MenuEntity)
	has, err := db.Orm.Where("`group` like ? ", grpname).And("pid = ? ", 0).And("domain like ?", domain).Get(sm)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if !has {
		return nil, errors.New(relations.CUS_ERR_4004)
	}
	return sm, nil
}

// 根据子项字符串，获得列表
func FindMenusByIdArr(ids string) ([]*MenuEntity, error) {
	res := make([]*MenuEntity, 0)
	sql := fmt.Sprintf("select * from menu_entity where id in ( %s ) order by sort ", ids)
	err := db.Orm.SQL(sql).Find(&res)
	if err != nil {
		log.Error(err.Error())

	}
	return res, err
}

func FindMenuByPid(pid int64) ([]*MenuEntity, error) {
	smlist := make([]*MenuEntity, 0)
	err := db.Orm.Where("pid = ? ", pid).OrderBy("sort").Find(&smlist)
	if err != nil {
		log.Error(err.Error())
	}
	return smlist, err
}

func FindMenuByPidAndDomain(pid int64, domain string) ([]*MenuEntity, error) {
	smlist := make([]*MenuEntity, 0)
	err := db.Orm.Where("pid = ? ", pid).And("domain like ?", domain).OrderBy("sort").Find(&smlist)
	if err != nil {
		log.Error(err.Error())
	}
	return smlist, err
}

func DeleteMenuById(id int64) error {
	child, err := FindMenuByPid(id)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if len(child) > 0 {
		return errors.New(relations.CUS_ERR_4005)
	}
	me := new(MenuEntity)
	has, err := db.Orm.ID(id).Unscoped().Get(me)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if has {
		_, err := db.Orm.ID(id).Unscoped().Delete(me)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func InitMenuData() error {
	res := &MenuEntity{
		Id: 1,
	}
	hasdata, err := db.Orm.Exist(res)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if hasdata {
		return nil
	} else {
		menulist := make([]*MenuEntity, 0)
		sysgrp := &MenuEntity{
			Id:          1,
			Name:        "System",
			Path:        "sys",
			Component:   "Layout",
			Redirect:    "noredirect",
			Url:         "",
			MetaTitle:   "系统管理",
			MetaIcon:    "international",
			MetaNocache: true,
			Alwaysshow:  false,
			MetaAffix:   0,
			Type:        2,
			Hidden:      false,
			Pid:         0,
			Sort:        0,
			Status:      0,
			Level:       0,
			IsLeaf:      false,
			Domain:      "back",
			Group:       "sysMgr",
			Remark:      "",
			CreatedBy:   0,
			CreatedAt:   1601022296,
			UpdatedBy:   0,
			UpdatedAt:   1601022296,
			Version:     1,
		}
		resgrp := &MenuEntity{
			Id:          2,
			Name:        "ResMgmt",
			Path:        "res",
			Component:   "resource/resmgmt",
			Redirect:    "",
			Url:         "",
			MetaTitle:   "资源维护",
			MetaIcon:    "international",
			MetaNocache: true,
			Alwaysshow:  false,
			MetaAffix:   0,
			Type:        2,
			Hidden:      false,
			Pid:         1,
			Sort:        0,
			Status:      0,
			Level:       1,
			IsLeaf:      true,
			Domain:      "back",
			Group:       "sysMgr",
			Remark:      "",
			CreatedBy:   0,
			CreatedAt:   1601022296,
			UpdatedBy:   0,
			UpdatedAt:   1601022296,
			Version:     1,
		}
		menugrp := &MenuEntity{
			Id:          3,
			Name:        "MenuMgmt",
			Path:        "menu",
			Component:   "menus/menumgmt",
			Redirect:    "",
			Url:         "",
			MetaTitle:   "菜单维护",
			MetaIcon:    "tree-table",
			MetaNocache: true,
			Alwaysshow:  false,
			MetaAffix:   0,
			Type:        2,
			Hidden:      false,
			Pid:         1,
			Sort:        0,
			Status:      0,
			Level:       1,
			IsLeaf:      true,
			Domain:      "back",
			Group:       "sysMgr",
			Remark:      "",
			CreatedBy:   0,
			CreatedAt:   1601022296,
			UpdatedBy:   0,
			UpdatedAt:   1601022296,
			Version:     1,
		}
		bindgrp := &MenuEntity{
			Id:          4,
			Name:        "BindRes",
			Path:        "bind",
			Component:   "menus/bindRes",
			Redirect:    "",
			Url:         "",
			MetaTitle:   "菜单绑定",
			MetaIcon:    "tree-table",
			MetaNocache: true,
			Alwaysshow:  false,
			MetaAffix:   0,
			Type:        2,
			Hidden:      false,
			Pid:         1,
			Sort:        0,
			Status:      0,
			Level:       1,
			IsLeaf:      true,
			Domain:      "back",
			Group:       "sysMgr",
			Remark:      "",
			CreatedBy:   0,
			CreatedAt:   1601022296,
			UpdatedBy:   0,
			UpdatedAt:   1601022296,
			Version:     1,
		}
		dictgrp := &MenuEntity{
			Id:          5,
			Name:        "DictCateMgmt",
			Path:        "mgmt",
			Component:   "dict/list",
			Redirect:    "",
			Url:         "",
			MetaTitle:   "字典维护",
			MetaIcon:    "tree-table",
			MetaNocache: true,
			Alwaysshow:  false,
			MetaAffix:   0,
			Type:        2,
			Hidden:      false,
			Pid:         1,
			Sort:        0,
			Status:      0,
			Level:       1,
			IsLeaf:      true,
			Domain:      "back",
			Group:       "sysMgr",
			Remark:      "",
			CreatedBy:   0,
			CreatedAt:   1601022296,
			UpdatedBy:   0,
			UpdatedAt:   1601022296,
			Version:     1,
		}

		rbacgrp := &MenuEntity{
			Id:          101,
			Name:        "RBAC",
			Path:        "rbac",
			Component:   "Layout",
			Redirect:    "noredirect",
			Url:         "",
			MetaTitle:   "授权管理",
			MetaIcon:    "password",
			MetaNocache: true,
			Alwaysshow:  false,
			MetaAffix:   0,
			Type:        2,
			Hidden:      false,
			Pid:         0,
			Sort:        0,
			Status:      0,
			Level:       0,
			IsLeaf:      false,
			Domain:      "back",
			Group:       "rbacMgr",
			Remark:      "",
			CreatedBy:   0,
			CreatedAt:   1601022296,
			UpdatedBy:   0,
			UpdatedAt:   1601022296,
			Version:     1,
		}

		rbac1 := &MenuEntity{
			Id:          102,
			Name:        "Perms",
			Path:        "perms",
			Component:   "rbac/perms",
			Redirect:    "",
			Url:         "",
			MetaTitle:   "权限维护",
			MetaIcon:    "",
			MetaNocache: true,
			Alwaysshow:  false,
			MetaAffix:   0,
			Type:        2,
			Hidden:      false,
			Pid:         101,
			Sort:        0,
			Status:      0,
			Level:       1,
			IsLeaf:      true,
			Domain:      "back",
			Group:       "rbacMgr",
			Remark:      "",
			CreatedBy:   0,
			CreatedAt:   1601022296,
			UpdatedBy:   0,
			UpdatedAt:   1601022296,
			Version:     1,
		}
		rbac2 := &MenuEntity{
			Id:          103,
			Name:        "Roles",
			Path:        "roles",
			Component:   "rbac/roles",
			Redirect:    "",
			Url:         "",
			MetaTitle:   "角色维护",
			MetaIcon:    "",
			MetaNocache: true,
			Alwaysshow:  false,
			MetaAffix:   0,
			Type:        2,
			Hidden:      false,
			Pid:         101,
			Sort:        0,
			Status:      0,
			Level:       1,
			IsLeaf:      true,
			Domain:      "back",
			Group:       "rbacMgr",
			Remark:      "",
			CreatedBy:   0,
			CreatedAt:   1601022296,
			UpdatedBy:   0,
			UpdatedAt:   1601022296,
			Version:     1,
		}
		rbac3 := &MenuEntity{
			Id:          104,
			Name:        "CreatePerm",
			Path:        "perm/create",
			Component:   "rbac/createPerm/:domain(\\w+)",
			Redirect:    "",
			Url:         "",
			MetaTitle:   "创建权限",
			MetaIcon:    "",
			MetaNocache: true,
			Alwaysshow:  false,
			MetaAffix:   0,
			Type:        2,
			Hidden:      true,
			Pid:         101,
			Sort:        0,
			Status:      0,
			Level:       1,
			IsLeaf:      true,
			Domain:      "back",
			Group:       "rbacMgr",
			Remark:      "",
			CreatedBy:   0,
			CreatedAt:   1601022296,
			UpdatedBy:   0,
			UpdatedAt:   1601022296,
			Version:     1,
		}
		rbac4 := &MenuEntity{
			Id:          105,
			Name:        "EditPerm",
			Path:        "perm/edit/:domain(\\w+)/:id(\\d+)",
			Component:   "rbac/editPerm",
			Redirect:    "",
			Url:         "",
			MetaTitle:   "修改权限",
			MetaIcon:    "",
			MetaNocache: true,
			Alwaysshow:  false,
			MetaAffix:   0,
			Type:        2,
			Hidden:      true,
			Pid:         101,
			Sort:        0,
			Status:      0,
			Level:       1,
			IsLeaf:      true,
			Domain:      "back",
			Group:       "rbacMgr",
			Remark:      "",
			CreatedBy:   0,
			CreatedAt:   1601022296,
			UpdatedBy:   0,
			UpdatedAt:   1601022296,
			Version:     1,
		}
		rbac5 := &MenuEntity{
			Id:          106,
			Name:        "CreateRole",
			Path:        "role/create/:domain(\\w+)",
			Component:   "rbac/createRole",
			Redirect:    "",
			Url:         "",
			MetaTitle:   "创建角色",
			MetaIcon:    "",
			MetaNocache: true,
			Alwaysshow:  false,
			MetaAffix:   0,
			Type:        2,
			Hidden:      true,
			Pid:         101,
			Sort:        0,
			Status:      0,
			Level:       1,
			IsLeaf:      true,
			Domain:      "back",
			Group:       "rbacMgr",
			Remark:      "",
			CreatedBy:   0,
			CreatedAt:   1601022296,
			UpdatedBy:   0,
			UpdatedAt:   1601022296,
			Version:     1,
		}
		rbac6 := &MenuEntity{
			Id:          107,
			Name:        "EditRole",
			Path:        "role/edit/:domain(\\w+)/:id(\\d+)",
			Component:   "rbac/editRole",
			Redirect:    "",
			Url:         "",
			MetaTitle:   "修改角色",
			MetaIcon:    "",
			MetaNocache: true,
			Alwaysshow:  false,
			MetaAffix:   0,
			Type:        2,
			Hidden:      true,
			Pid:         101,
			Sort:        0,
			Status:      0,
			Level:       1,
			IsLeaf:      true,
			Domain:      "back",
			Group:       "rbacMgr",
			Remark:      "",
			CreatedBy:   0,
			CreatedAt:   1601022296,
			UpdatedBy:   0,
			UpdatedAt:   1601022296,
			Version:     1,
		}

		menulist = append(append(append(append(append(menulist, sysgrp), resgrp), menugrp), bindgrp), dictgrp)
		menulist = append(append(append(append(append(append(append(menulist, rbacgrp), rbac1), rbac2), rbac3), rbac4), rbac5), rbac6)

		_, err := db.Orm.Insert(&menulist)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}
