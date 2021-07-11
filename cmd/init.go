package cmd

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"webfw/configs"
	"webfw/global"
	"webfw/pkg/logger"
	"webfw/pkg/mysql"
	"webfw/pkg/rabbitmq"
	"webfw/pkg/redis"
)

func InitApp() {
	app := &cli.App{
		Name:                 "webfw",
		Usage:                "web framework",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "configName",
				Aliases: []string{"cn"},
				Usage:   "配置文件名",
				Value:   "config",
			},
			&cli.StringFlag{
				Name:    "configType",
				Aliases: []string{"ct"},
				Usage:   "配置文件类型",
				Value:   "yaml",
			},
			&cli.StringFlag{
				Name:    "configPath",
				Aliases: []string{"cp"},
				Usage:   "配置文件路径",
				Value:   "./configs",
			},
		},
		Action: func(c *cli.Context) error {
			configName := c.String("configName")
			configType := c.String("configType")
			configPath := c.String("configPath")
			return func(configName, configType, configPath string) error {
				global.C = configs.NewConfig(configName, configType, configPath)

				global.L = logger.NewLogger(
					global.C.Logger.Filename,
					global.C.Logger.MaxSize,
					global.C.Logger.MaxBackups,
					global.C.Logger.MaxAge,
					global.C.Logger.Compress,
				)

				global.DB = mysql.NewMysql(
					global.C.Mysql.User,
					global.C.Mysql.Pass,
					global.C.Mysql.Host,
					global.C.Mysql.DBName,
				)

				global.Redis = redis.NewRedis(
					global.C.Redis.Host,
					global.C.Redis.Pass,
					global.C.Redis.DB,
				)

				rabbitmq.NewRabbitMq(
					global.C.RabbitMq.User,
					global.C.RabbitMq.Pass,
					global.C.RabbitMq.Host,
					global.C.RabbitMq.Vhost,
				)
				return nil
			}(configName, configType, configPath)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
