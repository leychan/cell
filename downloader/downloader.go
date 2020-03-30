package downloader

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var path = "/home/lei/tmp/sogou_cell/"

//根据传入的url,返回网页主体内容
func GetBody(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	return data
}

func Save(data []byte, name string) {
	dirExitsOrCreate()

	err := ioutil.WriteFile(path+name+".scel", data, 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func dirExitsOrCreate() bool {
	if _, err := os.Stat(path); err != nil {
		_ = os.Mkdir(path, 0777)
		return false
	}
	return true
}

func FileExists() map[string]int {
	var fileExist = map[string]int{}
	if dirExitsOrCreate() {
		d, _ := ioutil.ReadDir(path)
		for _, f := range d {
			fileExist[f.Name()] = 1
		}
	}
	return fileExist
}
