package urlrich

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chromedp/chromedp"
)

type ChromeDebugJson []struct {
	Description          string `json:"description"`
	DevtoolsFrontendURL  string `json:"devtoolsFrontendUrl"`
	ID                   string `json:"id"`
	Title                string `json:"title"`
	Type                 string `json:"type"`
	URL                  string `json:"url"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

// initChromedpCtx 初始化 chromedp context & 初始化 chromedp
func (o *UrlRich) initLocalChromedpCtx() {
	allocatorOptions := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"),
	)
	contextOption := make([]chromedp.ContextOption, 0)
	contextOption = append(contextOption,
		chromedp.WithLogf(log.Printf),
	)

	if o.debug {
		allocatorOptions = append(
			allocatorOptions,
			chromedp.Flag("headless", true),
		)
		contextOption = append(contextOption, chromedp.WithDebugf(log.Printf))
	}

	c, _ := chromedp.NewExecAllocator(context.Background(), allocatorOptions...)
	chromeCtx, cancel := chromedp.NewContext(
		c,
		contextOption...,
	)

	// init chromedp run first once
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	o.chromeInit = true
	o.chromedpCtx = chromeCtx
	o.ChromedpCancel = cancel
}

func (o *UrlRich) initRemoteChromeCtx() {
	c, _ := chromedp.NewRemoteAllocator(context.Background(), o.remoteChromeWS)
	chromeCtx, cancel := chromedp.NewContext(
		c,
	)

	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	o.chromeInit = true
	o.chromedpCtx = chromeCtx
	o.ChromedpCancel = cancel
}

// UpdateRemoteChromeDebugURL 从chrome的json地址，获取 ws debug 地址
// 当连接chrome失败的时候，可以调用这个重新获取一下 debug url 地址
func (o *UrlRich) UpdateRemoteChromeDebugURL(remoteChromeHTTP string) error {

	o.remoteChromeHTTP = remoteChromeHTTP
	resp, err := http.Get(o.remoteChromeHTTP + "/json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	jsonStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var chromeDebugJson ChromeDebugJson
	err = json.Unmarshal(jsonStr, &chromeDebugJson)
	if err != nil {
		return err
	}

	if len(chromeDebugJson) > 0 && len(chromeDebugJson[0].WebSocketDebuggerURL) > 0 {
		o.remoteChromeWS = chromeDebugJson[0].WebSocketDebuggerURL
		o.initRemoteChromeCtx() // 重新获取了 debug url 之后，需要重新 NewRemoteAllocator，所以需要再调用一次 init
	}

	// log.Println(string(jsonStr))

	return nil
}
