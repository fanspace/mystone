package settings

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var VOptions *viper.Viper

func initOptions() {
	log.Println("加载自定义配置文件......")
	VOptions = viper.New()
	VOptions.SetConfigName("options.yaml")
	VOptions.AddConfigPath("./conf")
	VOptions.SetConfigType("yaml")
	if err := VOptions.ReadInConfig(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

}
