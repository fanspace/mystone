package service

import (
	"backend/apps/res/model"
	pb "backend/training"
	"github.com/jinzhu/copier"
)

// 获得所有的菜单
func QueryAllMenus(req *pb.MenuQueryReq) (*pb.MenuListRes, error) {
	res := new(pb.MenuListRes)
	return res, nil
}

// 获得单个菜单
func FetchMenu(req *pb.MenuQueryReq) (*pb.MenuRes, error) {
	res := new(pb.MenuRes)
	menu, err := model.GetMenuById(req.Mid)
	if err != nil {
		return nil, err
	}
	meta := &pb.MenuMata{
		Title:   menu.MetaTitle,
		Icon:    menu.MetaIcon,
		NoCache: menu.MetaNocache,
	}
	copier.Copy(res, menu)
	res.Meta = meta
	res.NameCn = meta.Title

	return res, nil
}
