package configs

import "webfw/pkg/setting"

type AppConfig struct {
	Name    string
	Version string
}

type MysqlConfig struct {
	User   string
	Pass   string
	Host   string
	DBName string
}

type RedisConfig struct {
	Host string
	Pass string
	DB   int
}

type LoggerConfig struct {
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

type RabbitMqConfig struct {
	User string
	Pass string
	Host string
	Vhost string
}

func NewConfig(configName, configType, configPath string) *Config {
	var C *Config
	setting.NewConfig(C, configName, configType, configPath)
	return C
}

type Config struct {
	App   *AppConfig
	Mysql *MysqlConfig
	Redis *RedisConfig
	Logger *LoggerConfig
	RabbitMq *RabbitMqConfig
}

func (c *Config) Name() string {
	return "configuration"
}
