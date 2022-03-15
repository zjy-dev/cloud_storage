package algorithm

import (
    "io/ioutil"
    "os"
    "path"
    "strconv"
)

func MergeFile(filePath, fileName, storePath string, chunkCnt int)  {
    os.MkdirAll(path.Dir(storePath), os.ModePerm)

    file, _ := os.OpenFile(storePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
    defer file.Close()

    for i := 1; i <= chunkCnt; i++ {
        f, _ := os.OpenFile(filePath + "/" + fileName + "-" + strconv.Itoa(i), os.O_RDONLY, os.ModePerm)
        b, _ := ioutil.ReadAll(f)
        file.Write(b)
        f.Close()
    }
}
