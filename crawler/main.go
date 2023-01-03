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
	router.GET("/ask-for-login", askForLogin)

	done := make(chan bool)
	router.GET("/secret-kill-server", func(ctx *gin.Context) {
		if bot != nil {
			self, err := bot.GetCurrentUser()
			if err != nil {
				ctx.String(200, "Kill error GetCurrentUser")
				bot = nil
			}
			fh, err := self.FileHelper()
			if err != nil {
				ctx.String(200, "Kill error FileHelper")
				bot = nil
			}
			self.SendTextToFriend(fh, "server killed by request")
			bot.Logout()
		}
		done <- true
	})

	go func() {
		router.Run()
	}()
	<-done
}

// 判断所给路径文件/文件夹是否存在

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
