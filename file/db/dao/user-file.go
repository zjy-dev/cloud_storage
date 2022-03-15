package dao

import (
	"github.com/J-Y-Zhang/cloud-storage/file/db/model"
)



//UserAddFile 用户添加文件
func UserAddFile(userFile *model.UserFile) error {
	return GetMysqlConn().Create(userFile).Error
}

//UserDeleteFileByMd5 用户删除文件
func UserDeleteFileByMd5(uname, md5 string) error {
    return GetMysqlConn().Where("user_name = ? and file_id = ?", uname, md5).Delete(&model.UserFile{}).Error
}

//UserUpdateFileName 用户更新文件
func UserUpdateFileName(uname, md5, newname string) error {
	return GetMysqlConn().Model(&model.UserFile{}).Where("user_name = ? and file_id = ?", uname, md5).UpdateColumn("file_name", newname).Error
}

//UserQueryFiles 用户查询所有文件
func UserQueryFiles(uname string) ([]*model.UserFile, error) {
	userFiles := make([]*model.UserFile, 0)
	return userFiles, GetMysqlConn().Where("user_name = ?", uname).Find(&userFiles).Error
}

//UserQueryFileByMd5 用户查询指定文件
func UserQueryFileByMd5(uname, md5 string) (*model.UserFile, error) {
	userFile := &model.UserFile{}
	return userFile, GetMysqlConn().Where("user_name = ? and file_id = ?", uname, md5).First(userFile).Error
}




