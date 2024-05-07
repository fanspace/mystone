package model

import (
	"backend/db"
	log "backend/logger"
	"errors"
)

type Roles struct {
	Id         int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	RoleName   string `json:"roleName" xorm:"not null comment('角色名称') VARCHAR(64) unique('role_domain')"`
	RoleNameCn string `json:"roleNameCn" xorm:"not null comment('权限名称中文') VARCHAR(64)"`
	Perms      string `json:"perms" xorm:"not null comment('权限列表') VARCHAR(255)"`
	PermsName  string `json:"permsName" xorm:"not null comment('权限列表中文') VARCHAR(512)"`
	Domain     string `json:"domain" xorm:"comment('域') VARCHAR(64)  unique('role_domain')"` // 保留
	Remark     string `json:"remark" xorm:"comment('备注') VARCHAR(255)"`
	CreatedBy  int64  `json:"created_by" xorm:"BIGINT(20)"`
	CreatedAt  int64  `xorm:"created"`
	UpdatedBy  int64  `json:"updated_by" xorm:"BIGINT(20)"`
	UpdatedAt  int64  `xorm:"updated"`
	Version    int    `xorm:"version"`
}

/**
 *  @MethodName AddSysRoles
 *  @author: lf6128@163.com
 *  @Date: 2022/6/8 16:16
 *  @Description: Add  sysRoles
 */

func (sysRoles *Roles) AddSysRoles() (*Roles, error) {
	_, err := db.Orm.Insert(sysRoles)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return sysRoles, nil
}

/**
 *  @MethodName UpdateSysRoles
 *  @author: lf6128@163.com
 *  @Date: 2022/6/8 16:16
 *  @Description: Update  sysRoles
 */

func (sysRoles *Roles) UpdateSysRoles() error {
	_, err := db.Orm.ID(sysRoles.Id).Update(sysRoles)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

/**
 *  @MethodName StatSysRolesById
 *  @author: lf6128@163.com
 *  @Date: 2022/6/8 16:16
 *  @Description: Update  sysRoles  status
 */

func StatSysRolesById(id int64, stat int32) error {
	sql := "update sys_roles set status = ? where id = ? "
	_, err := db.Orm.Exec(sql, stat, id)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

/**
 *  @MethodName DeleteSysRolesById
 *  @author: lf6128@163.com
 *  @Date: 2022/6/8 16:16
 *  @Description: Delete  sysRoles  Physically
 */

func DeleteSysRolesById(id int64) error {
	sysRoles := new(Roles)
	has, err := db.Orm.ID(id).Unscoped().Get(sysRoles)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if has {
		_, err := db.Orm.ID(id).Unscoped().Delete(sysRoles)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

/**
 *  @MethodName GetSysRolesById
 *  @author: lf6128@163.com
 *  @Date: 2022/6/8 16:16
 *  @Description: Get  sysRoles  ById
 */

func GetSysRolesById(id int64) (*Roles, error) {
	res := new(Roles)
	has, err := db.Orm.ID(id).Get(res)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if !has {
		return nil, errors.New("data not exist")
	}
	return res, nil
}

/**
 *  @MethodName ListSysRoles
 *  @author: lf6128@163.com
 *  @Date: 2022/6/8 16:16
 *  @Description: List  sysRoles
 */
func FindSysRolesByInstance(sysRoles *Roles, page int, pagesize int) (int64, []*Roles, error) {
	res := make([]*Roles, 0)
	total, err := db.Orm.Count(sysRoles)
	if err != nil {
		log.Error(err.Error())
		return 0, res, err
	}
	if total > 0 {
		err := db.Orm.Where("id > ?", 0).Limit(pagesize, (page-1)*pagesize).Find(&res)
		if err != nil {
			log.Error(err.Error())
			return 0, res, err
		}
	}
	return total, res, nil
}

func GetSysRolesByNameAndDomain(rolename string, domain string) (*Roles, error) {
	res := new(Roles)
	has, err := db.Orm.Where("role_name like ?", rolename).And("domain like ?", domain).Get(res)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if !has {
		return nil, errors.New("data not exist")
	}
	return res, nil
}
