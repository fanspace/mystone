package service

import (
	. "backend/apps/res/model"
	log "backend/logger"
	pb "backend/training"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
)

// 初始化接口
func InitRes(req *pb.InitResourceReq) (*pb.ResourcesRes, error) {
	res := new(pb.ResourcesRes)
	if len(req.ResCate) == 0 || len(req.ResItem) == 0 {
		return res, errors.New("length of list is zero")
	}

	_ = addResMulti(req.ResCate)
	_ = addResMulti(req.ResItem)

	return readResMulti(req)

}

func addResMulti(list []*pb.Resource) error {

	for _, v := range list {
		it := &Resource{
			Name:      v.Name,
			GroupName: v.GroupName,
			Domain:    v.Domain,
		}
		pid, err := it.HasResExist()
		if err != nil {
			log.Error(err.Error())
			continue
		}
		if pid == 0 {
			err = copier.Copy(it, v)
			if err != nil {
				log.Error(err.Error())
				continue
			} else {
				_, err = it.InsertResorce()
				if err != nil {
					log.Error(err.Error())
					continue
				}

			}
		}
	}

	return nil
}

func readResMulti(req *pb.InitResourceReq) (*pb.ResourcesRes, error) {
	res := new(pb.ResourcesRes)
	if len(req.ResCate) == 0 || len(req.ResItem) == 0 {
		return res, errors.New("length of list is zero")
	}

	res.Resources = make([]*pb.Resource, 0)

	for _, v := range req.ResCate {
		cate, err := GetParentZjResourcesByGrpName(v.GroupName, v.Domain)
		if err != nil {
			log.Error(err.Error())
			continue

		}
		pbcate := new(pb.Resource)
		err = copier.Copy(pbcate, cate)
		fmt.Println("ssssssssssssssssssssssssssssssssssssssssss")
		fmt.Println(pbcate)
		pbcate.Children = make([]*pb.Resource, 0)
		reslist, err := FindResourcesByPid(pbcate.Id)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		fmt.Println("########################################")
		fmt.Println(reslist)
		if len(reslist) > 0 {
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			for _, v := range reslist {
				pbitem := new(pb.Resource)
				err = copier.Copy(pbitem, v)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				pbcate.Children = append(pbcate.Children, pbitem)
			}
		}
		res.Resources = append(res.Resources, pbcate)

	}

	fmt.Println(res)
	return res, nil
}
