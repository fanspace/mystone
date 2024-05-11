package relations

var (
	DOMAINS_LIMITED = [3]string{"foreend", "backend", "com"}
	GRPC_CREDENTIAL = map[string]string{
		"Ab365FDG": "9aBcDeFgHiJkLmNoPqRsTuVwXyZz123",
		"GSER3exs": "456789abCDefGhIjKlMnOpQrStUvWxYz",
	}
)

const (
	APP_NAME                   = "stonebackos"
	RED_CASBIN_CHANNEL_NAME    = "stonebackos_cas"
	RED_USER_BAN_CHANNEL_NAME  = "stonebackos_ban"
	APISITE_PREFIX             = "mgmt"
	JWT_SECRET_STRING_DEV      = "abcDef0135702468"
	JWT_SECRET_STRING_PROD_NOR = "aKamw02DfDlmC9vXzFiNHNscAyMi0wNS"
	JWT_SECRET_STRING_PROD_COM = "9Ag30dgh2sf14a3386b6105e90ed81at"
	JWT_SECRET_STRING_PROD_MAN = "bwe34saf63a593f70de7f3338581c2dd"
	WEB_STATUS_BACK            = 20000 // equlas http  code 200

	DEVICE_PC    = 1
	DEVICE_H5    = 2
	DEVICE_WX    = 3
	DEVICE_WXPRO = 4

	WX_SENCE_QRCODE_URL = "https://mp.weixin.qq.com/cgi-bin/showqrcode"
	ZERO                = 0
	Go_Time_Temp        = "2006-01-02T15 = 04 = 05Z"

	DICT_OPTIONS_CONST_IDTYPE = "idtype"

	FILE_404                 = "zjpub/assets/images/404.png"
	ARTICLE_SECTION_TRAINING = 1
	ARTICLE_SECTION_OA       = 2
	ARTICLE_SECTION_VERIFY   = 3

	Unit_Cate_Normal = 2
	Unit_Cate_Base   = 4
	Unit_Cate_Ind    = 8
	Unit_Cate_Gov    = 16

	PAYMENT_NAME_ALI = "ali"
	PAYMENT_NAME_WX  = "wx"

	PAYMENT_PLATFORM_ALI_WEB = 101
	PAYMENT_PLATFORM_ALI_H5  = 102
	PAYMENT_PLATFORM_ALI_APP = 103

	PAYMENT_PLATFORM_WX_WEB  = 201
	PAYMENT_PLATFORM_WX_H5   = 202
	PAYMENT_PLATFORM_WX_PROG = 203
	PAYMENT_PLATFORM_WX_APP  = 204

	PAYMENT_PLATFORM_APPLE_APP = 301
	MAIL_TYPE_REBACK_PWD       = 1
	MAIL_TYPE_NOTICE           = 2

	SSM_TYPE_REBACK_PWD = 1
	DEFAULT_PAGESIZE    = 20

	UNIFORM_TOKEN = "FrMfmBF9Us2Y98wH"

	STR_PREFIX    = "backos"
	STR_LOGIN_ERR = "loginerr"
	STR_BAN_USER  = "banuser"

	RABBIT_USER_ROUTING_KEY = "abcdefcc"
)
