package service

import (
	log "backgate/logger"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

func InitApis() (*ApiList, error) {
	applist := new(ApiList)
	yamlFile, err := ioutil.ReadFile("docs/swagger.yaml")
	if err != nil {
		log.Error("Error reading YAML file: " + err.Error())
		return applist, err
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
	_, err = AddMgrResFromSwagger(interfacesByTag, apiList)
	if err != nil {
		log.Error("Error adding mgr res from swagger: " + err.Error())
		return nil, err
	}
	if len(apiList) > 0 {
		applist.Apis = apiList
		jsoner, _ := json.Marshal(applist)
		//fmt.Println(string(jsoner))
		go func() {
			err := SetSimpleKey("apiList", 0, string(jsoner))
			if err != nil {
				log.Error("Error setting simple key when init apis : " + err.Error())
			}
		}()
	}
	return applist, nil
}

func ListAllApis() (*ApiList, error) {
	res := new(ApiList)
	tmpres, err := GetSimpleKeyWithoutTTL("apiList")
	if err != nil {
		log.Error("Error getting simple key when list apis : " + err.Error())
		return InitApis()

	}

	err = json.Unmarshal([]byte(tmpres.([]uint8)), res)
	if err != nil {
		log.Error("Error getting simple key when list apis : " + err.Error())
		return InitApis()

	}
	for _, item := range res.Apis {
		fmt.Println(item.GrpName, item.Path, item.HttpMethod)
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

type ApiList struct {
	Apis []*Apitem `json:"apis"`
}