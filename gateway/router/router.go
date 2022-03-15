package router

import (
    "github.com/J-Y-Zhang/cloud-storage/gateway/handler/user"
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
    r := gin.Default()
    // 解决跨域问题
    r.Use(cors())

    r.POST("/user/signup", user.UserSignUp)
    r.POST("/user/signin", user.UserSignIn)
    r.GET("/user/info", user.GetUserInfoByName)

    return r
}
