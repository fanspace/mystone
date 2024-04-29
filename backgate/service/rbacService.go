package service

import (
	log "backgate/logger"
	"backgate/rbac"
)

func UpdatePerms2Casbin() error {
	isdel, err := rbac.Casbin.DeleteRoleForUserInDomain("tmpuser2", "tmprole2", "back")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if !isdel {

		rbac.Casbin.AddRoleForUserInDomain("tmpuser2", "tmprole2", "back")
	}
	rbac.Casbin.LoadPolicy()
	return nil
}
