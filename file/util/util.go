package util

import (
    "crypto/md5"
    "crypto/sha1"
    "encoding/hex"
    "fmt"
    "io"
    "math"
    "os"
)

const (
    KB float64 = 1024
    MB float64 = 1024 * 1024
    GB float64 = 1024 * 1024 * 1024
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

func CompareInt(first int, args ...int) bool {
    for _, val := range args {
        if val != first {
            return false
        }
    }
    return true
}

func FomatFloat64(val float64, bit int) string {
    t := math.Pow10(bit)
    return fmt.Sprintf("%v", float64(int(val * t)) / t)
}