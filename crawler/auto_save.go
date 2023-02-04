//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func main() {
	gin_base_url := os.Getenv("GIN_BASE_URL")
	url_1 := "http://" + gin_base_url + "/save-group-info"
	url_2 := "http://" + gin_base_url + "/secret-kill-server"
	fmt.Println("auto save running...")
	resp, err := http.Get(url_1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url_1, err)
		os.Exit(1)
	}
	fmt.Printf("fetch: complete %s", b)
	_, _ = http.Get(url_2)
	// json 文件读取
	pwd, _ := os.Getwd()
	fileInfoList, _ := ioutil.ReadDir(pwd)
	// 过滤
	jsonFileList := make(MyFileInfo, 0)
	for f := range fileInfoList {
		file := fileInfoList[f]
		filename := file.Name()
		if strings.Contains(filename, ".json") && strings.Contains(filename, "-") {
			jsonFileList = append(jsonFileList, file)
		}
	}
	// 时间先后排序，取最后两个
	sort.Sort(jsonFileList)
	lenJson := len(jsonFileList)
	file1, file2 := jsonFileList[lenJson-2], jsonFileList[lenJson-1]
	fmt.Println("last two file are:", file1.Name(), file2.Name())
	// 打开文件 diff
	exec.Command("code", "-d", file1.Name(), file2.Name()).Run()
}

type MyFileInfo []fs.FileInfo

func (fi MyFileInfo) Len() int {
	return len(fi)
}

func (fi MyFileInfo) Swap(i, j int) {
	fi[i], fi[j] = fi[j], fi[i]
}

func (fi MyFileInfo) Less(i, j int) bool {
	return fi[i].ModTime().Before(fi[j].ModTime())
}
