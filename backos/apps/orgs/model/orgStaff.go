package model
import (
	"enlist_org/internal/common"
	log "enlist_org/logger"
	"errors"
)
type  OrgStaff struct  {
	Id  int64   `json:"id" xorm:"pk autoincr BIGINT(20)"`
	Usid  int64   `json:"usid" xorm:"NOT NULL BIGINT(20)"`
	UnitId  int64  `json:"unitId" xorm:"NOT NULL BIGINT(20)"`
	DeptId  int64  `json:"deptId" xorm:"NOT NULL BIGINT(20)"`
	Status  int32  `json:"status" xorm:"INT(5)"`
	UnitName  string  `json:"unitName xorm:"VARCHAR(64)"`
	DeptName string   `json:"deptName" xorm:"VARCHAR(64)"`
	IsManager  bool  `json:"isManager" xorm:"Bool"`
	IsRoot    bool  `json:"isRoot" xorm:"Bool"`
	Relations  string  `json:"relations" xorm:"VARCHAR(512)"`
	UpdatedBy       int64  `json:"updatedBy"  xorm:"BIGINT(20)"`
	UpdatedAt       int64  `json:"updatedAt"  xorm:"BIGINT(20)"`
	CreatedBy    int64  `json:"created_by" xorm:"BIGINT(20)"`
	CreatedAt    int64  `xorm:"created"`
	Version      int    `xorm:"version"`

}

/**
 *  @MethodName AddOrgStaff
 *  @author: lf6128@163.com
 *  @Date: 2022/7/21 16:48
 *  @Description: Add  orgStaff
 */

func (orgStaff *OrgStaff) AddOrgStaff()(*OrgStaff, error) {
    _, err := Orm.Insert(orgStaff)
    if err != nil {
       log.Error(err.Error())
       return nil, err
    }
    return orgStaff, nil
}



/**
 *  @MethodName UpdateOrgStaff
 *  @author: lf6128@163.com
 *  @Date: 2022/7/21 16:48
 *  @Description: Update  orgStaff
 */

func (orgStaff *OrgStaff) UpdateOrgStaff() error {
    _, err := Orm.ID(orgStaff.Id).AllCols().Update(orgStaff)
    if err != nil {
       log.Error(err.Error())
       return err
    }
    return nil
}



/**
 *  @MethodName StatOrgStaffById
 *  @author: lf6128@163.com
 *  @Date: 2022/7/21 16:48
 *  @Description: Update  orgStaff  status
 */

func StatOrgStaffById(id int64, stat int32) error {
	sql := "update org_staff set status = ? where id = ? "
	_, err := Orm.Exec(sql, stat, id)
	if err != nil {
	 	log.Error(err.Error())
	 }
	 return err
}



/**
 *  @MethodName DeleteOrgStaffById
 *  @author: lf6128@163.com
 *  @Date: 2022/7/21 16:48
 *  @Description: Delete  orgStaff  Physically
 */

func DeleteOrgStaffById(id int64) error {
	orgStaff := new(OrgStaff)
	has, err := Orm.ID(id).Unscoped().Get(orgStaff)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if has {
		_, err := Orm.ID(id).Unscoped().Delete(orgStaff)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}


/**
 *  @MethodName GetOrgStaffById
 *  @author: lf6128@163.com
 *  @Date: 2022/7/21 16:48
 *  @Description: Get  orgStaff  ById
 */

func GetOrgStaffById(id int64) (*OrgStaff, error) {
	res := new(OrgStaff)
	has, err := Orm.ID(id).Get(res)
	if err != nil {
	 	log.Error(err.Error())
		return nil, err
	 }
	 if !has {
	 return nil, errors.New("data not exist")
	 }
	 return res, nil
}


func GetOrgStaffByUsId(usid int64) (*OrgStaff, error) {
	res := new(OrgStaff)
	has, err := Orm.Where(" usid = ?", usid).Get(res)
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
 *  @MethodName ListOrgStaff
 *  @author: lf6128@163.com
 *  @Date: 2022/7/21 16:48
 *  @Description: List  orgStaff
 */
func FindOrgStaffByInstance(orgStaff *OrgStaff, page int, pagesize int) (int64, []*OrgStaff, error) {
    res :=make([]*OrgStaff, 0)
    total, err := Orm.Count(orgStaff)
    if err != nil {
    log.Error(err.Error())
  return 0, res, err
  }
 if total > 0 {
 err := Orm.Where("id > ?", 0).Limit(pagesize,(page-1)*pagesize).Find(&res)
 if err != nil {
   log.Error(err.Error())
   return 0, res, err
  }
 }
 return total, res, nil
}



func BatchUpdateDeptNameOrRemoveDept(req *OrgStaff, act int32) error {
	sql := ""
	var err error
	if act == 2 {
		sql = "update org_staff set dept_name = ? where unit_id = ? and dept_id = ? "
		_, err = Orm.Exec(sql, req.DeptName, req.UnitId, req.DeptId)
	} else if act == 3 {
		sql = "update org_staff set dept_name = '', dept_id = 0, is_manager = 0, relations='' where unit_id = ? and dept_id = ? "
		_, err = Orm.Exec(sql, req.UnitId, req.DeptId)
	} else {
		return errors.New(common.CUS_ERR_4002)
	}
	if err != nil {
		log.Error(err.Error())
		return err
	}
	Orm.ClearCache(new(OrgStaff))
	return nil
}


func BatchUpdateUnitNameOrRemoveUnit(req *OrgStaff, act int32) error {
	sql := ""
	var err error
	if act == 2 {
		sql = "update org_staff set unit_name = ? where unit_id = ? "
		_, err = Orm.Exec(sql, req.UnitName, req.UnitId)
	} else if act == 3 {
		sql = "update org_staff set unit_name = '', unit_id = 0, dept_id = 0, dept_name = '', is_root =0, is_manager = 0, relations='' where unit_id = ? "
		_, err = Orm.Exec(sql, req.UnitId)
	} else {
		return errors.New(common.CUS_ERR_4002)
	}


	if err != nil {
		log.Error(err.Error())
		return err
	}
	Orm.ClearCache(new(OrgStaff))
	return nil

}
