package service

import (
	log "backgate/logger"
	"backgate/relations"
	pb "backgate/training"
	"backgate/viewmodel"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

func InitApis() (*pb.ResourcesRes, error) {
	pblist := new(pb.ResourcesRes)
	yamlFile, err := ioutil.ReadFile("docs/swagger.yaml")
	if err != nil {
		log.Error("Error reading YAML file: " + err.Error())
		return pblist, err
	}

	var swagger SwaggerSpec
	err = yaml.Unmarshal(yamlFile, &swagger)
	if err != nil {
		log.Error("Error unmarshaling YAML: " + err.Error())
		return nil, err
	}
	apiList := make([]*Apitem, 0)
	interfacesByTag := make(map[string]string)
	for path, item := range swagger.Paths {
		hm := ""
		if item.Get != nil {
			hm = "GET"
		} else if item.Post != nil {
			hm = "POST"
		}
		for _, op := range []*Operation{item.Get, item.Post} { // 根据需要添加其他HTTP方法
			if op != nil && len(op.Tags) > 0 {
				tag := op.Tags[0] // 使用第一个tag进行分类

				if strings.LastIndex(tag, "Mgr") == len(tag)-3 {
					interfacesByTag[tag] = op.Description

					it := &Apitem{
						GrpName:    tag,
						Path:       path,
						HttpMethod: hm,
						Name:       op.Summary,
						Descr:      op.Description,
					}
					fmt.Println(it)
					apiList = append(apiList, it)
				}
			}
		}
	}
	pbapilist, err := AddMgrResFromSwagger(interfacesByTag, apiList)
	if err != nil {
		log.Error("Error adding mgr res from swagger: " + err.Error())
		return nil, err
	}
	go func() {
		jsonstr, err := proto.Marshal(pbapilist)
		if err != nil {
			log.Error(err.Error())
		}
		err = SetSimpleKey(fmt.Sprintf("%s_%s", relations.APP_NAME, "apiList"), 0, jsonstr)
		if err != nil {
			log.Error("Error setting simple key when init apis : " + err.Error())
		}
	}()

	return pbapilist, nil
}

// 获取缓存中的api列表
func ListAllApis() (*pb.ResourcesRes, error) {
	tmpres, err := GetSimpleKeyWithoutTTL(fmt.Sprintf("%s_%s", relations.APP_NAME, "apiList"))
	if err != nil {
		log.Error("Error getting simple key when list apis : " + err.Error())
		return InitApis()

	}
	res2 := new(pb.ResourcesRes)
	err = proto.Unmarshal(tmpres.([]byte), res2)
	return res2, nil
}

func RetrieveApiTree(pid int64) ([]*viewmodel.ApiTree, error) {
	res := make([]*viewmodel.ApiTree, 0)
	apilist, err := ListAllApis()
	if err != nil {
		log.Error(err.Error())
		return res, err
	}

	if pid == 0 {
		for _, v := range apilist.Resources {
			at := new(viewmodel.ApiTree)
			err = copier.Copy(at, v)
			res = append(res, at)
		}
	} else {
		for _, v := range apilist.Resources {
			fmt.Println(v.Children)
			if v.Id == pid && len(v.Children) > 0 {
				for _, k := range v.Children {
					at := new(viewmodel.ApiTree)
					err = copier.Copy(at, k)
					fmt.Println(at)
					res = append(res, at)
				}
			}
		}
	}

	return res, nil
}

// SwaggerSpec 代表Swagger规范的简化结构，仅用于演示目的
type SwaggerSpec struct {
	Paths map[string]*PathItem `yaml:"paths"`
	Tags  []Tag                `yaml:"tags"`
}

// PathItem 包含一个路径上的所有操作
type PathItem struct {
	Get  *Operation `yaml:"get,omitempty"`
	Post *Operation `yaml:"post,omitempty"`
	// 其他HTTP方法...
}

// Operation 表示一个操作详情
type Operation struct {
	Tags []string `yaml:"tags"`
	// 其他操作属性...
	Description string `yaml:"description"`
	Summary     string `yaml:"summary"`
}

// Tag 定义了一个标签
type Tag struct {
	Name string `yaml:"name"`
}

type Apitem struct {
	GrpName    string `json:"grpName"`
	Path       string `json:"path"`
	HttpMethod string `json:"httpMethod"`
	Name       string `json:"name"`
	NameCn     string `json:"nameCn"`
	Descr      string `json:"descr"`
}
