package mq

import "fmt"

var done chan bool = make(chan bool, 0)

//StartConsume 开始监听
func StartConsume(queueName, consumerName string, callback func(msg []byte) bool) {
    initChan()

    // 1.通过channel.Consume获得消息信道
    msgCh, err := channel.Consume(queueName, consumerName, true, false, false, false, nil)
    if err != nil {
        fmt.Println(err)
        return
    }

    // 2.循环获取新消息
    go func() {
        for msg := range msgCh {
            ok := callback(msg.Body)
            if !ok {
                fmt.Println("写入到oss失败")
                // TODO: 写到另一个异常处理队列中去
            }
        }
    }()

    // 阻塞
    <-done

    channel.Close()
}
