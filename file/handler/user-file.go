package handler

import (
    "fmt"
    "github.com/J-Y-Zhang/cloud-storage/file/ErrorCode"
    "github.com/J-Y-Zhang/cloud-storage/file/db/dao"
    "github.com/J-Y-Zhang/cloud-storage/file/db/model"
    "github.com/J-Y-Zhang/cloud-storage/file/oss"
    "github.com/J-Y-Zhang/cloud-storage/file/util"
    "github.com/gin-gonic/gin"
    "log"
    "net/http"
    "strconv"
    "strings"
)

//GetFileInfosByUserName 用户查询他创建的所有文件
func GetFileInfosByUserName(c *gin.Context) {
    userName := c.Query("username")
    files, err := dao.UserQueryFiles(userName)
    if err != nil || files == nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": ErrorCode.InvalidParams,
            "msg": "用户不存在",
        })
        return
    }

    type respFile struct {
        *model.UserFile
        TypedFileSize string
    }

    respFiles := make([]respFile, len(files))
    for i := range files {
        respFiles[i] = respFile{
            UserFile:          files[i],
        }
        size := float64(files[i].FileSize)
        if size < 100 {
            respFiles[i].TypedFileSize = strconv.Itoa(int(size)) + "B"
        } else if size < 100 * util.KB {
            respFiles[i].TypedFileSize = util.FomatFloat64(float64(size) / util.KB, 2) + "KB"
        } else if size < 100 * util.MB{
            respFiles[i].TypedFileSize = util.FomatFloat64(float64(size) / util.MB, 2) + "MB"
        } else {
            respFiles[i].TypedFileSize = util.FomatFloat64(float64(size) / util.GB, 2) + "GB"
        }

    }

    c.JSON(http.StatusOK, gin.H{
        "code": ErrorCode.Success,
        "files": respFiles,
    })
}

//DeleteFileByMd5 用户删除文件
func DeleteFileByMd5(c *gin.Context) {
    filemd5 := c.Query("filemd5")
    uname := c.Query("username")

    err := dao.UserDeleteFileByMd5(uname, filemd5)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": ErrorCode.UnknownErr,
            "msg": "删除文件失败",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code": ErrorCode.Success,
        "msg": "删除文件成功",
    })
}

//UserRenameFile 用户重命名文件
func UserRenameFile(c *gin.Context) {
    filemd5 := c.Query("filemd5")
    uname := c.Query("username")
    newname := c.Query("newname")

    err := dao.UserUpdateFileName(uname, filemd5, newname)
    if err != nil {
        log.Printf("%v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "code": ErrorCode.UnknownErr,
            "msg": "重命名文件失败",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code": ErrorCode.Success,
        "msg": "删除文件成功",
    })
}

//UserDownloadFile 用户下载文件
func UserDownloadFile(c *gin.Context) {
    fileMd5 := c.Query("filemd5")

    // 1.查询文件表获取文件位置
    fileInfo, err := dao.FindFileByMd5(fileMd5)
    if err != nil {
        log.Printf("查询文件表失败, 错误信息 %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "code": ErrorCode.UnknownErr,
            "msg": "获取下载链接失败",
        })
        return
    }

    // 2.如果还未存到oss, 则从服务器下载
    if !strings.HasPrefix(fileInfo.FilePath, "oss/") {
        c.Header("Content-Type", "application/octet-stream")
        c.Header("Content-Disposition", "attachment; filename=" + fileInfo.FileName)
        c.Header("Content-Transfer-Encoding", "binary")

        fmt.Println(fileInfo.FilePath)
        
        c.File(fileInfo.FilePath)
        return
    }

    // 3.如果存到了oss, 则返回oss的下载url给前端
    downloadURL := oss.DownloadURL(fileInfo.FilePath)
    if downloadURL == "" {
        log.Printf("获取oss下载链接失败\n")
        c.JSON(http.StatusBadRequest, gin.H{
            "code": ErrorCode.UnknownErr,
            "msg": "获取下载链接失败",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "code": ErrorCode.Success,
        "downloadURL": downloadURL,
    })
}
