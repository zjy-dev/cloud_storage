package model

import "gorm.io/gorm"


type UserFile struct {
    gorm.Model
    UserName string             `gorm:"not_null"`
    FileID string
    FileName string
    FileSize int
    Files []File                `gorm:"foreignKey:FileMd5;references:FileID"`
}
