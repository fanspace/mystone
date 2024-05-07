package utils

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const (
	XForwardedFor = "X-Forwarded-For"
	XRealIP       = "X-Real-IP"
)

// RemoteIp 返回远程客户端的 IP，如 192.168.1.1
func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get(XRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get(XForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

func randomString(l int) string {
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		if string(randInteger(65, 90)) != temp {
			temp = string(randInteger(65, 90))
			result.WriteString(temp)
			i++
		}
	}
	return result.String()
}
func RandomStr(l int) string {
	return randomString(l)
}

func RandomInt(l int) string {
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		temp = fmt.Sprintf("%d", randInteger(i+1, 9-i))
		result.WriteString(temp)
		i++
	}
	return result.String()
}
func randInteger(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

// 判断obj是否在target中，target支持的类型arrary,slice,map
func Contains(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	default:
		return false, errors.New("不支持的类型")

	}

	return false, nil
}

func RemoveRepeatElement(list []string) []string {
	// 创建一个临时map用来存储数组元素
	temp := make(map[string]bool)
	index := 0
	for _, v := range list {
		// 遍历数组元素，判断此元素是否已经存在map中
		_, ok := temp[v]
		if ok {
			list = append(list[:index], list[index+1:]...)
		} else {
			temp[v] = true
		}
		index++
	}
	return list
}

// resource列表
func IsNumIds(text string) bool {
	//r,_ := regexp.Compile(`(^\d+\|?\d+$){1,9}`)
	st := strings.Index(text, ",")
	if st == 0 || st == (len(text)-1) {
		return false
	}
	r, _ := regexp.Compile(`^\d+$`)
	text2 := strings.Replace(text, ",", "", -1)
	isCids := r.MatchString(text2)
	return isCids
}
