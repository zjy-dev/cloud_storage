package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/J-Y-Zhang/cloud-storage/file/ErrorCode"
	"github.com/J-Y-Zhang/cloud-storage/file/algorithm"
	"github.com/J-Y-Zhang/cloud-storage/file/cache/redis"
	"github.com/J-Y-Zhang/cloud-storage/file/config"
	"github.com/J-Y-Zhang/cloud-storage/file/db/dao"
	"github.com/J-Y-Zhang/cloud-storage/file/db/model"
	"github.com/J-Y-Zhang/cloud-storage/file/mq"
	"github.com/J-Y-Zhang/cloud-storage/file/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
)

var mux =  sync.Mutex{}

//UploadFileChunk 分块上传文件
func UploadFileChunk(c *gin.Context) {
	// 1.获取表单数据以及初始化redis连接
	fileName := c.PostForm("filename")
	fileMd5 := c.PostForm("identifier")
	relativePath := c.PostForm("relativePath")

	chunkIndex, _ := strconv.Atoi(c.PostForm("chunkNumber"))
	curChunkSize, _ := strconv.Atoi(c.PostForm("currentChunkSize"))
	fileSize, _ := strconv.Atoi(c.PostForm("totalSize"))
	chunkCnt, _ := strconv.Atoi(c.PostForm("totalChunks"))

	// 获取redis连接(别忘了关闭)
	rConn := redis.RedisPool().Get()
	defer rConn.Close()

	// 拦截重复请求
	UploadID := "CU_" + fileMd5
	isChunkExist, _ := rConn.Do("HEXISTS", UploadID, "CI_" + strconv.Itoa(chunkIndex))
	if isChunkExist.(int64) == 1 {
		c.JSON(http.StatusOK, gin.H{
			"code": ErrorCode.Success,
			"msg": fmt.Sprintf("文件%v的第%v个分块上传成功", fileName, chunkIndex),
		})
		return
	}

	// 2.拷贝文件分块

	// 创建文件夹
	chunkPath := "tmp/" + fileMd5 + "/" + relativePath + "-" + strconv.Itoa(chunkIndex)
	err := os.MkdirAll(path.Dir(chunkPath), os.ModePerm)
	if err != nil {
		log.Printf("创建文件夹失败, 错误信息 %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": ErrorCode.InternelServerError,
			"msg": "上传分块错误",
		})
		return
	}

	// 接收文件分块并拷贝
	chunk, err := c.FormFile("file")
	if err != nil {
		log.Printf("接收文件错误, 错误信息 %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": ErrorCode.InvalidParams,
			"msg": "文件传输失败",
		})
		return
	}

	if int(chunk.Size) != curChunkSize {
		log.Printf("文件受损")
		c.JSON(http.StatusBadRequest, gin.H{
			"code": ErrorCode.InvalidParams,
			"msg": "文件受损",
		})
		return
	}

	// 拷贝
	err = c.SaveUploadedFile(chunk, chunkPath)
	if err != nil {
		log.Printf("保存文件失败 %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": ErrorCode.InvalidParams,
			"msg": "文件保存失败",
		})
		return
	}

	// 3.将分块信息存入redis

	// 如果文件总信息还未初始化, 则初始化一下
	isExist, _ := rConn.Do("EXISTS", UploadID)

	for isExist.(int64) == 0 {
		mux.Lock()
		defer mux.Unlock()

		isExist, _ = rConn.Do("EXISTS", UploadID)
		if isExist.(int64) == 1 {
			break
		}

		rConn.Do("HSET", UploadID, "fileMd5", fileMd5)
		rConn.Do("HSET", UploadID, "chunkCnt", chunkCnt)
		rConn.Do("HSET", UploadID, "fileSize", fileSize)
		rConn.Do("HSET", UploadID, "filePath", path.Dir(chunkPath))
		rConn.Do("HSET", UploadID, "relativePath", relativePath)
		rConn.Do("HSET", UploadID, "fileName", fileName)

		isExist, _ = rConn.Do("EXISTS", UploadID)
	}

	// 更新分块信息, 接收到的分块数量和接收到的分块大小
	rConn.Do("HSET", UploadID, "CI_" + strconv.Itoa(chunkIndex), 1)
	rConn.Do("HINCRBY" +
		"", UploadID, "curSize", curChunkSize)
	rConn.Do("HINCRBY", UploadID, "curCnt", 1)

	// 4.返回上传分块成功
	c.JSON(http.StatusOK, gin.H{
		"code": ErrorCode.Success,
        "msg": fmt.Sprintf("文件%v的第%v个分块上传成功", fileName, chunkIndex),
    })
}

