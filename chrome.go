package urlrich

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chromedp/chromedp"
)

type ChromeDebugJson struct {
	// Browser              string `json:"Browser"`
	// Protocol_Version     string `json:"Protocol-Version"`
	// User_Agent           string `json:"User-Agent"`
	// V8_Version           string `json:"V8-Version"`
	// WebKit_Version       string `json:"WebKit-Version"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

// initChromedpCtx 初始化 chromedp context & 初始化 chromedp
func (o *UrlRich) initLocalChromedpCtx() error {
	allocatorOptions := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
	)

	if o.debug {
		allocatorOptions = append(
			allocatorOptions,
			chromedp.Flag("headless", true),
		)
	}

	*o.allocatorCtx, *o.allocatorCancel = chromedp.NewExecAllocator(context.Background(), allocatorOptions...)

	// 打开一个chrome tab，并维持，避免重复创建浏览器
	chromeCtx, chromeCancel := o.NewChromeTab()
	// init chromedp run first once, first run will be open a new browser, so dont, cancel it
	err := chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)
	if err != nil {
		chromeCancel()
		return err
	}

	o.chromeInit = true
	// o.chromedpCtx = &chromeCtx
	// o.ChromedpCancel = &cancel
	return nil
}

// connectRemoteChrome 初始化远程chrome allocator 并返回ctx和 cancle，如果远程的 wsurl 发生变化，则 cancel 掉，重新 init
func (o *UrlRich) connectRemoteChrome() {
	*o.allocatorCtx, *o.allocatorCancel = chromedp.NewRemoteAllocator(context.Background(), o.remoteChromeWS)

	o.chromeInit = true
}

// NewChromeTab 调用 chromedp.NewContext，相当于打开一个新的 chrome tab
func (o *UrlRich) NewChromeTab() (context.Context, context.CancelFunc) {
	if o.useRemoteChrome {
		chromeCtx, cancel := chromedp.NewContext(
			*o.allocatorCtx,
		)
		return chromeCtx, cancel
	}

	contextOption := []chromedp.ContextOption{
		chromedp.WithLogf(log.Printf),
	}
	if o.debug {
		contextOption = append(contextOption, chromedp.WithDebugf(log.Printf))
	}

	chromeCtx, cancel := chromedp.NewContext(
		*o.allocatorCtx,
		contextOption...,
	)

	return chromeCtx, cancel
}

// UpdateRemoteChromeDebugURL 从chrome的json地址，获取 ws debug 地址
// 当连接chrome失败的时候，可以调用这个重新获取一下 debug url 地址
func (o *UrlRich) UpdateRemoteChromeDebugURL(remoteChromeHTTP string) error {

	o.remoteChromeHTTP = remoteChromeHTTP
	resp, err := http.Get(o.remoteChromeHTTP + "/json/version")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	jsonStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// fmt.Println(string(jsonStr))

	var chromeDebugJson ChromeDebugJson
	err = json.Unmarshal(jsonStr, &chromeDebugJson)
	if err != nil {
		return err
	}

	if len(chromeDebugJson.WebSocketDebuggerURL) > 0 {
		o.remoteChromeWS = chromeDebugJson.WebSocketDebuggerURL
		// o.connectRemoteChrome() // 重新获取了 debug url 之后，需要重新 NewRemoteAllocator，所以需要再调用一次 init
	}

	// log.Println(string(jsonStr))
	return nil
}
