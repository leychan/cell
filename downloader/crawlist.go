package downloader

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

var list = []string{}
var index = []string{}
var index_url = "https://pinyin.sogou.com/dict/"
var regex = regexp.MustCompile(`/dict/cate/index/[\d]+\?rf=dictindex`)
var regex_cell = regexp.MustCompile(`/http:\/\/download\.pinyin\.sogou.com\/dict\/download_cell.php\?id\=[\d]+&name\=(\*)`)

func CrawIndex() {
	res, err := http.Get(index_url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	param := regex.FindAll(data, -1)
	for _, p := range param {
		index = append(index, string(p))
	}

}

func crawDicList() {
	index_len := len(index)
	wg := &sync.WaitGroup{}
	wg.Add(index_len)
	for i, list_entry := range index {
		go func(i int, list_entry string) {

		}(i, list_entry)
	}
}
