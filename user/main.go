package main

import (
	"go-micro.dev/v4"
	log "go-micro.dev/v4/logger"
	"user/config"
	"user/handler"
	config_center "user/plugins/config-center"
	registry_center "user/plugins/registry-center"

	pb "user/proto"
)

func main() {
	// 1.创建User微服务, 同时初始化配置中心, 服务发现/注册中心等插件
	srv := micro.NewService(
		micro.Name(config.ServiceName),
		micro.Version(config.ServiceVersion),
		micro.Address(config.ServiceAddress),
		// 配置中心
		micro.Config(config_center.GetConsulConfigCenter()),
		// 服务发现/注册中心
		micro.Registry(registry_center.GetConsulRegistryCenter()),
	)

	// 2.接收命令行参数初始化服务, 本项目中基本不需要
	srv.Init()

	// 3.为该服务注册handler
	pb.RegisterUserHandler(srv.Server(), new(handler.User))

	// 4.运行服务
	if err := srv.Run(); err != nil {
		log.Error(err)
	}
}
