package oss

import (
	"fmt"
	"github.com/J-Y-Zhang/cloud-storage/file/config"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
)

var ossCli *oss.Client

//OssClient 获取Oss的Client对象
func OssClient() *oss.Client {
	if ossCli != nil {
		return ossCli
	}
	t, err := oss.New(config.OSSEndpoint, config.OSSAccesskeyID, config.OSSAccessKeySecret)
	ossCli = t
	if err != nil {
		log.Printf("连接oss失败, 错误信息 %v\n", err)
		return nil
	}
	fmt.Println("获取oss连接成功")
	return ossCli
}

//OssBucket 获取Oss的Bucket, Bucket的Name在config/oss-config.go中
func OssBucket() *oss.Bucket {
	cli := OssClient()
	if cli != nil {
		bucket, err := cli.Bucket(config.OSSBucket)
		if err != nil {
			log.Printf("获取oss的bucket失败, 错误信息 %v\n", err)
			return nil
		}
		return bucket
	}
	return nil
}

//DownloadURL 通过对象名获取下载用的URL
func DownloadURL(objName string) string {
	signedURL, err := OssBucket().SignURL(objName, oss.HTTPGet, 3600)
	if err != nil {
		log.Printf("获取下载url失败, 错误信息 %v\n", err)
		return ""
	}
	return signedURL
}

// BuildLifecycleRule : 针对指定bucket设置生命周期规则
func BuildLifecycleRule(bucketName string) {
	// 表示前缀为test的对象(文件)距最后修改时间30天后过期。
	ruleTest1 := oss.BuildLifecycleRuleByDays("rule1", "test/", true, 29)
	rules := []oss.LifecycleRule{ruleTest1}

	OssClient().SetBucketLifecycle(bucketName, rules)
}
