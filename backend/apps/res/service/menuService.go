package service

import (
	. "backend/apps/res/model"
	"backend/db"
	log "backend/logger"
	"backend/relations"
	pb "backend/training"
	"backend/utils"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/slice"
	"github.com/jinzhu/copier"
	"regexp"
	"strings"
	"time"
)

// 获得单个菜单
func FetchMenu(req *pb.MenuQueryReq) (*pb.MenuRes, error) {
	res := new(pb.MenuRes)
	menu, err := GetMenuById(req.Mid)
	if err != nil {
		return nil, err
	}
	meta := &pb.MenuMata{
		Title:   menu.MetaTitle,
		Icon:    menu.MetaIcon,
		NoCache: menu.MetaNocache,
		Type:    menu.MetaType,
	}
	copier.Copy(res, menu)
	res.Meta = meta
	res.NameCn = meta.Title

	return res, nil
}

// 获得某用户授权的菜单
func QueryAuthMenus(req *pb.UserAuthMenuListReq) (*pb.MenuListRes, error) {
	res := new(pb.MenuListRes)
	isRoot := slice.Contain(req.Roles, "root")
	if isRoot {
		req2 := new(pb.MenuQueryReq)
		req2.Domain = req.Domain
		req2.Type = 0
		return QueryMenuList(req2)
	} else {
		//todo
		//  get  menu groups  from   rbacRpc
		//
		MenuHandleReq := new(pb.AuthMenuReq)
		MenuHandleReq.Domain = req.Domain
		MenuHandleReq.Roles = req.Roles
		MenuHandleReq.Usid = req.Usid
		menulist, err := QueryMenusByRole(MenuHandleReq)
		if err != nil {
			log.Error(err.Error())
			return res, err
		}

		for _, v := range menulist.AuthMenu {
			menu, err := findMenusByGrpNameAndConds(v.Group, v.MenuIds, req.Domain)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			res.Menus = append(res.Menus, menu)

		}

	}
	res.Total = int64(len(res.Menus))
	return res, nil
}

func QueryMenuList(req *pb.MenuQueryReq) (*pb.MenuListRes, error) {
	res := new(pb.MenuListRes)
	//类型 所有：0   某个pid:1   某个group : 2 (需要grpname)    自身？ 取消mid :3
	if req.Type == 0 {
		menus, err := QueryAllProtosMenu(req.Domain)
		if err != nil {
			return nil, err
		}
		res.Menus = menus
	} else if req.Type == 1 {
		if req.Pid == 0 {
			return nil, errors.New(relations.CUS_ERR_4002)
		}
		menus, err := findMenuProtosByPidAndDomain(req.Pid, req.Domain)
		if err != nil {
			return nil, err
		}
		res.Menus = menus

	} else if req.Type == 2 {
		if req.GroupName == "" {
			return nil, errors.New(relations.CUS_ERR_4002)
		}
		tmpres, err := getParentMenuByGrpName(req.GroupName, req.Domain)
		if err != nil {
			return nil, err
		}
		if tmpres == nil {
			return nil, errors.New(relations.CUS_ERR_4004)
		}
		resc, err := recurMenuVmById(tmpres.Id, req.Domain)
		if err != nil {
			return nil, err
		}
		res.Menus = append(res.Menus, resc)

	}
	res.Total = int64(len(res.Menus))
	return res, nil
}

// 获得所有菜单
func QueryAllProtosMenu(domain string) ([]*pb.MenuRes, error) {
	res := make([]*pb.MenuRes, 0)
	// 取得一级菜单，没实际用，只为了循环
	levelOne, err := findMenuProtosByPidAndDomain(relations.ZERO, domain)
	if err != nil {
		return res, err
	}
	if len(levelOne) > 0 {
		for _, v := range levelOne {
			menuitem, err := recurMenuVmById(v.Id, domain)
			if err != nil {
				return res, err
			}
			res = append(res, menuitem)
		}
	}

	return res, nil
}

// tree 方法
func getMenuProtossBytId(id int64) (*pb.MenuRes, error) {
	en, err := GetMenuById(id)
	if err != nil {
		return nil, err
	}
	res := new(pb.MenuRes)
	meta := &pb.MenuMata{
		Title:   en.MetaTitle,
		Icon:    en.MetaIcon,
		NoCache: en.MetaNocache,
		Type:    en.MetaType,
	}
	copier.Copy(res, en)
	res.Meta = meta
	res.NameCn = meta.Title
	return res, nil
}

