package model

import (
	"backend/db"
	"backend/relations"
	"backend/utils"

	log "backend/logger"
	"errors"
	"fmt"
)

type Resource struct {
	Id        int64       `json:"id" xorm:"pk autoincr BIGINT(20)"`
	Name      string      `json:"name"  xorm:"not null VARCHAR(64) INDEX(grp_name_unique_idx)"`
	NameCn    string      `json:"nameCn" xorm:"VARCHAR(255)"`
	Url       string      `json:"url"  xorm:"not null VARCHAR(255)"`
	Act       string      `json:"act"  xorm:"not null VARCHAR(32)"`
	Pid       int64       `json:"pid" xorm:"not null BIGINT(20)"`
	IsLeaf    bool        `json:"isLeaf" xorm:"not null default 1 TINYINT(1)"`
	Domain    string      `json:"domain"  xorm:"VARCHAR(255) INDEX(grp_name_unique_idx)" `
	Remark    string      `json:"remark"  xorm:"VARCHAR(255)"`
	GroupName string      `json:"groupName" xorm:"not null VARCHAR(64) INDEX(grp_name_unique_idx)"`
	Level     int32       `json:"level" xorm:"not null default 0 comment('层级') TINYINT(4)"`
	CreatedBy int64       `json:"created_by" xorm:"BIGINT(20)"`
	CreatedAt int64       `xorm:"created"`
	UpdatedBy int64       `json:"updated_by" xorm:"BIGINT(20)"`
	UpdatedAt int64       `xorm:"updated"`
	Version   int         `xorm:"version"`
	Children  []*Resource `json:"children" xorm:"-"`
}

func (zr *Resource) HasResExist() (int64, error) {
	has, err := db.Orm.Limit(0).Get(zr)
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}
	if !has {
		return 0, nil
	}

	return zr.Id, nil
}

func (zr *Resource) InsertResorce() (int64, error) {
	lastzr := new(Resource)
	has, err := db.Orm.Where("`group_name` like ? ", zr.GroupName).And("domain like ?", zr.Domain).Desc("id").Limit(1).Get(lastzr)
	if err != nil {
		return 0, err
	}
	if has && lastzr != nil && lastzr.Id > 0 {
		zr.Id = lastzr.Id + 1
	} else {
		lastdata := new(Resource)
		has, err := db.Orm.Where("id > ? ", 0).Desc("id").Limit(1).Get(lastdata)
		if err != nil {
			return 0, err
		}
		if !has {
			zr.Id = 101
		}
		zr.Id = ((lastdata.Id/100)+1)*100 + 1
	}
	if zr.Level > 1 && zr.Pid == 0 {
		grpitem := new(Resource)
		hasgrp, err := db.Orm.Where("`group_name` like ? ", zr.GroupName).And("domain like ?", zr.Domain).And("level = ?", 1).Desc("id").Limit(1).Get(grpitem)
		if err != nil {
			log.Error(err.Error())
			return 0, err
		}
		if !hasgrp {
			return 0, errors.New(relations.CUS_ERR_4004)
		}
		zr.Pid = grpitem.Id
	}

	_, err = db.Orm.Insert(zr)
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}

	return zr.Id, nil
}

func (zr *Resource) UpdateResource() (bool, error) {
	_, err := db.Orm.ID(zr.Id).AllCols().Update(zr)
	if err != nil {
		log.Error(err.Error())
		return false, err
	}

	return true, nil
}

