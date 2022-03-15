package dao

import (
    "github.com/J-Y-Zhang/cloud-storage/common/config-center"
    log "go-micro.dev/v4/logger"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "user/config"
)

var mysqlConn *gorm.DB = nil

func GetMysqlConn() *gorm.DB {
    // 单例模式
    if mysqlConn != nil {
        return mysqlConn
    }

    // 1.获取配置中心
    configCenter, err := config_center.GetConsulConfigCenter()
    if err != nil {
        log.Errorf("获取配置中心失败, 错误信息 ", err)
        return nil
    }

    // 2.从配置中心获取mysql配置信息
    mysqlConf := config.GetMysqlConfig(configCenter)
    if mysqlConf == nil {
        return nil
    }

    // 3.根据mysql配置信息去连接mysql
    mysqlStr := mysqlConf.User + ":" + mysqlConf.Pwd + "@tcp(" + mysqlConf.Host +
        ":" + mysqlConf.Port + ")/" + mysqlConf.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
    mysqlConn, err = gorm.Open(mysql.Open(mysqlStr), &gorm.Config{})
    if err != nil {
        log.Errorf("连接mysql数据库失败, 错误信息 ", err)
        return nil
    }

    return mysqlConn
}
