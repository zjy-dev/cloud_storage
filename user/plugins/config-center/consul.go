package config_center

import (
	"github.com/asim/go-micro/plugins/config/source/consul/v4"
	"go-micro.dev/v4/config"
	log "go-micro.dev/v4/logger"
)

func GetConsulConfigCenter() config.Config {
	// 1.获取一个consul源
	consulSource := consul.NewSource(
		// Consul的地址, 通过配置文件实现解耦
		consul.WithAddress(ConsulHost+":"+ConsulPort),

		// Consul可能同时作为很多应用的配置中心, 所以要有前缀来区分
		consul.WithPrefix(ConsulPrefix),

		// 不加前缀也能获取value
		consul.StripPrefix(true),
	)

	// 2.获取一个通用的配置中心
	conf, err := config.NewConfig()

	if err != nil {
		log.Errorf("获取通用配置中心失败, err %v", err)
		return nil
	}

	// 3.将consul源载入这个配置中心
	err = conf.Load(consulSource)
	if err != nil {
		log.Errorf("载入consul配置中心失败, err %v", err)
		return nil
	}

	return conf
}
