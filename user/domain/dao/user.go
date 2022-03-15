package dao

import (
    log "go-micro.dev/v4/logger"
    "user/domain/model"
)

//init 建表
func init() {
    err := GetMysqlConn().AutoMigrate(&model.User{})
    if err != nil {
        log.Errorf("建表失败, 错误信息 %v", err)
        return
    }
    log.Infof("建立User表成功")
}

//AddUser 向mysql中添加用户
func AddUser(user *model.User) (id int, err error) {
    return int(user.ID), GetMysqlConn().Create(&user).Error
}

//DeleteUserByID 通过ID删除指定用户
func DeleteUserByID(id int) error {
    return GetMysqlConn().Delete(&model.User{}, uint(id)).Error
}

//UpdateUserPwd 更新用户密码
func UpdateUserPwd(user *model.User) error {
    return GetMysqlConn().Where("name = ", user.Name).Update("encoded_pwd", user.Encoded_Pwd).Error
}

//UpdateUserNickName 更新用户昵称
func UpdateUserNickName (user *model.User) error {
    return GetMysqlConn().Where("name = ", user.Name).Update("nick_name", user.NickName).Error
}

//FindUserByName 通过Name查找用户
func FindUserByName(name string) (*model.User, error) {
    user := &model.User{}
    return user, GetMysqlConn().Where("name = ?", name).First(user).Error
}

//FindUserByID 通过ID查找用户
func FindUserByID(id int) (*model.User, error) {
    user := &model.User{}
    return user, GetMysqlConn().First(user, uint(id)).Error
}