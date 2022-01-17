package configs

import (
	"sync"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Port     int `yaml:"port"`
	Database struct {
		Driver   string `yaml:"driver"`
		Name     string `yaml:"name"`
		Address  string `yaml:"address"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
}

var lock = &sync.Mutex{}
var appConfig *AppConfig

func GetConfig(env string) *AppConfig {
	lock.Lock()
	defer lock.Unlock()

	if appConfig == nil {
		appConfig = initConfig(env)
	}

	return appConfig
}

func initConfig(env string) *AppConfig {
	var defaultConfig AppConfig
	defaultConfig.Port = 8000
	defaultConfig.Database.Driver = "mysql"
	defaultConfig.Database.Address = "localhost"
	defaultConfig.Database.Port = 3306
	defaultConfig.Database.Username = "root"
	defaultConfig.Database.Password = ""

	if env == "test" {
		defaultConfig.Database.Name = "altastore"
	} else {
		defaultConfig.Database.Name = "altastore"
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs/")

	if err := viper.ReadInConfig(); err != nil {
		return &defaultConfig
	}

	var finalConfig AppConfig
	err := viper.Unmarshal(&finalConfig)
	if err != nil {
		log.Info("failed to extract config, will use default value")
		return &defaultConfig
	}

	return &finalConfig
}