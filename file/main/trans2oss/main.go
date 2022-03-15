package main

import (
    "encoding/json"
    "fmt"
    "github.com/J-Y-Zhang/cloud-storage/file/config"
    "github.com/J-Y-Zhang/cloud-storage/file/db/dao"
    "github.com/J-Y-Zhang/cloud-storage/file/mq"
    "github.com/J-Y-Zhang/cloud-storage/file/oss"
    "log"
    "os"
    "path"
)

func ProcessTransfer(msg []byte) bool {
    // 1.解析msg
    data := mq.TransferData{}
    err := json.Unmarshal(msg, &data)
    if err != nil {
        log.Printf("解析rabbitmq中的msg失败, 错误信息%v\n", err)
        return false
    }

    // 2.创建文件句柄
    file, err := os.Open(data.CurLocation)
    if err != nil {
        log.Printf("打开文件失败, 错误信息 %v\n", err)
    }

    // 3.读文件并发送到oss
    err = oss.OssBucket().PutObject(data.DestLocation, file)

    if err != nil {
        log.Printf("发送到oss失败, 错误信息%v\n", err)
        return false
    }

    // 4.在文件表中更改文件存储路径
    err = dao.UpdateFileByMd5(data.FileMd5, data.DestLocation)
    if err != nil {
        log.Printf("更新文件位置到oss失败, 错误信息%v\n", err)
        return false
    }

    // 5.删除源文件
    err = os.RemoveAll(path.Dir(data.CurLocation))
    if err != nil {
        log.Printf("删除store目录下的文件失败, 错误信息%v", err)
        return false
    }

    fmt.Println("消费了一条消息")
    return true
}

func main() {
    fmt.Println("开始监听任务转移队列")
    mq.StartConsume(config.TransOSSQueueName, "transfer_oss", ProcessTransfer)
}
