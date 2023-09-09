package main

import (
	"fmt"
	"go.uber.org/zap"
	"web_app/controller"
	"web_app/dao/mysql"
	"web_app/dao/rabbitMQ"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/model"
	"web_app/routes"
	"web_app/serve"
	"web_app/settings"
)

// Go Web开发较通用脚手架模板
func main() {

	// 加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed,err:%v\n", err)
		return
	}
	// 初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.ProjectConfig.Mode); err != nil {
		fmt.Printf("init logger failed,err:%v\n", err)
		return
	}
	//初始化雪花id
	if err := controller.Init("2023-03-17", 1); err != nil {
		fmt.Printf("init snow failed,err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")

	// 初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed,err:%v\n", err)
		return
	}
	//初始化连接redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed ,err:%v\n", err)
		zap.L().Error(fmt.Sprintf("init redis failed ,err:%v\n", err))
		return
	}
	//初始化RabbitMQ
	if err := rabbitMQ.Init(settings.Conf.RabbitMQConfig); err != nil {
		fmt.Printf("init RabbitMQ failed ,err:%v\n", err)
		zap.L().Error(fmt.Sprintf("init RabbitMQ failed ,err:%v\n", err))
		return
	}
	var err error
	//err = mysql.CreateTable()
	//if err != nil {
	//	return
	//}
	//开启一个协程用于聊天服务
	go serve.Start(&model.Manager)
	// 注册路由
	r := routes.Setup()
	err = r.Run(fmt.Sprintf(":%d", settings.Conf.ProjectConfig.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
	//开启一个RabbitMQ的消费者

}
