package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/vissong/urlrich"
)

func main() {
	urlrichIns, err := urlrich.New(
		// urlrich.WithUseChrome(true),
		// urlrich.WithRemoteChrome("http://9.144.134.236:80"),
		// urlrich.WithDebug(true),
		urlrich.WithTimeout(3*time.Second),
		urlrich.WithDowngrading(true),
		// urlrich.WithHttpProxy("http://9.144.134.236:8080"),
	)
	// 域名黑名单
	urlrichIns.SetBlack([]string{"*.inner.com"}, nil)
	if err != nil {
		log.Fatal(err)
	}

	urls := []string{
		"https://www.inner.com",
		// "https://www.sina.com.cn",
		// "http://news.fznews.com.cn/dsxw/20200825/5f44b9b09fa35.shtml",
		// "https://blog.csdn.net/qq_33285730/article/details/73239263",
		"https://www.baidu.com",
		// "https://new.qq.com/rain/a/20200820A0K3WB00",
		// "https://blog.csdn.net/qq_33285730/article/details/73239263",
		// "https://www.qq.com",
		// // "http://news.163.com/20/0821/12/FKI9Q79G000189FH.html",
		// "https://xue.youdao.com/sw/m/1946563?keyfrom=dict2.index",
		// "https://www.google.com",
		// "https://new.qq.com/rain/a/20200820A0K3WB00",
		// "https://new.qq.com/rain/a/20200819A0U90I00",
		// "https://im.qq.com",
		// "https://qun.qq.com",
		// "https://docs.qq.com",
		// "https://vip.qq.com",
	}

	for _, url := range urls {
		a, err := urlrichIns.Do(url)
		fmt.Println(err)
		jsonstr, _ := json.Marshal(a)
		fmt.Println(string(jsonstr))
		fmt.Println()
	}

	// fmt.Println("reconnect....")

	// urlrichIns.ReConnectRemote("http://127.0.0.1:9222")
	// for _, url := range urls {
	// 	a, _ := urlrichIns.Do(url)
	// 	jsonstr, _ := json.Marshal(a)
	// 	fmt.Println(string(jsonstr))
	// }

	// a, _ = urlrichIns.Do("https://xue.youdao.com/sw/m/1946563?keyfrom=dict2.index")

	// jsonstr, _ = json.Marshal(a)
	// fmt.Println(string(jsonstr))

}
