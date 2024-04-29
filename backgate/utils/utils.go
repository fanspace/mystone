package utils

import (
	log "backgate/logger"
	"backgate/relations"
	"backgate/settings"
	"backgate/validations"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	XForwardedFor = "X-Forwarded-For"
	XRealIP       = "X-Real-IP"
)

type MyCustomClaims struct {
	Domain   string `json:"domain"`
	Usid     int64  `json:"usid"`
	AgencyId int64  `json:"agency_id"`
	Username string `json:"username"`
	UserType int32  `json:"userType"`
	Status   int32  `json:"status"`
	Device   int32  `json:"device"`
	jwt.StandardClaims
}
type MyClaim struct {
	Domain   string `json:"domain"`
	Usid     int64  `json:"usid"`
	AgencyId int64  `json:"agency_id"`
	Username string `json:"username"`
	UserType int32  `json:"userType"`
	Status   int32  `json:"status"`
	Device   int32  `json:"device"`
}

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

// 前端解析jwt
func ParseJwt(tokenstr string, logintype string) (bool, *MyClaim) {
	token, err := jwt.ParseWithClaims(strings.TrimSpace(tokenstr), &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if settings.Cfg.ReleaseMode {
			switch logintype {
			case "fore":
				return []byte(relations.JWT_SECRET_STRING_PROD_NOR), nil
			case "com":
				return []byte(relations.JWT_SECRET_STRING_PROD_COM), nil
			case "back":
				return []byte(relations.JWT_SECRET_STRING_PROD_MAN), nil
			default:
				return nil, errors.New(relations.CUS_ERR_4002)
			}

		} else {
			return []byte(relations.JWT_SECRET_STRING_DEV), nil
		}
	})
	if err != nil {
		log.Error(err.Error())
		return false, nil
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		mc := new(MyClaim)
		mc.Domain = claims.Domain
		mc.Usid = claims.Usid
		mc.Username = claims.Username
		mc.UserType = claims.UserType
		mc.Status = claims.Status
		mc.AgencyId = claims.AgencyId
		mc.Device = claims.Device
		return true, mc
	} else {
		log.Info(fmt.Sprintf("%v %v", claims.Usid, claims.StandardClaims.ExpiresAt))
		return false, nil
	}
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

func ConvertExt2ContentType(ext string) (string, error) {

	upext := strings.ToUpper(ext)
	switch upext {
	case "JPG", "JPEG":
		return "image/jpeg", nil
	case "PNG", "GIF":
		return fmt.Sprintf("image/%s", strings.ToLower(ext)), nil
	// application
	case "ZIP", "CVS", "PDF":
		return fmt.Sprintf("application/%s", strings.ToLower(ext)), nil
		// text
	case "HTM", "HTML", "CSS", "XML":
		return fmt.Sprintf("text/%s", strings.ToLower(ext)), nil
	case "TXT":
		return "text/plain", nil
		// video
	case "AVI", "MP4":
		return fmt.Sprintf("video/%s", strings.ToLower(ext)), nil
		// audio
	case "MP3":
		return fmt.Sprintf("audio/%s", strings.ToLower(ext)), nil
		// office
	case "DOC":
		return "application/msword", nil
	case "DOCX":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document", nil
	case "XLS":
		return "application/x-xls", nil
	case "XLSX":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
	case "PPT", "PPS":
		return "application/vnd.ms-powerpoint", nil
	case "SVG":
		return "text/xml", nil
	default:
		return "", errors.New(relations.CUS_ERR_4100)
	}
}

// 根据职称值，返回当前职称等级  -1:无效值 ,0：无职称， 1：初级， 2：中级， 3：高级
// 取余
/*
   c := 1511
   for c >= 10 {
	c = GetRemainder(c)
	}

*/
func GetRemainder(val int) int {
	if val < 10 {
		return val
	}
	bits := len(strconv.Itoa(val)) - 1
	offset := int(math.Pow(10, float64(bits)))
	return val % offset
}

func ConvertJsonStrNoQuote(str string) string {
	re := regexp.MustCompile(`([a-zA-Z]+?):`)
	rep1 := `"${1}":`
	res := fmt.Sprintf("%q\n", re.ReplaceAllString(str, rep1))
	res2 := strings.Replace(res, "\\", "", -1)
	if string(res2[0]) == `"` {
		res3 := strings.Replace(res2, `"{`, "{", -1)
		res3 = strings.Replace(res3, `}"`, "}", -1)
		//fmt.Println(res3)
		return res3
	}
	return res2
}

// 将get请求的参数进行转义
func getParseParam(param string) string {
	return url.PathEscape(param)
}

func HttpHandle(method, urlVal, data, tokens string) (string, error) {
	// `这里请注意，使用 InsecureSkipVerify: true 来跳过证书验证`
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}

	var req *http.Request

	if data == "" {
		urlArr := strings.Split(urlVal, "?")
		if len(urlArr) == 2 {
			urlVal = urlArr[0] + "?" + getParseParam(urlArr[1])
		}
		req, _ = http.NewRequest(method, urlVal, nil)
	} else {
		req, _ = http.NewRequest(method, urlVal, strings.NewReader(data))
	}

	//添加cookie，key为X-Xsrftoken，value为df41ba54db5011e89861002324e63af81
	//可以添加多个cookie
	//cookie1 := &http.Cookie{Name: "X-Xsrftoken",Value: "df41ba54db5011e89861002324e63af81", HttpOnly: true}
	//req.AddCookie(cookie1)

	//添加header，key为X-Xsrftoken，value为b6d695bbdcd111e8b681002324e63af81
	if tokens != "" {
		req.Header.Add("Authorization", tokens)
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Error(err.Error())
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	defer client.CloseIdleConnections()
	return string(b), nil
}

func ConvertIdcodeToGender(idcode string) int {
	if !validations.IsIdCardStrict(idcode) {
		return 0
	}
	gender, err := strconv.Atoi(idcode[16:17])
	if err != nil {
		return 0
	}
	if gender%2 == 0 {
		return 2
	}
	return 1
}

func ConvertIdcodeToBirth(idcode string) string {
	if !validations.IsIdCardStrict(idcode) {
		return ""
	}
	return fmt.Sprintf("%s-%s-%s", idcode[6:10], idcode[10:12], idcode[12:14])
}

func DistributeAreaByDistrict(dist int32) int64 {
	/*
		1353  市南
		1354  市北
		1355  黄岛区
		1356  崂山区
		1357  李沧区
		1358  城阳区
		1359  即墨区
		1360  胶州市
		1361  平度市
		1362  莱西市
	*/
	switch dist {
	case 1353:
		return 69
	case 1354:
		return 70
	case 1355, 3222, 3224, 3225:
		return 76
	case 1356:
		return 73
	case 1357:
		return 71
	case 1358, 3220:
		return 72
	case 1359, 3223:
		return 74
	case 1360, 3221:
		return 75
	case 1361:
		return 77
	case 1363:
		return 78
	default:
		return 1
	}
}
