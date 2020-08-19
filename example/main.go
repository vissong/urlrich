package main

import (
	"encoding/json"
	"fmt"

	"github.com/vissong/urlrich"
)

func main() {
	urlrichIns := &urlrich.UrlRich{}
	urlrichIns.Init(
		urlrich.WithUseChrome(true),
		urlrich.WithRemoteChrome("http://127.0.0.1:9222"),
		// urlrich.WithDebug(true),
	)

	urls := []string{
		"https://www.baidu.com",
		"https://www.qq.com",
		"https://xue.youdao.com/sw/m/1946563?keyfrom=dict2.index",
	}

	for _, url := range urls {
		a, _ := urlrichIns.Do(url)
		jsonstr, _ := json.Marshal(a)
		fmt.Println(string(jsonstr))
	}

	// a, _ = urlrichIns.Do("https://xue.youdao.com/sw/m/1946563?keyfrom=dict2.index")

	// jsonstr, _ = json.Marshal(a)
	// fmt.Println(string(jsonstr))
}
