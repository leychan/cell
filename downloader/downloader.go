package downloader

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var List = []string{}

var path = "/home/chenlei/tmp/sogou_cell/"

func Download() {
	wg := sync.WaitGroup{}
	list_len := len(List)
	wg.Add(list_len)
	for i, l := range List {
		go func(l string, i int) {
			res, _ := http.Get(l)
			defer res.Body.Close()
			data, _ := ioutil.ReadAll(res.Body)
			dirExits()
			err := ioutil.WriteFile(path + strconv.Itoa(i)+".scel", data, 0777)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(l, i)
	}
	wg.Wait()
}

func dirExits() {
	if _, err := os.Stat(path); err != nil {
		_ = os.Mkdir(path, 0777)
	}
}
