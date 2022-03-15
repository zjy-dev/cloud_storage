package mq

import (
    "fmt"
    "github.com/J-Y-Zhang/cloud-storage/file/config"
    amqp "github.com/rabbitmq/amqp091-go"
)

var (
    conn *amqp.Connection
    channel *amqp.Channel
)

func initChan() bool {
    if channel != nil {
        return true
    }

    // 1.获得一个rabbitmq的连接
    conn, err := amqp.Dial(config.RabbitURL)
    if err != nil {
        fmt.Println(err)
        return false
    }
    fmt.Println("连接rabbitmq成功")

    // 2.打开一个channel, 用于消息发布
    channel, err = conn.Channel()
    if err != nil {
        return false
    }

    return true
}
