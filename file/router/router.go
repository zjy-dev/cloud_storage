package router

import (
    "github.com/J-Y-Zhang/cloud-storage/file/handler"
    "github.com/gin-gonic/gin"
    "net/http"
)

//cors 一招鲜解决Ajax跨域
func cors() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
        c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
        c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
        c.Header("Access-Control-Allow-Credentials", "true")

        method := c.Request.Method

        //放行所有OPTIONS方法
        if method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
        }

        c.Next()
    }
}


func Router() *gin.Engine {
    r := gin.New()
    // 解决跨域问题
    r.Use(cors())

    r.GET("file/upload/chunk", handler.TryFastAndBreakPointUpload)
    r.POST("file/upload/chunk", handler.UploadFileChunk)
    r.POST("file/upload/done", handler.UploadFileDone)
    r.GET("file/user", handler.GetFileInfosByUserName)
    r.DELETE("file/user", handler.DeleteFileByMd5)
    r.PUT("file/user", handler.UserRenameFile)
    r.GET("file/downloadurl", handler.UserDownloadFile)

    return r
}