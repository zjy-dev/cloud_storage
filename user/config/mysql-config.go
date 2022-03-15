package config

import (
    "go-micro.dev/v4/config"
    log "go-micro.dev/v4/logger"
)

type MysqlConfig struct {
    Host     string `json:"host"`
    Port     string `json:"port"`
    User     string `json:"user"`
    Pwd      string `json:"pwd"`
    Database string `json:"database"`
}


func GetMysqlConfig(configCenter config.Config) *MysqlConfig{
    mysqlConf := &MysqlConfig{}

    // 从consul获取mysql的配置信息, 并存入结构体mysqlConf
    err := configCenter.Get("mysql").Scan(mysqlConf)

    if err != nil {
        log.Errorf("获取mysql配置信息失败, 错误信息 %v", err)
        return nil
    }

    return mysqlConf
}