// 单层,子级
func findMenuProtosByPidAndDomain(pid int64, domain string) ([]*pb.MenuRes, error) {
	mlist := make([]*pb.MenuRes, 0)
	elist, err := FindMenuByPidAndDomain(pid, domain)
	if err != nil {
		return mlist, err
	}
	for _, v := range elist {
		zj := new(pb.MenuRes)
		copier.Copy(zj, v)
		meta := &pb.MenuMata{
			Title:   v.MetaTitle,
			Icon:    v.MetaIcon,
			NoCache: v.MetaNocache,
			Type:    v.MetaType,
		}
		zj.Meta = meta
		zj.NameCn = meta.Title
		mlist = append(mlist, zj)
	}
	return mlist, nil
}

// 根据groupname  获得该组的1级菜单
func getParentMenuByGrpName(grpname string, domain string) (*pb.MenuRes, error) {
	pmenu, err := GetParentMenuByGrpName(grpname, domain)
	if err != nil {
		return nil, err
	}
	res := new(pb.MenuRes)
	copier.Copy(res, pmenu)
	meta := &pb.MenuMata{
		Title:   pmenu.MetaTitle,
		Icon:    pmenu.MetaIcon,
		NoCache: pmenu.MetaNocache,
		Type:    pmenu.MetaType,
	}
	res.Meta = meta
	res.NameCn = meta.Title
	return res, nil
}

