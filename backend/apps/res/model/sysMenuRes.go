package model

import (
	"backend/db"
	log "backend/logger"
)

type MenuRes struct {
	Id        int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	MenuId    int64  `json:"menu_id"`
	ResId     int64  `json:"res_id"`
	Domain    string `json:"domain"`
	CreatedAt int64  `xorm:"created"`
	UpdatedAt int64  `xorm:"updated"`
	Version   int    `xorm:"version"`
}

func GetResByMenuId(mid int64, domain string) ([]*Resource, error) {
	res := make([]*Resource, 0)

	err := db.Orm.Table(&Resource{}).Join("INNER", "menu_res", "resource.id = menu_res.res_id").Where("menu_res.menu_id = ?", mid).And("menu_res.domain like ?", domain).Find(&res)
	if err != nil {
		log.Error(err.Error())
		return res, err
	}
	return res, nil
}

func DeteleBindByMenuId(mid int64) error {
	_, err := db.Orm.Exec("delete from menu_res where menu_id = ?", mid)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	db.Orm.ClearCache(new(MenuRes))
	return nil
}

type BindResReq struct {
	// must 1
	Action int     `json:"action"`
	MenuId int64   `json:"menu_id"`
	ResIds []int64 `json:"res_ids"`
}

func InitMenuRes() error {
	res := &MenuRes{
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
		mrlist := make([]*MenuRes, 0)
		mrarr := [...]*MenuRes{
			&MenuRes{
				Id:        1,
				MenuId:    2,
				ResId:     2,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        2,
				MenuId:    2,
				ResId:     3,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        3,
				MenuId:    2,
				ResId:     4,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        4,
				MenuId:    3,
				ResId:     102,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        5,
				MenuId:    3,
				ResId:     103,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        6,
				MenuId:    3,
				ResId:     104,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        7,
				MenuId:    4,
				ResId:     5,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        8,
				MenuId:    4,
				ResId:     6,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        9,
				MenuId:    5,
				ResId:     202,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        10,
				MenuId:    5,
				ResId:     203,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        11,
				MenuId:    5,
				ResId:     204,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        12,
				MenuId:    102,
				ResId:     303,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        13,
				MenuId:    103,
				ResId:     305,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        14,
				MenuId:    104,
				ResId:     302,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        15,
				MenuId:    105,
				ResId:     302,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        16,
				MenuId:    105,
				ResId:     304,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        17,
				MenuId:    106,
				ResId:     307,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        18,
				MenuId:    107,
				ResId:     307,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			}, &MenuRes{
				Id:        19,
				MenuId:    107,
				ResId:     306,
				Domain:    "back",
				CreatedAt: 1632734800,
				UpdatedAt: 1632734800,
				Version:   1,
			},
		}
		mrlist = mrarr[:]
		_, err := db.Orm.Insert(mrlist)
		if err != nil {
			log.Error(err.Error())
			return err
		}

	}
	return nil
}
