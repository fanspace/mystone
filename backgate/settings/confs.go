package settings

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"os"
)

func InitCfg() {
	log.Println("加载系统配置文件......config.yaml")
	vs := viper.New()
	vs.SetConfigName("config.yaml")
	vs.AddConfigPath("./conf")
	vs.SetConfigType("yaml")
	if err := vs.ReadInConfig(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	if err := vs.Unmarshal(Cfg); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	// 以下为判断
	if Cfg.Smark == "" {
		Cfg.Smark = GetIpStr()
	}
	// database
	if Cfg.ReleaseMode {
		mysqlstr, err := DealConfDecMysql(Cfg.Database.MysqlSettings.Url)
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}
		Cfg.Database.MysqlSettings.Url = mysqlstr

		redispwd := UnPackEnc(Cfg.Database.RedisSettings.Password)
		if redispwd == "" {
			log.Println("读取用户数据失败")
			os.Exit(1)
		}
		Cfg.Database.RedisSettings.Password = redispwd
		// 缓存和持久化分开放;
		/*redispwd2 :=UnPackEnc(Cfg.RedisSecondSettings.Password)
		if redispwd2 == "" {
			log.Error("读取用户数据失败")
			os.Exit(1)
		}
		Cfg.RedisSecondSettings.Password = redispwd2*/

		//  rabbitmq

		rmuser := UnPackEnc(Cfg.RabbitMqSettings.User)
		if rmuser == "" {
			log.Println("Readding RabbitMq config failed")
			//os.Exit(1)
		}
		Cfg.RabbitMqSettings.User = rmuser
		rmpwd := UnPackEnc(Cfg.RabbitMqSettings.Pwd)
		if rmpwd == "" {
			log.Println("Readding RabbitMq config failed")
			//os.Exit(1)
		}
		Cfg.RabbitMqSettings.Pwd = rmpwd
		// minio
		minioacid := UnPackEnc(Cfg.MinioSettings.KeyID)
		if minioacid == "" {
			log.Println("Readding Minio config failed")
		}
		Cfg.MinioSettings.KeyID = minioacid
		minioackey := UnPackEnc(Cfg.MinioSettings.AccessKey)
		if minioackey == "" {
			log.Println("Readding Minio config failed")
		}
		Cfg.MinioSettings.AccessKey = minioackey
		gin.SetMode(gin.ReleaseMode)
	}
	initOptions()
	Cfg.PrintConfig()

}

func (s *LogSettings) PrintLogConfig() {
	log.Println(fmt.Sprintf("EnableConsole:%v", s.EnableConsole))
	log.Println(fmt.Sprintf("ConsoleLevel:%v", s.ConsoleLevel))
	log.Println(fmt.Sprintf("ConsoleJson:%v", *s.ConsoleJson))
	log.Println(fmt.Sprintf("EnableFile:%v", s.EnableFile))
	log.Println(fmt.Sprintf("FileLevel:%v", s.FileLevel))
	log.Println(fmt.Sprintf("FileJson:%v", *s.FileJson))
	log.Println(fmt.Sprintf("FileLocation:%v", s.FileLocation))
}

func (o *Config) PrintConfig() {
	o.LogSettings.PrintLogConfig()
	log.Println(fmt.Sprintf("ReleaseMode:%v", o.ReleaseMode))
}
