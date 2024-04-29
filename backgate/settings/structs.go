package settings

type Config struct {
	ReleaseMode     bool
	DaprMode        bool
	Smark           string
	AppName         string
	DbName          string
	ServiceSettings *ServiceSettings
	Database        *Database
	CorsSettings    *CorsSettings
	LogSettings     *LogSettings
	//RpcSettings		*RpcSettings
	GrpcSettings     map[string]string
	RabbitMqSettings *RabbitMqSettings
	MinioSettings    *MinioSettings
}
type Database struct {
	MysqlSettings  *MysqlSettings
	RedisSettings  *RedisSettings
	DamengSettings *DamengSettings
}
type MysqlSettings struct {
	DriverName   string
	Url          string
	MaxIdleConns int32
	MaxOpenConns int32
	QueryTimeout int32
}

type RedisSettings struct {
	Addr        string
	DB          int
	Password    string
	PoolSize    int
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
}
type LogSettings struct {
	EnableConsole bool
	ConsoleLevel  string
	ConsoleJson   *bool
	EnableFile    bool
	FileLevel     string
	FileJson      *bool
	FileLocation  string
}
type CorsSettings struct {
	Allows []string
}
type ServiceSettings struct {
	ListenAddress string
	ReadTimeout   int
	WriteTimeout  int
}
type FtpSettings struct {
	Ftpurl  string
	Ftpuser string
	Ftppwd  string
}
type RpcSettings struct {
	//ResRpc       string
	ArticleRpc string
	CourseRpc  string
	ExamRpc    string
	InterRpc   string
	LiveApi    string
	AccountRpc string
	SysRpc     string
	CloudRpc   string
	LogRpc     string
	CronRpc    string
	OrgRpc     string
	FundRpc    string
	OfficeRpc  string
	MsgRpc     string
	NoticeRpc  string
	CheckRpc   string
	ReportRpc  string
	WxmpRpc    string
}
type DamengSettings struct {
	DriverName string
	Url        string
}
type RabbitMqSettings struct {
	Url  string
	Port int32
	User string
	Pwd  string
}

type MinioSettings struct {
	KeyID         string
	AccessKey     string
	PublicBucket  string // 公用
	PrivateBucket string // 私有
	EndPoint      string
	EntryPoint    string
}

var Cfg *Config = &Config{}
