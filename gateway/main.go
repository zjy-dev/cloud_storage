package main

import (
    "github.com/J-Y-Zhang/cloud-storage/gateway/router"
)


func main() {
    r := router.Router()

    r.Run("0.0.0.0:8080")
}