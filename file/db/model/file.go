package model

import (
    "gorm.io/gorm"
)

type
File struct {
    gorm.Model
    FileName string
    FileMd5 string          `gorm:"unique"`
    FilePath string
    FileSize int
}


