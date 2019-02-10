package config

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"

	"os"
	"fmt"
	"github.com/antigloss/go/logger"
)

type Config struct {
	Postgres postgres
	Keys	keys
	loaded bool
	LogConfig	logConfig
}

type logConfig struct {
	Log string
	LogDirectory	string
	LogFilePrefix	string
	EnableLogTrace	bool
}

type keys struct {

	PrivateKeyPath	string
	PublicKeyPath	string
}

type postgres struct {
	Database string
	Username string
	Password string
	Host string
	Port string
}

var Conf Config

//InitConfig reads the configuration of the application from the config.toml file
func InitConfig(){

	viper.SetConfigName("config")
	viper.AddConfigPath("others")


	viper.SetConfigType("toml")
	err := viper.ReadInConfig()

	if err != nil {

		fmt.Println(err.Error())
		fmt.Printf("Exiting application. got error %v",err.Error())
		os.Exit(1)
	}

	err = viper.Unmarshal(&Conf)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Printf("config %v",Conf)
	Conf.loaded = true

	logFileMaxSize := 100

	err = logger.Init(Conf.LogConfig.LogDirectory,
		800,
		20,
		uint32(logFileMaxSize),
		Conf.LogConfig.EnableLogTrace,
		Conf.LogConfig.Log)

	if err != nil {
		fmt.Println("Error in intializing logger, is : ", err)
	}

	err = logger.SetFilenamePrefix(Conf.LogConfig.LogFilePrefix, Conf.LogConfig.LogFilePrefix)
	logger.SetLogThrough(false)

	if err != nil {
		fmt.Println("Error is : ", err.Error())
	}

}

//GetConfig returns the Config object
func GetConfig() (Config){

	if Conf.loaded == false{
		InitConfig()
	}

	return Conf
}