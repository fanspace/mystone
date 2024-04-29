package model

import (
	"backos/db"
	log "backos/logger"
	"backos/relations"
	"errors"
	"time"
)

type Org struct {
	Id           int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	UnitName     string `json:"unit_name" xorm:"unique not null comment('单位名称中文') VARCHAR(64)"`
	NamePy       string `json:"name_py" xorm:"not null comment('单位名称拼音') VARCHAR(255)"`
	UnitIndustry string `json:"unit_industry" xorm:"not null comment('行业') VARCHAR(500)"`
	UnitCharater int32  `json:"unit_charater" xorm:"not null default 0 comment('单位性质') INT(10)"`
	Status       int32  `json:"status" xorm:"not null default 1 INT(11)"` // 插入1   审核2
	Province     int32  `json:"province" xorm:"comment('省') INT(10)"`
	City         int32  `json:"city" xorm:"comment('市') INT(10)"`
	District     int32  `json:"district" xorm:"comment('区') INT(10)"`
	Location     string `json:"location" xorm:"comment('省市区中文') VARCHAR(255)"` //所在地区	省市区 的字符串
	Address      string `json:"address"  xorm:"comment('地址：道路门牌号') VARCHAR(255)"`
	CreditCode   string `json:"credit_code" xorm:"unique comment('信用号') VARCHAR(64)"`
	LinkedMan    string `json:"linked_man" xorm:"not null comment('单位负责人') VARCHAR(64)"`
	LinkedMobile string `json:"linked_mobile" xorm:"not null comment('单位负责人手机') VARCHAR(64)"`
	LinkedPhone  string `json:"linked_phone" xorm:"not null comment('单位负责人电话') VARCHAR(64)"`
	LinkedEmail  string `json:"linked_email" xorm:"not null comment('单位负责人邮箱') VARCHAR(64)"`
	CreatedBy    int64  `json:"created_by" xorm:"BIGINT(20)"`
	CreatedAt    int64  `xorm:"created"`
	UpdatedBy    int64  `json:"updatedBy"  xorm:"BIGINT(20)"`
	UpdatedAt    int64  `json:"updatedAt"  xorm:"BIGINT(20)"`
	IsHidden     bool   `json:"isHidden" xorm:"BOOL"`
	Version      int    `xorm:"version"`
	Cover        string `json:"cover" xorm:"comment('单位标识') VARCHAR(255)"`
	Icon         string `json:"icon" xorm:"comment('单位icon') VARCHAR(255)"`
	Domain       string `json:"domain" xorm:"VARCHAR(64)"`
	Remark       string `json:"remark" xorm:"comment('备注') VARCHAR(255)"`
}

/**
 *  @MethodName AddOrg
 *  @author: lf6128@163.com
 *  @Date: 2022/8/8 16:01
 *  @Description: Add  org
 */

func (org *Org) AddOrg() (*Org, error) {
	_, err := db.Orm.Insert(org)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return org, nil
}

/**
 *  @MethodName UpdateOrg
 *  @author: lf6128@163.com
 *  @Date: 2022/8/8 16:01
 *  @Description: Update  org
 */

func (org *Org) UpdateOrg() error {
	_, err := db.Orm.ID(org.Id).AllCols().Update(org)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

/**
 *  @MethodName StatOrgById
 *  @author: lf6128@163.com
 *  @Date: 2022/8/8 16:01
 *  @Description: Update  org  status
 */

func StatOrgById(id int64, stat int32, operator int64, version int) error {

	sql := "update org set status = ?, updated_at = ?, updated_by = ?, version = ? where id = ?  and  version = ? "
	_, err := db.Orm.Exec(sql, stat, time.Now().Unix(), operator, version+1, id, version)
	if err != nil {
		log.Error(err.Error())
	}
	db.Orm.ClearCache(new(Org))
	return err
}

/**
 *  @MethodName DeleteOrgById
 *  @author: lf6128@163.com
 *  @Date: 2022/8/8 16:01
 *  @Description: Delete  org  Physically
 */

func DeleteOrgById(id int64) error {
	org := new(Org)
	has, err := db.Orm.ID(id).Unscoped().Get(org)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if has {
		_, err := db.Orm.ID(id).Unscoped().Delete(org)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

/**
 *  @MethodName GetOrgById
 *  @author: lf6128@163.com
 *  @Date: 2022/8/8 16:01
 *  @Description: Get  org  ById
 */

func GetOrgById(id int64) (*Org, error) {
	res := new(Org)
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

func IsOrgExist(query string, qtype int32) error {
	sql := "select count(*) from org  where id > 0 "
	switch qtype {
	case 1:
		sql += " and unit_name like ? "
	case 2:
		sql += " and credit_code like ? "
	default:
		return errors.New(relations.CUS_ERR_4002)
	}
	org := new(Org)
	total, err := db.Orm.SQL(sql, query).Count(org)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if total > 0 {
		return errors.New("duplicate data found")
	}
	return nil
}

/**
 *  @MethodName ListOrg
 *  @author: lf6128@163.com
 *  @Date: 2022/8/8 16:01
 *  @Description: List  org
 */
func FindOrgByInstance(org *Org, page int, pagesize int) (int64, []*Org, error) {
	res := make([]*Org, 0)
	total, err := db.Orm.Count(org)
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
