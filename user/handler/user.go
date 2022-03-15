package handler

import (
    "context"
    "user/domain/model"
    "user/domain/service"
    pb "user/proto"
)

type User struct{
    userService service.UserService
}

//SignUp 用户注册服务handler
func (u User) SignUp(ctx context.Context, req *pb.UserSignUpReq, resp *pb.UserSignUpResp) error {
    _, err := u.userService.AddUserService(&model.User{
        Name:        req.GetUserName(),
        Encoded_Pwd: req.GetUserPwd(),
        NickName:    req.GetNickName(),
    })

    resp.IsUserExist = (err != nil)

    return err
}

//SignIn 用户登录handler
func (u User) SignIn(ctx context.Context, req *pb.UserSignInReq, resp *pb.UserSignInResp) error {
    // 1.检查用户是否存在
    user, err := u.userService.FindUserByNameService(req.GetUserName())

    if user == nil || err != nil {
        resp.IsUserExist = false
        return err
    }
    resp.IsUserExist = true

    // 2.用户存在, 则校验密码
    ok, err := u.userService.CheckPwdService(req.GetUserName(), req.GetUserPwd())
    resp.IsPwdError = (!ok || err != nil)
    return err
}

//GetUserInfo 获取用户信息handler
func (u User) GetUserInfo(ctx context.Context, req *pb.UserInfoReq, resp *pb.UserInfoResp) error {
    user, err := u.userService.FindUserByNameService(req.GetUserName())
    resp.UserId = int64(user.ID)
    resp.UserName = user.Name
    resp.NickName = user.NickName

    return err
}


