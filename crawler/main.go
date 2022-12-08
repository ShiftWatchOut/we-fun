package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

const groupName = "国星宇航"

func main() {
	// loadJSON()

	router := gin.Default()

	router.GET("/logout-wechat", logoutWechat)
	router.GET("/save-group-info", saveGroupInfo)
	router.GET("/sql-test", sqlTest)

	router.Run()
}

// 判断所给路径文件/文件夹是否存在

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
