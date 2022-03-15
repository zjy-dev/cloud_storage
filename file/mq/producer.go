package mq

import (
    "fmt"
    amqp "github.com/rabbitmq/amqp091-go"
)



//Publish 发布消息
func Publish(exchange, routingKey string, msg []byte) bool {
    // 1.判断channel是否正常
    if !initChan() {
        return false
    }

    // 2.发布消息
    err := channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
        ContentType: "text/plain",
        Body:        msg,
    })

    if err != nil {
        fmt.Println(err)
        return false
    }

    return true
}