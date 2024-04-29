package model

import (
	log "enlist_org/logger"
	"errors"
	"time"
)

type OrgDept struct {

	Id           int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	UnitId      int64  `json:"unit_id" xorm:"unique(unit_parent_name)  BIGINT(20) Not null"`
	ParentId    int64   `json:"parent_id" xorm:"unique(unit_parent_name) BIGINT(20) not null"`
	DeptName  string  `json:"dept_name" xorm:"unique(unit_parent_name) not null comment('名称中文') VARCHAR(64)"`
	Level    int32    `json:"level" xorm:"INT(5)"`
	Sort    int32    `json:"sort" xorm:"INT(8)"`
	Lft    int32    `json:"lft" xorm:"INT(5)"`
	Rgt    int32    `json:"rgt" xorm:"INT(5)"`
	Status       int32		  `json:"status" xorm:"not null default 1 INT(11)"`
	NamePy    string   `json:"name_py" xorm:"not null comment('单位名称拼音') VARCHAR(255)"`
	DeptPath    string   `json:"dept_path" xorm:"not null comment('部门path') VARCHAR(255)"`
	CreatedBy    int64  		`json:"created_by" xorm:"BIGINT(20)"`
	CreatedAt int64 		`xorm:"created"`
	UpdatedBy	 int64       `json:"updatedBy"  xorm:"BIGINT(20)"`
	UpdatedAt    int64		`json:"updatedAt"  xorm:"BIGINT(20)"`
	Version      int		`xorm:"version"`
}

/**
 *  @MethodName AddDept
 *  @author: lf6128@163.com
 *  @Date: 2022/8/9 13:25
 *  @Description: Add  dept
 */ 

func (dept *OrgDept) AddOrgDept()(*OrgDept, error) {
    _, err := Orm.Insert(dept)
    if err != nil {
       log.Error(err.Error())
       return nil, err
    }
    return dept, nil
}



/**
 *  @MethodName UpdateDept
 *  @author: lf6128@163.com
 *  @Date: 2022/8/9 13:25
 *  @Description: Update  dept
 */ 

func (dept *OrgDept) UpdateOrgDept() error {
    _, err := Orm.ID(dept.Id).AllCols().Update(dept)
    if err != nil {
       log.Error(err.Error())
       return err
    }
    return nil
}



/**
 *  @MethodName StatDeptById
 *  @author: lf6128@163.com
 *  @Date: 2022/8/9 13:25
 *  @Description: Update  dept  status
 */

func StatOrgDeptById(id int64, stat int32, operator int64, version int) error {

	sql := "update org_dept set status = ?, updated_at = ?, updated_by = ?, version = ? where id = ?  and  version = ? "
	_, err := Orm.Exec(sql, stat, time.Now().Unix(), operator, version+1, id, version)
	if err != nil {
		log.Error(err.Error())
	}
	Orm.ClearCache(new(OrgDept))
	return err
}


/**
 *  @MethodName DeleteDeptById
 *  @author: lf6128@163.com
 *  @Date: 2022/8/9 13:25
 *  @Description: Delete  dept  Physically
 */ 

func DeleteOrgDeptById(id int64) error {
	dept := new(OrgDept)
	has, err := Orm.ID(id).Unscoped().Get(dept)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if has {
		_, err := Orm.ID(id).Unscoped().Delete(dept)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}


/**
 *  @MethodName GetDeptById
 *  @author: lf6128@163.com
 *  @Date: 2022/8/9 13:25
 *  @Description: Get  dept  ById
 */ 

func GetOrgDeptById(id int64) (*OrgDept, error) {
	res := new(OrgDept)
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



/**
 *  @MethodName ListDept
 *  @author: lf6128@163.com
 *  @Date: 2022/8/9 13:25
 *  @Description: List  dept  
 */
func FindOrgDeptByInstance(dept *OrgDept, page int, pagesize int) (int64, []*OrgDept, error) {
    res :=make([]*OrgDept, 0)
    total, err := Orm.Count(dept)
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