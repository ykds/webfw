package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

type Configuration interface {
	Name() string
}


func NewConfig(config Configuration, configName, configType, configPath string) {
	if configName == "" {
		configName = "config"
	}

	if configType == "" {
		configType = "yaml"
	}

	if configPath == "" {
		configPath = "./configs"
	}

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed. err：%v", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("parse config failed. err：%v", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		err := viper.Unmarshal(config)
		if err != nil {
			log.Fatalf("reload config failed. err：%v", err)
		}
	})
}
