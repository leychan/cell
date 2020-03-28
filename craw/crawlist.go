package craw

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/leychan/cell/downloader"
)

type Dic struct {
	dicUrl []byte
	id     []byte
	name   []byte
}

type Index struct {
	indexUrl []byte
	id       []byte
}

var index = []Index{}
var indexUrl = "https://pinyin.sogou.com/dict/"
var regex = regexp.MustCompile(`/dict/cate/index/([\d]+)\?rf=dictindex`)
var regexCell = regexp.MustCompile(`http:\/\/download\.pinyin\.sogou\.com\/dict\/download_cell\.php\?id=([\d]+)&name=(.*)\"`)
var regexPageCount = regexp.MustCompile(`分类下共有([\d]+)个词库`)
var ch = make(chan Dic, 10)
var uniqueUrl = map[string]string{}
var wgGlobal = &sync.WaitGroup{}
var fileExist = downloader.FileExists()

//启动
func CrawStart() {
	data := downloader.GetBody(indexUrl)
	param := regex.FindAllSubmatch(data, -1)
	for _, p := range param {
		indexTmp := Index{
			indexUrl: p[0],
			id:       p[1],
		}
		index = append(index, indexTmp)
	}
	wgGlobal.Add(2)
	go func() {
		crawDicList()
		close(ch)
		wgGlobal.Done()
	}()
	go func() {
		for i := range ch {
			name, _ := url.QueryUnescape(string(i.name))
			//去重判断
			if _, ok := uniqueUrl[name]; ok {
				continue
			}
			dicUrl := string(i.dicUrl)
			uniqueUrl[name] = dicUrl
			name_decode, _ := url.QueryUnescape(name)
			fmt.Println(name_decode)
			if _, ok := fileExist[name_decode]; !ok {
				data := downloader.GetBody(dicUrl)
				downloader.Save(data, formatFileName(name_decode))
				time.Sleep(time.Millisecond * 300)
			}

		}
		wgGlobal.Done()
	}()

	wgGlobal.Wait()
}

//爬取
func crawDicList() {
	index_len := len(index)
	wg := &sync.WaitGroup{}
	wg.Add(index_len)
	for _, entry := range index {
		//子分类, 此页面可得到下载链接
		go func(entry Index) {
			crawDicPage(entry)
			wg.Done()
		}(entry)
	}
	wg.Wait()
}

func crawDicPage(entry Index) {
	data := downloader.GetBody("https://pinyin.sogou.com" + string(entry.indexUrl))
	wg := &sync.WaitGroup{}
	dicPage := regexPageCount.FindSubmatch(data)
	if len(dicPage) != 0 {
		dicCount, _ := strconv.Atoi(string(dicPage[1]))
		page := calculatePage(dicCount)
		wg.Add(page)
		for i := 1; i <= page; i++ {
			pageUrl := "https://pinyin.sogou.com/dict/cate/index/" + string(entry.id) + "/default/" + strconv.Itoa(i)
			go func(pageUrl string) {
				crawDicPageUrl(pageUrl)
				wg.Done()
			}(pageUrl)
		}
		wg.Wait()
	}
}

//爬取每一页的字典url
func crawDicPageUrl(url string) {
	data := downloader.GetBody(url)
	dicCell := regexCell.FindAllSubmatch(data, -1)
	if len(dicCell) != 0 {
		for _, d := range dicCell {
			dic := Dic{d[0][0 : len(d[0])-1], d[1], d[2]}
			ch <- dic
		}
	}
}

//根据总条数计算总页数
func calculatePage(p int) int {
	if p%10 > 0 {
		return p/10 + 1
	}
	return p / 10
}

func formatFileName(name string) string {
	match, _ := regexp.MatchString("/", name)
	if match {
		name = strings.Replace(name, "/", "", -1)
	}
	return name
}