func DeleteResourceById(id int64) error {
	child, err := FindResourcesByPid(id)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if len(child) > 0 {
		return errors.New(relations.CUS_ERR_4005)
	}
	me := new(Resource)
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

func GetResourceById(id int64) (*Resource, error) {
	zr := new(Resource)
	has, err := db.Orm.ID(id).Get(zr)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if !has {
		return nil, errors.New(relations.CUS_ERR_4004)
	}
	return zr, nil
}

//  某id下的子资源

func FindResourcesByPid(pid int64) ([]*Resource, error) {
	smlist := make([]*Resource, 0)
	err := db.Orm.Where("pid = ? ", pid).Find(&smlist)
	if err != nil {
		log.Error(err.Error())
	}
	return smlist, err
}

// 根据gourpname  获得 该类主资源
func GetParentZjResourcesByGrpName(grpname, domain string) (*Resource, error) {
	sm := new(Resource)
	has, err := db.Orm.Where(" `group_name` like ? ", grpname).And("domain like ?", domain).And("pid = ? ", 0).Get(sm)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if !has {
		return nil, errors.New(relations.CUS_ERR_4004)
	}
	return sm, nil
}

// 根据pid得所有
func FindAllResourcesByPid(pid int64, domain string) ([]*Resource, error) {
	reslist := make([]*Resource, 0)
	err := db.Orm.Where("pid = ? ", pid).And("domain like ?", domain).Find(&reslist)
	if err != nil {
		log.Error(err.Error())
	}
	return reslist, err
}

// 根据子项字符串，获得列表
func FindResourcesByIdArr(ids string) ([]*Resource, error) {
	res := make([]*Resource, 0)
	sql := fmt.Sprintf("select * from resource where id in ( %s ) ", ids)
	err := db.Orm.SQL(sql).Find(&res)
	if err != nil {
		log.Error(err.Error())

	}
	return res, err
}

func QueryMenuStrByUsid(usid int) (string, error) {
	result, err := db.Orm.Query("select GROUP_CONCAT(menus) as menus from user_menu where usid = ?", usid)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	menus := string(result[0]["menus"])
	if menus == "" {
		return "", nil
	}
	if !utils.IsNumIds(menus) {
		return "", errors.New(relations.CUS_ERR_3011)
	}
	return menus, nil
}

func QueryResByUsid(usid int) ([]*Resource, error) {
	res := make([]*Resource, 0)
	menus, err := QueryMenuStrByUsid(usid)
	if err != nil {
		return res, err
	}
	if menus == "" {
		return res, errors.New(relations.CUS_ERR_4004)
	}
	err = db.Orm.SQL("select * from resource a,(select res_id from menu_res where menu_id in (" + menus + ") group by res_id )b  where a.id = b.res_id").Find(&res)
	if err != nil {
		log.Error(err.Error())
		return res, err
	}

	return res, nil
}

func QueryResByConditions(req *Resource) (int64, []*Resource, error) {
	res := make([]*Resource, 0)
	session := db.Orm.Where("1=1")
	if req.GroupName != "" {
		session = session.And("group_name = ?", req.GroupName)
	}
	if req.Level != 0 {
		session = session.And("level = ?", req.Level)
	}

	if req.Domain != "" {
		session = session.And("domain = ?", req.Domain)
	}
	if req.Pid != -1 {
		session = session.And("pid = ?", req.Pid)
	}
	total, err := session.OrderBy("created_at desc").FindAndCount(&res)
	if err != nil {
		log.Error(err.Error())
		return 0, res, err
	}
	return total, res, nil
}

// 初始化资源
func InitResData() error {
	res := &Resource{
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
		reslist := make([]*Resource, 0)
		resgrp := &Resource{
			Id:        1,
			Name:      "资源管理",
			Url:       "/*",
			Act:       "*",
			Pid:       0,
			IsLeaf:    false,
			Domain:    "back",
			Remark:    "/sys/res/entry",
			GroupName: "resMgr",
			Level:     0,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}

		resitem1 := &Resource{
			Id:        2,
			Name:      "查询资源列表",
			Url:       "/list",
			Act:       "POST",
			Pid:       1,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/sys/res/entry",
			GroupName: "resMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}

		resitem2 := &Resource{
			Id:        3,
			Name:      "查询资源详情",
			Url:       "/query",
			Act:       "POST",
			Pid:       1,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/sys/res/entry",
			GroupName: "resMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}

		resitem3 := &Resource{
			Id:        4,
			Name:      "编辑资源",
			Url:       "/edit",
			Act:       "POST",
			Pid:       1,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/sys/res/entry",
			GroupName: "resMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}
		resitem4 := &Resource{
			Id:        5,
			Name:      "查询绑定的资源",
			Url:       "/find/:menuid",
			Act:       "GET",
			Pid:       1,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/sys/res/entry",
			GroupName: "resMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}
		resitem5 := &Resource{
			Id:        6,
			Name:      "绑定资源到菜单",
			Url:       "/bind",
			Act:       "POST",
			Pid:       1,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/sys/res/entry",
			GroupName: "resMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}

		menugrp := &Resource{
			Id:        101,
			Name:      "菜单管理",
			Url:       "/*",
			Act:       "*",
			Pid:       0,
			IsLeaf:    false,
			Domain:    "back",
			Remark:    "/sys/res/menu",
			GroupName: "menuMgr",
			Level:     0,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}
		menuitem1 := &Resource{
			Id:        102,
			Name:      "查询菜单列表",
			Url:       "/list",
			Act:       "POST",
			Pid:       101,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/sys/res/menu",
			GroupName: "menuMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}

		menuitem2 := &Resource{
			Id:        103,
			Name:      "查询菜单详情",
			Url:       "/query",
			Act:       "POST",
			Pid:       101,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/sys/res/menu",
			GroupName: "menuMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}

		menuitem3 := &Resource{
			Id:        104,
			Name:      "编辑菜单",
			Url:       "/edit",
			Act:       "POST",
			Pid:       101,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/sys/res/menu",
			GroupName: "menuMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}

		dictgrp := &Resource{
			Id:        201,
			Name:      "字典管理",
			Url:       "/*",
			Act:       "*",
			Pid:       0,
			IsLeaf:    false,
			Domain:    "back",
			Remark:    "/sys/res/dict",
			GroupName: "dictMgr",
			Level:     0,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}
		dictitem1 := &Resource{
			Id:        202,
			Name:      "查询字典分类",
			Url:       "/cate/:keyword",
			Act:       "GET",
			Pid:       201,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/sys/res/dict",
			GroupName: "dictMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}
		dictitem2 := &Resource{
			Id:        203,
			Name:      "编辑字典分类",
			Url:       "/cate/edit",
			Act:       "POST",
			Pid:       201,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/sys/res/dict",
			GroupName: "dictMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}
		dictitem3 := &Resource{
			Id:        204,
			Name:      "编辑字典子项",
			Url:       "/item/edit",
			Act:       "POST",
			Pid:       201,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/sys/res/dict",
			GroupName: "dictMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}

		rbacgrp := &Resource{
			Id:        301,
			Name:      "权限管理",
			Url:       "/*",
			Act:       "*",
			Pid:       0,
			IsLeaf:    false,
			Domain:    "back",
			Remark:    "/acc/rbac",
			GroupName: "rbacMgr",
			Level:     0,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}
		rbacitem1 := &Resource{
			Id:        302,
			Name:      "编辑权限",
			Url:       "/perm/edit",
			Act:       "POST",
			Pid:       301,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/acc/rbac",
			GroupName: "rbacMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}
		rbacitem2 := &Resource{
			Id:        303,
			Name:      "查询权限列表",
			Url:       "/perm/list",
			Act:       "POST",
			Pid:       301,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/acc/rbac",
			GroupName: "rbacMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}
		rbacitem3 := &Resource{
			Id:        304,
			Name:      "查询权限",
			Url:       "/perm/query/:id",
			Act:       "GET",
			Pid:       301,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/acc/rbac",
			GroupName: "rbacMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}
		rbacitem4 := &Resource{
			Id:        305,
			Name:      "查询角色列表",
			Url:       "/role/list",
			Act:       "POST",
			Pid:       301,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/acc/rbac",
			GroupName: "rbacMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}
		rbacitem5 := &Resource{
			Id:        306,
			Name:      "查询角色",
			Url:       "/role/query/:id",
			Act:       "GET",
			Pid:       301,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/acc/rbac",
			GroupName: "rbacMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}
		rbacitem6 := &Resource{
			Id:        307,
			Name:      "编辑角色",
			Url:       "/role/edit",
			Act:       "POST",
			Pid:       301,
			IsLeaf:    true,
			Domain:    "back",
			Remark:    "/acc/rbac",
			GroupName: "rbacMgr",
			Level:     1,
			CreatedBy: 0,
			CreatedAt: 1600838950,
			UpdatedBy: 0,
			UpdatedAt: 1600838950,
			Version:   1,
		}
		reslist = append(append(append(append(append(append(reslist, resgrp), resitem1), resitem2), resitem3), resitem4), resitem5)
		reslist = append(append(append(append(reslist, menugrp), menuitem1), menuitem2), menuitem3)
		reslist = append(append(append(append(reslist, dictgrp), dictitem1), dictitem2), dictitem3)
		reslist = append(append(append(append(append(append(append(reslist, rbacgrp), rbacitem1), rbacitem2), rbacitem3), rbacitem4), rbacitem5), rbacitem6)
		_, err := db.Orm.Insert(&reslist)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}
