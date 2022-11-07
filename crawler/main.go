package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
)

const groupName = "国星宇航"

func main() {
	bot := openwechat.DefaultBot(openwechat.Desktop)
	bot.UUIDCallback = func(uuid string) {
		q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
		fmt.Println(q.ToString(true))
	}
	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	}
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")
	err := bot.HotLogin(reloadStorage, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	memberMap := make(map[string]string)

	groups, _ := self.Groups()
	for _, group := range groups.SearchByNickName(1, groupName) {
		fmt.Println(group)
		members, _ := group.Members()
		for _, member := range members {
			if _, ok := memberMap[member.NickName]; ok {
				memberMap[member.NickName+member.DisplayName] = member.DisplayName
			} else {
				memberMap[member.NickName] = member.DisplayName
			}
		}
	}
	year, month, day := time.Now().Date()
	filename := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
	fmt.Println(len(memberMap), filename)
	b, err := json.MarshalIndent(memberMap, "", "  ")
	if err != nil {
		fmt.Println("json err: ", err)
		return
	}
	file, err := os.Create(filename + ".json")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	file.Write(b)
	file.Close()
	bot.Block()
}

// 判断所给路径文件/文件夹是否存在

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
