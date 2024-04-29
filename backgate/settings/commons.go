package settings

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var versions = []string{
	"1.0",
}

var CurrentVersion string = versions[0]
var BuildDate string
var BuildHash string

const (
	LOG_FILENAME   = "critic.log"
	CONF_DES_TOKEN = "38df159ts535e258"
)

var (
	commonBaseSearchPaths = []string{
		".",
		"..",
		"../..",
		"../../..",
	}
)

const (
	InvalidParam        = 4000
	InvalidFrontImg     = 4001
	InvalidBackImg      = 4002
	InvalidEmail        = 4003
	InvalidMobile       = 4004
	InvalidWechat       = 4005
	InvalidQQ           = 4006
	InvalidDescription  = 4007
	InvalidType         = 4008
	AlreadyApproved     = 4009
	AlreadySubmit       = 4010
	InvalidStatus       = 4011
	AlreadyHandle       = 4012
	RecordNotFound      = 4013
	InvalidScrshot      = 4014
	InvalidLink         = 4015
	InternalServerError = 5000
)

type Err struct {
	Code   int    `json:"code"`
	Result string `json:"result"`
}

func (e *Err) Error() string {
	if e == nil {
		return ""
	}
	return "code:" + strconv.Itoa(e.Code) + ", result:" + e.Result
}
func NewInvalidParamError(code int, where, parameter, details string) *Err {
	var message string
	if details != "" {
		message = ", details:" + details + ", where:" + where
	} else {
		message = ", where:" + where
	}
	return &Err{Code: code, Result: "Invalid " + parameter + " patameter" + message}
}
func NewInternalServerError(where, details string) *Err {
	var message string
	if details != "" {
		message = "details:" + details + ", where:" + where
	} else {
		message = "where:" + where
	}
	return &Err{Code: InternalServerError, Result: "Internal Server Error," + message}
}

// GetMillis is a convience method to get milliseconds since epoch.
func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
func FindConfigFile(fileName string) (path string) {
	found := FindFile(filepath.Join("conf", fileName))
	if found == "" {
		found = FindPath(fileName, []string{"."}, nil)
	}

	return found
}
func FindPath(path string, baseSearchPaths []string, filter func(os.FileInfo) bool) string {
	//判斷是否是絕對路徑
	if filepath.IsAbs(path) {
		if _, err := os.Stat(path); err == nil {
			return path
		}

		return ""
	}

	searchPaths := []string{}
	for _, baseSearchPath := range baseSearchPaths {
		searchPaths = append(searchPaths, baseSearchPath)
	}

	var binaryDir string
	//返回启动当前进程的可执行文件的路径名称。
	if exe, err := os.Executable(); err == nil {
		if exe, err = filepath.EvalSymlinks(exe); err == nil {
			if exe, err = filepath.Abs(exe); err == nil {
				binaryDir = filepath.Dir(exe)
			}
		}
	}
	if binaryDir != "" {
		for _, baseSearchPath := range baseSearchPaths {
			searchPaths = append(
				searchPaths,
				filepath.Join(binaryDir, baseSearchPath),
			)
		}
	}

	for _, parent := range searchPaths {
		found, err := filepath.Abs(filepath.Join(parent, path))
		if err != nil {
			continue
		} else if fileInfo, err := os.Stat(found); err == nil {
			if filter != nil {
				if filter(fileInfo) {
					return found
				}
			} else {
				return found
			}
		}
	}

	return ""
}

func GetLogFileLocation(fileLocation string) string {
	if fileLocation == "" {
		fileLocation, _ = FindDir("logs")
	}

	return filepath.Join(fileLocation, LOG_FILENAME)
}

