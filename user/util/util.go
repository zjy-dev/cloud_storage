package util

import (
    "crypto/md5"
    "crypto/sha1"
    "encoding/hex"
    "io"
    "os"
)

func Sha1(data []byte) string {
    hash := sha1.New()
    hash.Write(data)
    return hex.EncodeToString(hash.Sum([]byte("")))
}

func FileMd5(file *os.File) string {
    hash := md5.New()
    io.Copy(hash, file)
    return hex.EncodeToString(hash.Sum(nil))
}