package main

import (
	"encoding/json"
	"fmt"

	"github.com/vissong/urlrich"
)

func main() {
	urlrichIns := &urlrich.UrlRich{}
	urlrichIns.Init(urlrich.WithUseChrome(true))

	a, _ := urlrichIns.Do("https:www.baidu.com")

	a, _ = urlrichIns.Do("https://xue.youdao.com/sw/m/1946563?keyfrom=dict2.index")

	jsonstr, _ := json.Marshal(a)
	fmt.Println(string(jsonstr))
}
