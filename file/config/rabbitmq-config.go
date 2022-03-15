package config

const (
    // RabbitURL : rabbitmq服务的入口url
    RabbitURL = "<YourRabbitURL>"

    // TransExchangeName : 用于文件transfer的交换机
    TransExchangeName = "uploadserver.trans"

    // TransOSSQueueName : oss转移队列名
    TransOSSQueueName = "uploadserver.trans.oss"

    // TransOSSRoutingKey : routingkey
    TransOSSRoutingKey = "oss"
)
