package model

import (
	"backend/db"
	log "backend/logger"
	"backend/relations"
	"errors"
	"fmt"
)

type Perms struct {
	Id         int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	PermName   string `json:"permName" xorm:"not null comment('权限名称') VARCHAR(64)  unique('perm_domain')"`
	PermNameCn string `json:"permNameCn" xorm:"not null comment('权限名称中文') VARCHAR(64)"` // 例如  账户查询权限(只读）、  账户管理权限（全，包含写，不再另设）、
	Menus      string `json:"menus" xorm:"not null comment('菜单列表') VARCHAR(255)"`       //  注意，菜单与资源不同，父子必须同时存在，本菜单列表中，仅列举子菜单，需要与group联合出菜单
	MenuNames  string `json:"menuNames" xorm:"not null comment('菜单列表中文') VARCHAR(255)"`
	Group      string `json:"group" xorm:"not null comment('群组名称') VARCHAR(64)"` //菜单项必须用到
	GroupId    int64  `json:"group_id"`                                          // 这一版本用不到，暂时保留
	IsAll      bool   `json:"isAll"  xorm:"default 0 BOOL"`                      // 是否有全部权限, 这一版本用不到，暂时保留
	Domain     string `json:"domain" xorm:"comment('域') VARCHAR(64)  unique('perm_domain')"`
	Remark     string `json:"remark" xorm:"comment('备注') VARCHAR(255)"`
	CreatedBy  int64  `json:"created_by" xorm:"BIGINT(20)"`
	CreatedAt  int64  `xorm:"created"`
	UpdatedBy  int64  `json:"updated_by" xorm:"BIGINT(20)"`
	UpdatedAt  int64  `xorm:"updated"`
	Version    int    `xorm:"version"`
}

func (zr *Perms) InsertPerms() (int64, error) {
	_, err := db.Orm.Insert(zr)
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}
	return zr.Id, nil
}

func (zr *Perms) UpdatePerms() (bool, error) {
	_, err := db.Orm.ID(zr.Id).AllCols().Update(zr)
	if err != nil {
		log.Error(err.Error())
		return false, err
	}

	return true, nil
}

func DeletePermsById(id int64) error {
	me := new(Perms)
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

func GetPermsById(id int64) (*Perms, error) {
	zr := new(Perms)
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

func GetPermsByName(name string) (*Perms, error) {
	zr := new(Perms)
	has, err := db.Orm.Where("perm_name like ? ", name).Get(zr)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if !has {
		return nil, errors.New(relations.CUS_ERR_4004)
	}
	return zr, nil
}

func FindPermsByIds(ids string) ([]*Perms, error) {
	res := make([]*Perms, 0)
	sql := fmt.Sprintf("select * from perms where id in ( %s ) ", ids)
	err := db.Orm.SQL(sql).Find(&res)
	if err != nil {
		log.Error(err.Error())

	}
	return res, err
}

// 查询组名，这里暂时不用，使用配置来替代
func FindGroupsFromPerms() ([]*Perms, error) {
	res := make([]*Perms, 0)
	err := db.Orm.Distinct("group").Find(&res)
	if err != nil {
		log.Error(err.Error())
	}
	return res, err
}

// 根据子项字符串，获得列表
func FindPermsByIdArr(ids string) ([]*Perms, error) {
	res := make([]*Perms, 0)
	sql := fmt.Sprintf("select * from perms where id in ( %s ) ", ids)
	err := db.Orm.SQL(sql).Find(&res)
	if err != nil {
		log.Error(err.Error())

	}
	return res, err
}

// 查询包含menuid的perms,
func QueryPermsByMenuId(mid int64, domain string) ([]*Perms, error) {
	res := make([]*Perms, 0)
	sql := fmt.Sprintf(" select * from perms where domain like ? and find_in_set(%d, menus)", mid)
	err := db.Orm.SQL(sql, domain).Find(&res)
	if err != nil {
		log.Error(err.Error())

	}
	return res, err
}
