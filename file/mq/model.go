package mq


// 写入到RabbitMq中的数据结构
type TransferData struct {
    FileMd5 string
    CurLocation string
    DestLocation string
}

