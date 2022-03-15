package service

import (
    "user/domain/dao"
    "user/domain/model"
    "user/util"
)

type UserServiceInterface interface {
    AddUserService(user *model.User) (int, error)
    DeleteUserByIDService(id int) error
    UpdateUserPwdService(user *model.User) error
    FindUserByNameService(name string) (*model.User, error)
    CheckPwdService(string, string) (bool, error)
}

type UserService struct {

}

//AddUserService 新增用户服务
func (u UserService) AddUserService(user *model.User) (int, error) {
    // 加密用户密码, 不要把裸密直接存入数据库
    user.Encoded_Pwd = util.Sha1([]byte(user.Encoded_Pwd))
    return dao.AddUser(user)
}

//DeleteUserByIDService 通过ID删除用户服务
func (u UserService) DeleteUserByIDService(id int) error {
    return dao.DeleteUserByID(id)
}

//UpdateUserService 更新用户密码服务
func (u UserService) UpdateUserPwdService(user *model.User) error {
    // 加密用户密码, 不要把裸密直接存入数据库
    user.Encoded_Pwd = util.Sha1([]byte(user.Encoded_Pwd))
    return dao.UpdateUserPwd(user)
}

//FindUserByNameService 通过Name查找用户服务
func (u UserService) FindUserByNameService(name string) (*model.User, error) {
    return dao.FindUserByName(name)
}

//CheckPwdService 校验密码服务
func (u UserService) CheckPwdService(name string, pwd string) (bool, error) {
    user, err := dao.FindUserByName(name)
    if err != nil {
        return false, err
    }

    return util.Sha1([]byte(pwd)) == user.Encoded_Pwd, nil
}

