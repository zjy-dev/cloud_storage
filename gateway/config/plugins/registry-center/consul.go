package registry_center

import (
    "github.com/asim/go-micro/plugins/registry/consul/v4"
    log "go-micro.dev/v4/logger"
    "go-micro.dev/v4/registry"
)

func GetConsulRegistryCenter() registry.Registry {
    consulRegistryCenter := consul.NewRegistry(func(options *registry.Options) {
        options.Addrs = []string{
            ConsulHost + ":" + ConsulPort,
        }
    })

    if consulRegistryCenter == nil {
        log.Errorf("获取consul注册中心失败")
        return nil
    }

    return consulRegistryCenter
}
