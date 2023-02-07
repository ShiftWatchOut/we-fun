package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/skip2/go-qrcode"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var bot *openwechat.Bot

func logoutWechat(c *gin.Context) {
	if bot != nil {
		err := bot.Logout()
		if err != nil {
			fmt.Printf("%s", err)
		}
		bot = nil
	}
	c.String(http.StatusOK, "Logged out")
}

func askForLogin(c *gin.Context) {
	bot = openwechat.DefaultBot(openwechat.Desktop)
	done := make(chan bool)
	bot.UUIDCallback = func(uuid string) {
		q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
		imgPath := "./tmp.jpg"
		q.WriteFile(200, imgPath)
		c.File(imgPath)
		done <- true
	}
	go func() {
		reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")
		err := bot.HotLogin(reloadStorage, true)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	<-done
}

func saveGroupInfo(c *gin.Context) {
	bot = openwechat.DefaultBot(openwechat.Desktop)
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
	c.String(http.StatusOK, "Logged in well")
}

func sqlTest(c *gin.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}
	sqlStr := os.Getenv("SQL_STR")
	fmt.Println("sql connect: ", sqlStr)
	db, err := sql.Open("mysql", sqlStr)
	if err != nil {
		log.Fatal(err)
		return
	}
	db.Ping()
	defer db.Close()
	_, err = db.Exec("CREATE TABLE " + "(id INT NOT NULL , name VARCHAR(20), PRIMARY KEY(ID));")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Successfully Created")
}

func loadJSON() {
	// 数据库
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}
	dsn := os.Getenv("SQL_STR")
	fmt.Println("sql connect: ", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}

	// 读文件
	pwd, _ := os.Getwd()
	fileInfoList, _ := os.ReadDir(pwd)
	for f := range fileInfoList {
		file := fileInfoList[f]
		filename := file.Name()
		if strings.Contains(filename, ".json") && strings.Contains(filename, "-") {
			content, err := os.ReadFile(filename)
			if err != nil {
				panic(err)
			}
			memberMap := make(map[string]string)
			err = json.Unmarshal(content, &memberMap)
			if err != nil {
				panic(err)
			}
			tableName := strings.Split(filename, ".")[0]
			db.Table(tableName).AutoMigrate(&GroupMember{})
			memberList := make([]GroupMember, 0)
			for k, v := range memberMap {
				memberList = append(memberList, GroupMember{
					Nickname:    k,
					IngroupName: v,
				})
			}
			db.Table(tableName).Create(&memberList)
		}
	}
}

type GroupMember struct {
	Nickname    string
	IngroupName string
}
