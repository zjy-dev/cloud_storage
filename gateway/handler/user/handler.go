package user

import (
    "context"
    "fmt"
    "github.com/J-Y-Zhang/cloud-storage/gateway/ErrorCode"
    registry_center "github.com/J-Y-Zhang/cloud-storage/gateway/config/plugins/registry-center"
    "github.com/J-Y-Zhang/cloud-storage/gateway/config/services"
    pb "github.com/J-Y-Zhang/cloud-storage/gateway/handler/user/proto"
    "github.com/gin-gonic/gin"
    "go-micro.dev/v4"
    log "go-micro.dev/v4/logger"
    "net/http"
)

var (
    microUserClient pb.UserService
)

// 加载micro的User客户端
func init()  {
    srv := micro.NewService(
        // 注册中心, 有它才能发现服务端的地址
        micro.Registry(registry_center.GetConsulRegistryCenter()),
    )

    srv.Init()
    microUserClient = pb.NewUserService(services.UserServiceName, srv.Client())
    log.Infof("初始化User微服务客户端成功")
}

//UserSignIn 用户登录
func UserSignIn(c *gin.Context) {
    uname := c.PostForm("username")
    upwd := c.PostForm("password")

    // 1.请求登录服务
    userSignInResp, err := microUserClient.SignIn(context.Background(), &pb.UserSignInReq{
        UserName: uname,
        UserPwd:  upwd,
    })

    // 2.错误处理
    if userSignInResp != nil && !userSignInResp.IsUserExist{
        log.Errorf("用户" + uname + "不存在")

        c.JSON(http.StatusBadRequest, gin.H{
            "code": ErrorCode.InvalidParams,
            "msg":  "用户" + uname + "不存在",
        })

        return
    }

    if userSignInResp != nil && userSignInResp.IsPwdError{
        log.Errorf("用户" + uname + "密码错误")

        c.JSON(http.StatusBadRequest, gin.H{
            "code": ErrorCode.InvalidParams,
            "msg":  "用户" + uname + "密码错误",
        })

        return
    }

    if userSignInResp == nil || err != nil{
        if userSignInResp == nil {
            log.Errorf("未知错误")
        } else {
            log.Errorf("未知错误, 错误信息 %v", err)
        }

        c.JSON(http.StatusBadRequest, gin.H{
            "code": ErrorCode.InternelServerError,
            "msg":  "服务器内部错误, 抱歉啦",
        })

        return
    }


    // 3.登录成功
    c.JSON(http.StatusOK, gin.H{
        "code": ErrorCode.Success,
        "msg":  "用户" + uname + "登录成功",
    })
}

//UserSignUp 用户注册
func UserSignUp(c *gin.Context) {
    uname := c.PostForm("username")
    upwd := c.PostForm("password")
    nickName := c.PostForm("nickname")

    // 1.请求注册服务
    userSignUpResp, err := microUserClient.SignUp(context.Background(), &pb.UserSignUpReq{
        UserName: uname,
        UserPwd:  upwd,
        NickName: nickName,
    })

    // 2.错误处理
    if userSignUpResp != nil && userSignUpResp.IsUserExist{
        log.Errorf("用户" + uname + "已存在")

        c.JSON(http.StatusBadRequest, gin.H{
            "code": ErrorCode.InvalidParams,
            "msg":  "用户" + uname + "已存在, 请换一个名字哦",
        })

        return
    }

    if userSignUpResp == nil || err != nil {
        if userSignUpResp == nil {
            log.Errorf("未知错误")
        } else {
            log.Errorf("未知错误, 错误信息 %v", err)
        }

        c.JSON(http.StatusBadRequest, gin.H{
            "code": ErrorCode.InternelServerError,
            "msg":  "服务器内部错误, 抱歉啦",
        })

        return
    }

    // 3.注册成功
    c.JSON(http.StatusOK, gin.H{
        "code": ErrorCode.Success,
        "msg":  "用户" + uname + "注册成功, 记好密码哦",
    })
}

//GetUserInfoByName 通过用户名获取用户信息
func GetUserInfoByName(c *gin.Context) {
    uname := c.Query("username")
    userInfo, err := microUserClient.GetUserInfo(context.Background(), &pb.UserInfoReq{
        UserName: uname,
    })

    if err != nil || userInfo == nil {
        fmt.Println(err)
        fmt.Println(userInfo)
        c.JSON(http.StatusBadRequest, gin.H{
            "code": ErrorCode.InvalidParams,
            "msg":  "获取用户信息失败, 用户名: " + uname,
        })

        return
    }


    
    c.JSON(http.StatusOK, gin.H{
        "code": ErrorCode.Success,
        "msg": "获取用户信息成功",
        "userInfo": userInfo,
    })
}
