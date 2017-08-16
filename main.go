package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var strTable = map[string]string{
	"_z2C$q": ":",
	"_z&e3B": ".",
	"AzdH3F": "/",
}

var charTable = map[string]string{
	"w": "a",
	"k": "b",
	"v": "c",
	"1": "d",
	"j": "e",
	"u": "f",
	"2": "g",
	"i": "h",
	"t": "i",
	"3": "j",
	"h": "k",
	"s": "l",
	"4": "m",
	"g": "n",
	"5": "o",
	"r": "p",
	"q": "q",
	"6": "r",
	"f": "s",
	"p": "t",
	"7": "u",
	"e": "v",
	"o": "w",
	"8": "1",
	"d": "2",
	"n": "3",
	"9": "4",
	"c": "5",
	"m": "6",
	"0": "7",
	"b": "8",
	"l": "9",
	"a": "0",
}

var cnum chan int
var failedCount = 0
var lock = &sync.Mutex{}

func main() {
	fmt.Println(fozu())
	doit()
	// for {

	// }
}

func doit() {
	fmt.Println("请输入关键字,回车后开始下载图片")
	reader := bufio.NewReader(os.Stdin)
	var word string
	for {
		b, _, _ := reader.ReadLine()
		word = string(b)
		word = strings.Replace(word, "\n", "", -1)
		if len(strings.TrimSpace(word)) != 0 {
			break
		}
	}
	fmt.Printf("输入的数据：%s\n", word)
	var baseURL = "http://image.baidu.com/search/acjson?tn=resultjson_com&ipn=rj&ct=201326592&fp=result&cl=2&lm=-1&ie=utf-8&oe=utf-8&st=-1&ic=0&face=0&istype=2nc=1&rn=60&word=" + word + "&pn="
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = strings.Replace(dir, "\\", "/", -1)
	dir = dir + "/results/" + word + "/"
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0777)
	var total = 0
	cnum = make(chan int, 60)
	for i := 0; ; i++ {
		rn := 60
		pn := rn*i + 1
		reqURL := baseURL + strconv.Itoa(pn)
		resp, err := http.Get(reqURL)
		if err != nil {
			fmt.Println("error ")
			fmt.Println(err.Error())
			return
		}
		defer resp.Body.Close()
		content, _ := ioutil.ReadAll(resp.Body)
		reg := regexp.MustCompile("\"objURL\":\"(.*?)\"")
		var imgList = reg.FindAllString(string(content), -1)

		if len(imgList) == 0 {
			fmt.Println("no data ")
			fmt.Println("图片数量：", total)
			break
		}
		for _, img := range imgList {
			total++
			img = img[10 : len(img)-1]
			img = buildURL(img)
			fileName := dir + strconv.Itoa(total) + ".jpg"
			fmt.Printf("开始获取第%d张图片，url:%s\n", total, img)
			go download(fileName, img)
		}
		fmt.Println("是否继续获取图片（Y/N）")
		cb, _, _ := reader.ReadLine()
		isContinue := string(cb)
		isContinue = strings.Replace(isContinue, "\n", "", -1)
		if !strings.EqualFold(isContinue, "Y") {
			break
		}
	}
	<-cnum
	fmt.Printf("total :%d, failed :%d\n", total, failedCount)

	fmt.Println("再来一次（Y/N）")
	cb, _, _ := reader.ReadLine()
	isContinue := string(cb)
	isContinue = strings.Replace(isContinue, "\n", "", -1)
	if strings.EqualFold(isContinue, "Y") {
		doit()
	}
}

// orginalURL 构建访问路径
func buildURL(orginalURL string) string {
	for k, v := range strTable {
		orginalURL = strings.Replace(orginalURL, k, v, -1)
	}
	var newURL = ""
	for _, v := range orginalURL {
		v1 := string(v)
		v2, succ := charTable[v1]
		if succ {
			v1 = v2
		}
		newURL = newURL + v1
	}

	return newURL
}

func download(filePath, url string) {
	r, err := http.Get(url)
	defer func() {
		if r != nil && r.Body != nil && !r.Close {
			r.Body.Close()
		}
		cnum <- 1
	}()

	if err != nil || r.StatusCode != 200 {
		fmt.Println("下载失败啦：" + url)
		lock.Lock()
		failedCount++
		lock.Unlock()
		return
	}
	imagebyte, _ := ioutil.ReadAll(r.Body)
	ioutil.WriteFile(filePath, imagebyte, 0777)
}

func fozu() string {

	buff := bytes.NewBufferString("")
	buff.WriteString("                   _ooOoo_ \n")
	buff.WriteString("                  o8888888o \n")
	buff.WriteString("                  88\" . \"88 \n")
	buff.WriteString("                  (| -_- |) \n")
	buff.WriteString("                  O\\  =  /O \n")
	buff.WriteString("               ____/`---'\\____ \n")
	buff.WriteString("             .'  \\|     |//  `. \n")
	buff.WriteString("            /  \\|||  :  |||//  \\ \n")
	buff.WriteString("           /  _||||| -:- |||||-  \\ \n")
	buff.WriteString("           |   | \\\\  -  /// |   | \n")
	buff.WriteString("           | \\_|  ''\\---/''  |   | \n")
	buff.WriteString("           \\  .-\\__  `-`  ___/-. / \n")
	buff.WriteString("         ___`. .'  /--.--\\  `. . __ \n")
	buff.WriteString("      .\"\" '<  `.___\\_<|>_/___.'  >'\"\". \n")
	buff.WriteString("     | | :  `- \\`.;`\\ _ /`;.`/ - ` : | | \n")
	buff.WriteString("     \\  ]\\ `-.   \\_ __\\ /__ _/   .-` /  / \n")
	buff.WriteString("======`-.____`-.___\\_____/___.-`____.-'====== \n")
	buff.WriteString("                  `=---=' \n")
	buff.WriteString("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ \n")
	buff.WriteString("         佛祖保佑       永无BUG \n")
	return buff.String()
}
