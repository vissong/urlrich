package urlrich

import (
	"context"
	"encoding/json"
	"errors"
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
			chromedp.Flag("headless", false),
		)
	}

	o.allocatorCtx, _ = chromedp.NewExecAllocator(context.Background(), allocatorOptions...)

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
func (o *UrlRich) connectRemoteChrome() error {

	if len(o.remoteChromeWS) == 0 {
		return errors.New("remoteChromeWS is empty")
	}
	o.allocatorCtx, o.allocatorCancel = chromedp.NewRemoteAllocator(context.Background(), o.remoteChromeWS)

	o.chromeInit = true

	return nil
}

// NewChromeTab 调用 chromedp.NewContext，相当于打开一个新的 chrome tab
func (o *UrlRich) NewChromeTab() (context.Context, context.CancelFunc) {
	if o.useRemoteChrome {
		chromeCtx, cancel := chromedp.NewContext(
			o.allocatorCtx,
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
		o.allocatorCtx,
		contextOption...,
	)

	return chromeCtx, cancel
}

// UpdateRemoteChromeDebugURL 从chrome的json地址，获取 ws debug 地址
// 当连接chrome失败的时候，可以调用这个重新获取一下 debug url 地址
func (o *UrlRich) updateRemoteChromeWS(remoteChromeHTTP string) error {

	o.remoteChromeHTTP = remoteChromeHTTP
	resp, err := http.Get(o.remoteChromeHTTP + "/json/version")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	jsonStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Debug("json version return: %s", jsonStr)
		logger.Error("json version err: %v", err)
		// 获取 wsurl 失败，比如被代理拦截等，请求不可达等
		return err
	}
	// fmt.Println(string(jsonStr))

	var chromeDebugJson ChromeDebugJson
	err = json.Unmarshal(jsonStr, &chromeDebugJson)
	if err != nil {
		logger.Error("updateRemoteChromeWS err: %v", err)
		return err
	}

	logger.Debug("chromeDebugJson %v", chromeDebugJson)

	if len(chromeDebugJson.WebSocketDebuggerURL) > 0 {
		o.remoteChromeWS = chromeDebugJson.WebSocketDebuggerURL
	}

	return nil
}

// ReConnectRemote 当使用远程 chrome 的时候，如果上一个 chrome 报错了
// （比如运行chrome的机器下线了）需要先从名字服务中取出ip，然后再重新连接一下远程 chrome
func (o *UrlRich) ReConnectRemote(remoteChromeHTTP string) {
	if o.useRemoteChrome != true {
		return
	}

	if o.chromeInit && o.useRemoteChrome && o.allocatorCancel != nil {
		// chromedp.Cancel(o.allocatorCtx)
		o.allocatorCancel()
	}
	o.remoteChromeHTTP = remoteChromeHTTP
	o.updateRemoteChromeWS(remoteChromeHTTP)
	o.connectRemoteChrome()
}