func FindDir(dir string) (string, bool) {
	found := FindPath(dir, commonBaseSearchPaths, func(fileInfo os.FileInfo) bool {
		return fileInfo.IsDir()
	})
	if found == "" {
		return "./", false
	}

	return found, true
}
func FindFile(path string) string {
	return FindPath(path, commonBaseSearchPaths, func(fileInfo os.FileInfo) bool {
		return !fileInfo.IsDir()
	})
}
func checkFileExit(filename string) bool {
	var exit = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exit = false
	}
	return exit
}
func GetFile(filename string) (*os.File, error) {
	if checkFileExit(filename) {
		return os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	} else {
		index := strings.LastIndex(filename, "/")
		// 创建文件夹
		err := os.MkdirAll(filename[:index], os.ModePerm)
		if err != nil {
			return nil, err
		}
		return os.Create(filename)
	}

}
func randomString(l int) string {
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		if string(randInt(65, 90)) != temp {
			temp = string(randInt(65, 90))
			result.WriteString(temp)
			i++
		}
	}
	return result.String()
}
func RandomStr(l int) string {
	return randomString(l)
}
func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func RandNum(min int, max int) int {
	return randInt(min, max)
}

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

// Ip2long 将 IPv4 字符串形式转为 uint32
func Ip2long(ipstr string) uint32 {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}
func PaddingByZero(num int, lens int) string {
	strnum := fmt.Sprintf("%d", num)
	if len(strnum) >= lens {
		return strnum
	}
	strlens := fmt.Sprintf("%d", lens)
	tmps := "%0" + strlens + "d\n"
	return strings.TrimSpace(fmt.Sprintf(tmps, num))
}
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func aesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func DealConfDecMysql(s string) (string, error) {
	if s == "" {
		return "", errors.New("设置数据失败")
	}
	sarr := strings.Split(s, "@")
	if len(sarr) != 2 {
		return "", errors.New("主机格式不对，设置数据失败")
	}
	finarr := strings.Split(sarr[0], ":")
	if len(finarr) != 2 {
		return "", errors.New("用户格式不对，设置数据失败")
	}
	uname := UnPackEnc(finarr[0])
	upwd := UnPackEnc(finarr[1])
	if uname == "" || upwd == "" {
		return "", errors.New("数据为空，设置数据失败")
	}
	return fmt.Sprintf("%s:%s@%s", uname, upwd, sarr[1]), nil
}

func UnPackEnc(s string) string {
	s64, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Println(err.Error())
		return ""

	}

	sfin, err := aesDecrypt(s64, []byte(CONF_DES_TOKEN))
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return string(sfin)
}

func GetIpStr() string {
	ips, err := get_internal()
	fmt.Println("当前ipv4 ", ips)
	if ips == "" || err != nil {
		fmt.Println("返回ipv4字符串为随机码 ")
		return RandomStr(8)
	} else {
		arr := strings.Split(ips, ".")
		if len(arr) != 4 {
			return RandomStr(8)
		}
		res := ""
		for _, v := range arr {
			i64, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return RandomStr(8)
			}
			res += dec2Hex(i64)
		}
		fmt.Println("返回ipv4字符串为16进制码 ", res)
		if len(res) != 8 {
			return RandomStr(8)
		}
		return res
	}

}
func get_internal() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//fmt.Println(ipnet.IP.To4())
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", nil
}

func dec2Hex(n int64) string {
	if n < 0 {
		log.Println("Decimal to hexadecimal error: the argument must be greater than zero.")
		return ""
	}
	if n == 0 {
		return "00"
	}
	hex := map[int64]int64{10: 65, 11: 66, 12: 67, 13: 68, 14: 69, 15: 70}
	s := ""
	for q := n; q > 0; q = q / 16 {
		m := q % 16
		if m > 9 && m < 16 {
			m = hex[m]
			s = fmt.Sprintf("%v%v", string(m), s)
			continue
		}
		s = fmt.Sprintf("%v%v", m, s)
	}
	if len(s) == 1 {
		s = "0" + s
	}
	return s
}

// 全角转换半角
func DBCtoSBC(s string) string {
	retstr := ""
	for _, i := range s {
		inside_code := i
		if inside_code == 12288 {
			inside_code = 32
		} else {
			inside_code -= 65248
		}
		if inside_code < 32 || inside_code > 126 {
			retstr += string(i)
		} else {
			retstr += string(inside_code)
		}
	}
	return retstr
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

// 判断obj是否在target中，target支持的类型arrary,slice,map
func Contains2(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	default:
		return false

	}

	return false
}
