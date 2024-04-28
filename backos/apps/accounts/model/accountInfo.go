package model

import (
	"backos/db"
	log "backos/logger"
	"errors"
)

/**
 *  @Classname AccountInfo
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 15:02
 *  @Description:
 */

type AccountInfo struct {
	Id int64 `json:"id" xorm:" pk  BIGINT(20) not null unique"`
	// 是否实名
	IsIdValid bool   `json:"is_id_valid" xorm:"Bool"`
	Realname  string `json:"realname" xorm:"not null VARCHAR(32)"`
	Idtype    int32  `json:"idtype" xorm:"INT(10) comment('证件类型') unique('id_agent')"`
	Idcode    string `json:"idcode" xorm:"VARCHAR(32) comment('证件号') unique('id_agent')"`
	// json 自定义字段
	AccountDetail map[string]interface{} `json:"account_detail" xorm:"MediumText JSON 'account_detail'"`
	UpdatedBy     int64                  `json:"updatedBy"  xorm:"BIGINT(20)"`
	UpdatedAt     int64                  `json:"updatedAt"  xorm:"BIGINT(20)"`
	Version       int                    `json:"version"  xorm:"version"`
}

/**
 *  @MethodName TableName
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 15:08
 *  @Description: return TableName
 */

func (accountInfo *AccountInfo) TableName() string {
	return "account_info"
}

/**
 *  @MethodName NewAccountInfo
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 15:08
 *  @Description: New  accountInfo
 */

func (accountInfo *AccountInfo) NewAccountInfo() (*AccountInfo, error) {
	_, err := db.Orm.Insert(accountInfo)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return accountInfo, nil
}

/**
 *  @MethodName UpdateAccountInfo
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 15:08
 *  @Description: Update  accountInfo
 */

func (accountInfo *AccountInfo) UpdateAccountInfo() error {
	_, err := db.Orm.ID(accountInfo.Id).Update(accountInfo)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

/**
 *  @MethodName StatAccountInfoById
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 15:08
 *  @Description: Update  accountInfo  status
 */

func StatAccountInfoById(id int64, stat int32) error {
	sql := "update account_info set status = ? where id = ? "
	_, err := db.Orm.Exec(sql, stat, id)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

/**
 *  @MethodName RemoveAccountInfoById
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 15:08
 *  @Description: Logical Delete accountInfo (),  set status = -99 (deleted)  -1(invalid)
 */

func RemoveAccountInfoById(id int64) error {
	sql := "update account_info set status = -99 where id = ? "
	_, err := db.Orm.Exec(sql, id)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

/**
 *  @MethodName DeleteAccountInfoById
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 15:08
 *  @Description: Physical Delete  accountInfo  Physically
 */

func DeleteAccountInfoById(id int64) error {
	accountInfo := new(AccountInfo)
	has, err := db.Orm.ID(id).Unscoped().Get(accountInfo)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if has {
		_, err := db.Orm.ID(id).Unscoped().Delete(accountInfo)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

/**
 *  @MethodName GetAccountInfoById
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 15:08
 *  @Description: Get  accountInfo  ById
 */

func GetAccountInfoById(id int64) (*AccountInfo, error) {
	res := new(AccountInfo)
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
 *  @MethodName ListAccountInfo
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 15:08
 *  @Description: List  accountInfo
 */
func FindAccountInfoByInstance(accountInfo *AccountInfo, page int, pagesize int) (int64, []*AccountInfo, error) {
	res := make([]*AccountInfo, 0)
	total, err := db.Orm.Count(accountInfo)
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
