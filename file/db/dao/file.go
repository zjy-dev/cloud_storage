package dao

import (
    "github.com/J-Y-Zhang/cloud-storage/file/db/model"
    "log"
)

//init 建表
func init() {
    err := GetMysqlConn().AutoMigrate(&model.File{})
    if err != nil {
        log.Printf("建立File表失败, 错误信息 %v\n", err)
        return
    }
    log.Println("建立File表成功")
    err = GetMysqlConn().AutoMigrate(&model.UserFile{})
    if err != nil {
        log.Printf("建立UserFile表失败, 错误信息 %v\n", err)
        return
    }
    log.Println("建立UserFile表成功")
}

//AddFile 向mysql中添加文件
func AddFile(file *model.File) error {
    return GetMysqlConn().Create(&file).Error
}

//DeleteFileByMd5 通过Md5删除指定文件
func DeleteFileByMd5(md5 string) error {
    return GetMysqlConn().Where("file_md5 = ?", md5).Delete(&model.File{}).Error
}

//UpdateFileByMd5 更新文件Path
func UpdateFileByMd5(md5, path string) error {
    return GetMysqlConn().Model(&model.File{}).Where("file_md5 = ?", md5).Update("file_path", path).Error
}

//FindFileByMd5 通过Md5查找文件
func FindFileByMd5(md5 string) (*model.File, error) {
    file := &model.File{}
    return file, GetMysqlConn().Where("file_md5 = ?", md5).First(file).Error
}