//TryFastAndBreakPointUpload 秒传和断点续传
func TryFastAndBreakPointUpload(c *gin.Context) {
	// 1.解析url参数
	fileMd5 := c.Query("identifier")
	chunkCnt, _ := strconv.Atoi(c.Query("totalChunks"))

	// 2.尝试秒传
	file, _ := dao.FindFileByMd5(fileMd5)
	if file.FileName != "" {
		fmt.Println("秒传, 用户已经爽飞了")
		c.JSON(http.StatusOK, gin.H{
			"code": ErrorCode.Success,
			"isExist": true,
		})
		return
	}

	// 3.尝试断点续传

	// 获取redis连接(别忘了关闭)
	rConn := redis.RedisPool().Get()
	defer rConn.Close()
	UploadID := "CU_" + fileMd5
	idx := 1
	for ; idx < chunkCnt; idx++ {
		isChunkExist, _ := rConn.Do("HEXISTS", UploadID, "CI_" + strconv.Itoa(idx))
		if isChunkExist.(int64) == 0 {
			break
		}
	}

	uploaded := make([]int, idx - 1)
	for i := 1; i < idx; i++ {
		uploaded[i - 1] = i
	}

	c.JSON(http.StatusOK, gin.H{
		"code": ErrorCode.Success,
		"isExist": false,
		"uploaded": uploaded,
	})

}

//UploadFileDone 所有分块上传完毕后合并文件, 将文件转移到oss的请求说写入rabbitmq, 并清除redis缓存和文件分片
func UploadFileDone(c *gin.Context) {
	// 1.获取表单数据
	userName := c.PostForm("username")
	fileName := c.PostForm("filename")
	fileMd5 := c.PostForm("identifier")
	fileSize, _ := strconv.Atoi(c.PostForm("totalSize"))
	chunkCnt, _ := strconv.Atoi(c.PostForm("totalChunks"))

	// 2.与redis存储的信息做校验

	// 获取redis连接(别忘了关闭)
	rConn := redis.RedisPool().Get()
	defer rConn.Close()
	UploadID := "CU_" + fileMd5

	isExist, _ := rConn.Do("EXISTS", UploadID)
	if isExist.(int64) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": ErrorCode.Success,
			"msg": "秒传",
		})

		dao.UserAddFile(&model.UserFile{
			UserName: userName,
			FileID:   fileMd5,
			FileName: fileName,
			FileSize: fileSize,
		})

		return
	}


	rCurCnt, _ := rConn.Do("HGET", UploadID, "curCnt")
	rCurSize, _ := rConn.Do("HGET", UploadID, "curSize")
	rFileSize, _ := rConn.Do("HGET", UploadID, "fileSize")
	rFilePath, _ := rConn.Do("HGET", UploadID, "filePath")
	rFileName, _ := rConn.Do("HGET", UploadID, "fileName")
	rChunkCnt, _ := rConn.Do("HGET", UploadID, "chunkCnt")
	//relativePath, _ := rConn.Do("HGET", UploadID, "relativePath")

	// 最后要清一下redis和删除临时缓存
	defer func() {
		p, ok := rFilePath.([]byte)
		if ok {
			err := os.RemoveAll(string(p))
			if err != nil {
				log.Printf("删除文件夹%v失败", string(rFilePath.([]byte)))
			}
		}

		_, err := rConn.Do("DEL", UploadID)
		if err != nil {
			log.Printf("清除Redis缓存(key = %v)失败", UploadID)
		}

	}()

	r_cur_cnt, _ := strconv.Atoi(string(rCurCnt.([]byte)))
	r_chunk_cnt, _ := strconv.Atoi(string(rChunkCnt.([]byte)))
	r_cur_size, _ := strconv.Atoi(string(rCurSize.([]byte)))
	r_file_size, _ := strconv.Atoi(string(rFileSize.([]byte)))

	if !util.CompareInt(r_cur_cnt, r_chunk_cnt, chunkCnt) ||
		!util.CompareInt(r_cur_size, r_file_size, fileSize) ||
		!bytes.Equal(rFileName.([]byte), []byte(fileName)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": ErrorCode.UnknownErr,
			"msg": "文件校验失败, 请重新上传",
		})
		return
	}

	// 3.合并文件
	storePath := "store/" + fileMd5 + "/" + fileName
	algorithm.MergeFile(string(rFilePath.([]byte)), fileName, storePath, chunkCnt)

	// 4.更新文件表和用户文件表
	dao.AddFile(&model.File{
		FileName:     fileName,
		FileMd5:      fileMd5,
		FilePath:	  storePath,
		FileSize:     fileSize,
	})

	dao.UserAddFile(&model.UserFile{
		UserName: userName,
		FileID:   fileMd5,
		FileName: fileName,
		FileSize: fileSize,
	})

	// 5.将转移请求写入rabbitmq, 后台会异步接收, 然后将文件转移到oss
	ossPath := "oss/" + fileMd5 + "/" + fileName
	transData := mq.TransferData{
		FileMd5:     fileMd5,
		CurLocation:  storePath,
		DestLocation: ossPath,
	}
	jsonedTransData, _ := json.Marshal(transData)
	mq.Publish(config.TransExchangeName, config.TransOSSRoutingKey, jsonedTransData)

	c.JSON(http.StatusOK, gin.H{
		"code": ErrorCode.Success,
		"msg": "上传文件成功",
	})
}
