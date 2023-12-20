package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type RabbitMQ struct {
	IP       string
	User     string
	Password string
	Port     string
}

var (
	RabbitMQConfig  RabbitMQ
	RabbitMQURL     string
	ArticleExchange = "article_exchange"
	ArticleQueue    = "article_queue"
)

//初始化配置信息
func init() {
	// 设置配置文件的名字和路径（如果文件不在当前目录）
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	RabbitMQConfig = RabbitMQ{
		IP:       viper.GetString("rabbitmq.IP"),
		User:     viper.GetString("rabbitmq.User"),
		Password: viper.GetString("rabbitmq.Password"),
		Port:     viper.GetString("rabbitmq.Port"),
	}
	RabbitMQURL = fmt.Sprintf("amqp://%v:%v@%v:%v/", RabbitMQConfig.User, RabbitMQConfig.Password, RabbitMQConfig.IP, RabbitMQConfig.Port)
}
