package model

import (
	"backos/db"
	log "backos/logger"
	"backos/relations"
	"errors"
)

/**
 *  @Classname ACCOUNT_GO
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 14:26
 *  @Description:
 */

type Account struct {
	Id           int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	Username     string `json:"username" xorm:"unique not null VARCHAR(64)"`
	Password     string `json:"password" xorm:"not null VARCHAR(64)"`
	PasswordHash string `xorm:"not null VARCHAR(32)"`
	// 用户类型  /guest/normal/company/manager
	UserType int32 `json:"usertype" xorm:"unique(mobile_type) not null default 0 INT(8)"`
	Status   int32 `json:"status" xorm:"not null default 0 TINYINT(3)"`
	//
	Domaim    string `json:"Domain" xorm:"VARCHAR(32)"`
	Mobile    string `json:"mobile" xorm:"unique(mobile_type) not null VARCHAR(12)"`
	Email     string `json:"email" xorm:"VARCHAR(255)"`
	CreatedBy int64  `json:"created_by" xorm:"BIGINT(20)"`
	CreatedAt int64  `xorm:"created_at"`
	Version   int    `xorm:"version"`
}

/**
 *  @MethodName TableName
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 14:52
 *  @Description: return TableName
 */

func (account *Account) TableName() string {
	return "account"
}

/**
 *  @MethodName NewAccount
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 14:52
 *  @Description: New  account
 */

func (account *Account) NewAccount() (*Account, error) {
	_, err := db.Orm.Insert(account)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return account, nil
}

/**
 *  @MethodName UpdateAccount
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 14:52
 *  @Description: Update  account
 */

func (account *Account) UpdateAccount() error {
	_, err := db.Orm.ID(account.Id).Update(account)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

/**
 *  @MethodName StatAccountById
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 14:52
 *  @Description: Update  account  status
 */

func StatAccountById(id int64, stat int32) error {
	sql := "update account set status = ? where id = ? "
	_, err := db.Orm.Exec(sql, stat, id)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

/**
 *  @MethodName RemoveAccountById
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 14:52
 *  @Description: Logical Delete account (),  set status = -99 (deleted)  -1(invalid)
 */

func RemoveAccountById(id int64) error {
	sql := "update account set status = -99 where id = ? "
	_, err := db.Orm.Exec(sql, id)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

/**
 *  @MethodName DeleteAccountById
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 14:52
 *  @Description: Physical Delete  account  Physically
 */

func DeleteAccountById(id int64) error {
	account := new(Account)
	has, err := db.Orm.ID(id).Unscoped().Get(account)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if has {
		_, err := db.Orm.ID(id).Unscoped().Delete(account)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

/**
 *  @MethodName GetAccountById
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 14:52
 *  @Description: Get  account  ById
 */

func GetAccountById(id int64) (*Account, error) {
	res := new(Account)
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
 *  @MethodName ListAccount
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 14:52
 *  @Description: List  account
 */
func FindAccountByInstance(account *Account, page int, pagesize int) (int64, []*Account, error) {
	res := make([]*Account, 0)
	total, err := db.Orm.Count(account)
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

func GetAccountByUsername(uname string) (*Account, error) {
	res := new(Account)
	has, err := db.Orm.Where("username like ?", uname).Get(res)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if !has {
		return nil, errors.New(relations.CUS_ERR_1010)
	}
	return res, nil
}
