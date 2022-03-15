package dao

import (
    "github.com/J-Y-Zhang/cloud-storage/file/config"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "log"
)

var mysqlConn *gorm.DB = nil

func GetMysqlConn() *gorm.DB {
    // 单例模式
    if mysqlConn != nil {
        return mysqlConn
    }

    mysqlStr := config.MysqlUser + ":" + config.MysqlPwd + "@tcp(" + config.MysqlHost +
        ":" + config.MysqlPort + ")/" + config.MysqlDateBase + "?charset=utf8mb4&parseTime=True&loc=Local"

    mysqlConn, err := gorm.Open(mysql.Open(mysqlStr), &gorm.Config{})
    if err != nil {
        log.Printf("连接mysql数据库失败, 错误信息 %v\n", err)
        return nil
    }

    return mysqlConn
}
