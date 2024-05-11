package service

import (
	log "backgate/logger"
	pb "backgate/training"
	"errors"
	"fmt"
	"strings"
)

// 不需要事务
func AddMgrResFromSwagger(cates map[string]string, items []*Apitem) (*pb.ResourcesRes, error) {
	res := &pb.ResourcesRes{}
	if len(items) == 0 || len(cates) == 0 {
		return res, errors.New("res items or cates is empty")
	}

	req := new(pb.InitResourceReq)
	req.ResCate = make([]*pb.Resource, 0)
	req.ResItem = make([]*pb.Resource, 0)

	for tag, _ := range cates {
		fmt.Println(tag)
		fmt.Println(cates[tag])
		if tag == "" || strings.Index(tag, "Mgr") <= 1 {
			return res, errors.New("cates tag is empty")
		}
		if cates[tag] == "" || strings.Index(cates[tag], "|") <= 1 {
			return res, errors.New("cates value is empty")
		}
		req.ResCate = append(req.ResCate, &pb.Resource{
			Name:      strings.ToLower(tag[:len(tag)-3] + ":all"),
			NameCn:    strings.Split(cates[tag], "|")[0],
			Url:       tag + "/*",
			Act:       "*",
			Pid:       0,
			IsLeaf:    false,
			Domain:    "backend",
			Remark:    "",
			GroupName: tag,
			Level:     0,
		})
	}
	fmt.Println("**********************************************************")
	fmt.Println(req.ResCate)
	fmt.Println("**********************************************************")

	for _, item := range items {
		fmt.Println(item)
		if item.GrpName == "" || strings.Index(item.GrpName, "Mgr") <= 1 {
			return res, errors.New("item grpname is empty")
		}
		if item.Descr == "" || strings.Index(item.Descr, "|") <= 1 {
			return res, errors.New("item descr is empty")
		}
		if item.Path == "" || item.HttpMethod == "" {
			return res, errors.New("item path or method is empty")
		}
		req.ResItem = append(req.ResItem, &pb.Resource{
			Name:      item.Name,
			NameCn:    item.NameCn,
			Url:       item.Path,
			Act:       item.HttpMethod,
			Pid:       0,
			IsLeaf:    true,
			Domain:    "backend",
			Remark:    "",
			GroupName: item.GrpName,
			Level:     1,
		})
	}
	fmt.Println("**********************************************************")
	fmt.Println(req.ResItem)
	fmt.Println("**********************************************************")

	res2, err := DealGrpcCall(req, "InitRes", "backendrpc")
	if err != nil {
		log.Error(err.Error())
		return res, err
	}

	return res2.(*pb.ResourcesRes), nil
}