// 递归 单个
func recurMenuVmById(id int64, domain string) (*pb.MenuRes, error) {
	pmenu, err := getMenuProtossBytId(id)
	if err != nil {
		return nil, err
	}
	if pmenu != nil {
		if !pmenu.IsLeaf {
			pmenu.Children, err = findMenuProtosByPidAndDomain(pmenu.Id, domain)
		}
		if err != nil {
			return nil, err
		}
		if len(pmenu.Children) > 0 {
			for k, _ := range pmenu.Children {
				pmenu.Children[k], err = recurMenuVmById(pmenu.Children[k].Id, domain)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return pmenu, nil
}

func MutationMenu(req *pb.MenuHandleReq) (*pb.CommonResponse, error) {
	res := new(pb.CommonResponse)
	if req.Action >= 1 && req.Action <= 3 && req.Menu == nil {
		fmt.Println("222222222222")
		return res, errors.New(relations.CUS_ERR_4002)
	}
	if req.Action == 9 && len(req.MenuIds) == 0 {
		fmt.Println("33333333333333333333")
		return res, errors.New(relations.CUS_ERR_4002)
	}
	menu := new(MenuEntity)
	if req.Action != pb.Action_REMOVE {
		copier.Copy(menu, req.Menu)
		menu.Redirect = req.Menu.Redirect
		menu.MetaNocache = req.Menu.Meta.NoCache
		menu.MetaIcon = req.Menu.Meta.Icon
		menu.MetaTitle = req.Menu.Meta.Title
		menu.MetaType = req.Menu.Meta.Type
	}

	switch req.Action {
	case pb.Action_NEW:
		if menu.Name == "" || menu.Domain == "" {
			//res.Success = false
			//res.Msg = common.CUS_ERR_4002
			return nil, errors.New(relations.CUS_ERR_4002)
		}
		menu.CreatedAt = time.Now().Unix()
		menu.CreatedBy = req.Operator
		id, err := addMenuSrv(menu)
		if err != nil {
			//res.Success = false
			//res.Msg = err.Error()
			return nil, errors.New(err.Error())
		}
		res.Id = id
	case pb.Action_UPDATE:
		if menu.Path == "" || menu.Name == "" || menu.Domain == "" || menu.Id == 0 {
			//res.Success = false
			//res.Msg = common.CUS_ERR_4002
			return nil, errors.New(relations.CUS_ERR_4002)
		}
		menu.UpdatedAt = time.Now().Unix()
		menu.UpdatedBy = req.Operator
		err := updateMenuSrv(menu)
		if err != nil {
			//res.Success = false
			//res.Msg = err.Error()
			return nil, errors.New(err.Error())
		}
		res.Id = menu.Id
	case pb.Action_DELETE:
		if menu.Id == 0 {
			//res.Success = false
			//res.Msg = common.CUS_ERR_4002
			return nil, errors.New(relations.CUS_ERR_4002)
		}
		if !menu.IsLeaf {
			//res.Success = false
			//res.Msg = common.CUS_ERR_4005
			return nil, errors.New(relations.CUS_ERR_4005)
		}
		err := DeleteMenuById(menu.Id)
		if err != nil {
			//res.Success = false
			//res.Msg = err.Error()
			return nil, errors.New(err.Error())
		}
		go DeteleBindByMenuId(menu.Id)
		res.Id = menu.Id
	case pb.Action_REMOVE:
		if len(req.MenuIds) == 0 {
			return nil, errors.New(relations.CUS_ERR_4002)
		}
		for _, v := range req.MenuIds {
			err := DeleteMenuById(v)
			if err != nil {
				return nil, errors.New(err.Error())
			}
			go DeteleBindByMenuId(v)
		}

	default:
		//res.Success = false
		//res.Msg = common.CUS_ERR_4002
		return nil, errors.New(relations.CUS_ERR_4002)
	}
	res.Success = true
	return res, nil
}

func addMenuSrv(menu *MenuEntity) (int64, error) {
	if menu.Pid == 0 {
		return menu.InsertMenu()
	} else {
		pmenu, err := GetMenuById(menu.Pid)
		if err != nil {
			return relations.ZERO, err
		}
		if pmenu.IsLeaf {
			pmenu.IsLeaf = false
			_, err = pmenu.UpdateMenu()
			if err != nil {
				return relations.ZERO, err
			}
		}
		return menu.InsertMenu()

	}
}
func updateMenuSrv(menu *MenuEntity) error {
	children := make([]*MenuEntity, 0)
	oldmenu, err := GetMenuById(menu.Id)
	if err != nil {
		return err
	}
	if menu.Version != oldmenu.Version {
		menu.Version = oldmenu.Version
		log.Info("更新菜单时有乐观锁异常，暂时强制更新，需要检查")
	}

	if oldmenu == nil {
		return errors.New(relations.CUS_ERR_4004)
	}

	if !oldmenu.IsLeaf {
		children, err = FindMenuByPid(oldmenu.Id)
		if err != nil {
			return err
		}
	}
	if len(children) > 0 && menu.IsLeaf {
		return errors.New(relations.CUS_ERR_4016)
	}
	if menu.GroupName != oldmenu.GroupName {
		err = batchUpdateMenuGroup(menu)
		if err != nil {
			return err
		}
	} else if menu.Level != oldmenu.Level && len(children) > 0 {
		err = batchUpdateMenuLevel(menu)
		if err != nil {
			return err
		}
	} else {
		_, err = menu.UpdateMenu()
		if err != nil {
			return err
		}
	}

	return nil
}

func batchUpdateMenuGroup(menu *MenuEntity) error {
	if menu.GroupName == "" {
		return errors.New(relations.CUS_ERR_1008)
	}
	session := db.Orm.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	oldmenu := new(MenuEntity)
	_, err = session.ID(menu.Id).Get(oldmenu)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if oldmenu == nil {
		return errors.New(relations.CUS_ERR_4004)
	}

	_, err = session.ID(menu.Id).Update(menu)
	if err != nil {
		session.Rollback()
		log.Error(err.Error())
		return err
	}
	sql := "update menu_entity set `group_name` = ? where pid =  ? and pid > 0 "
	_, err = session.Exec(sql, menu.GroupName, menu.Id)
	if err != nil {
		session.Rollback()
		log.Error(err.Error())
		return err
	}
	err = session.Commit()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	db.Orm.ClearCache(new(MenuEntity))
	return nil
}

func batchUpdateMenuLevel(menu *MenuEntity) error {
	session := db.Orm.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	oldmenu := new(MenuEntity)
	_, err = session.ID(menu.Id).Get(oldmenu)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if oldmenu == nil {
		return errors.New(relations.CUS_ERR_4004)
	}

	offset := menu.Level - oldmenu.Level
	sql := ""
	if offset < 0 {
		sql = fmt.Sprintf("update zj_menu_entity set level = level %d  where pid = ? ", offset)
	} else {
		sql = fmt.Sprintf("update zj_menu_entity set level = level + %d  where pid = ? ", offset)
	}

	_, err = session.Exec(sql, menu.Id)
	if err != nil {
		session.Rollback()
		log.Error(err.Error())
		return err
	}

	_, err = session.ID(menu.Id).Update(menu)
	if err != nil {
		session.Rollback()
		log.Error(err.Error())
		return err
	}

	err = session.Commit()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	db.Orm.ClearCache(new(MenuEntity))
	return nil
}

// 绑定资源列表到菜单
func BindResWithMenuId(req *pb.MenuResBindReq) (*pb.CommonResponse, error) {
	res := new(pb.CommonResponse)
	res.Success = false
	if req.MenuId == 0 || len(req.ResIds) == 0 {
		return res, errors.New(relations.CUS_ERR_4002)
	}

	session := db.Orm.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		log.Error(err.Error())
		return res, err
	}
	_, err = session.Exec("delete from menu_res where menu_id = ? ", req.MenuId)
	if err != nil {
		log.Error(err.Error())
		session.Rollback()
		return res, err
	}
	db.Orm.ClearCache(new(MenuRes))
	mrlist := make([]*MenuRes, 0)
	for _, v := range req.ResIds {
		tt := time.Now().Unix()
		mr := &MenuRes{
			Id:        0,
			MenuId:    req.MenuId,
			ResId:     v,
			CreatedAt: tt,
			UpdatedAt: tt,
			Version:   1,
		}
		mrlist = append(mrlist, mr)
	}
	_, err = session.Insert(mrlist)
	if err != nil {
		log.Error(err.Error())
		session.Rollback()
		return res, err
	}

	err = session.Commit()
	if err != nil {
		log.Error(err.Error())
		return res, err
	}
	//todo
	//re build rbac in sysapi
	//todo
	//go reDisPerms(menuid, opid, ip)
	res.Success = true
	res.Id = 1
	return res, nil
}

// 根据某个组名，获得该组完整菜单
func findMenusByGrpNameAndConds(grpname string, ids string, domain string) (*pb.MenuRes, error) {
	if grpname == "" {
		return nil, errors.New(relations.CUS_ERR_4002)
	}
	rootmenu, err := getParentMenuByGrpName(grpname, domain)
	if err != nil {
		return nil, err
	}
	// 获取全部
	if ids == "" {
		menuitem, err := recurMenuVmById(rootmenu.Id, domain)
		if err != nil {
			return nil, err
		}
		return menuitem, nil
	} else {
		isok := validateMenuIdArr(ids)
		if !isok {
			return nil, errors.New(relations.CUS_ERR_4002)
		}
		menuChildren := make([]*pb.MenuRes, 0)
		menuarr, err := FindMenusByIdArr(ids)
		if err != nil {
			return nil, err
		}
		for _, v := range menuarr {
			if v.GroupName == grpname && v.Pid == rootmenu.Id {
				menu := new(pb.MenuRes)
				meta := &pb.MenuMata{
					Title:   v.MetaTitle,
					Icon:    v.MetaIcon,
					NoCache: v.MetaNocache,
					Type:    v.MetaType,
				}
				copier.Copy(menu, v)
				menu.Meta = meta
				menu.NameCn = meta.Title
				menuChildren = append(menuChildren, menu)
			}
		}
		for k, v := range menuChildren {
			if !v.IsLeaf {
				leve2menu, err := findMenuProtosByPidAndDomain(v.Id, domain)
				if err != nil {
					log.Error(err.Error())
				} else {
					menuChildren[k].Children = leve2menu
				}
			}
		}
		rootmenu.Children = menuChildren
	}
	return rootmenu, nil
}

func validateMenuIdArr(arr string) bool {
	r, _ := regexp.Compile(`^\d+(\,\d+)*$`)
	return r.MatchString(arr)

}

// 根据role 查询授权用户可显示的菜单
func QueryMenusByRole(req *pb.AuthMenuReq) (*pb.AuthMenuRes, error) {
	res := new(pb.AuthMenuRes)
	if !slice.Contain(relations.DOMAINS_LIMITED, req.Domain) || len(req.Roles) == 0 || req.Usid == 0 {
		return res, errors.New(relations.CUS_ERR_4008)
	}
	orignmenus := make([]*pb.AuthMenu, 0)
	for _, v := range req.Roles {
		perms, err := getMenusByRole(v, req.Domain)
		if err != nil {
			log.Error(err.Error())
			continue
		} else {
			for _, it := range perms {
				menu := new(pb.AuthMenu)
				menu.Group = it.Group
				menu.MenuIds = it.Menus
				orignmenus = append(orignmenus, menu)
			}
		}
	}
	res.AuthMenu = rebuildAuthMenus(orignmenus)
	return res, nil
}

func rebuildAuthMenus(old []*pb.AuthMenu) []*pb.AuthMenu {
	res := make([]*pb.AuthMenu, 0)
	for _, v := range old {
		k := isContainsGroup(v.Group, res)
		if k == -1 {
			res = append(res, v)
		} else {
			if utils.IsNumIds(res[k].MenuIds) {
				res[k].MenuIds += "," + v.MenuIds
			}
		}
	}
	for k, v := range res {
		newarr := utils.RemoveRepeatElement(strings.Split(v.MenuIds, ","))
		res[k].MenuIds = strings.Join(newarr, ",")

	}
	return res
}

func isContainsGroup(group string, arr []*pb.AuthMenu) int {
	if len(arr) == 0 {
		return -1
	}
	for k, v := range arr {
		if v.Group == group {
			return k
		}
	}
	return -1
}

func getMenusByRole(rolename string, domain string) ([]*Perms, error) {
	res := make([]*Perms, 0)
	rolereq := new(pb.RolesListReq)
	rolereq.Domain = domain
	rolereq.RoleName = rolename
	role, err := GetSysRolesByNameAndDomain(rolename, domain)
	if err != nil {
		log.Error(err.Error())
		return res, err
	}
	if role.Perms == "" {
		return res, errors.New(relations.CUS_ERR_4004)
	}
	if !utils.IsNumIds(role.Perms) {
		return res, errors.New(relations.CUS_ERR_4004)
	}
	return FindPermsByIdArr(role.Perms)
}
