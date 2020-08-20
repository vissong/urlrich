package urlrich

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/vissong/go-readability"
)

type UrlRich struct {
	userAgent        string
	debug            bool
	timeout          time.Duration
	useChrome        bool
	useRemoteChrome  bool
	remoteChromeHTTP string // 远程 chrome 的 ws 地址，比如 http://127.0.0.1:9222
	remoteChromeWS   string // 请求 http://127.0.0.1:9222/json 后，返回中有一个完整的 url，比如 ws://127.0.0.1:9222/devtools/page/6AAF75357FA5B76E36E50C2C7B3FC284
	chromeInit       bool
	allocatorCtx     *context.Context
	allocatorCancel  *context.CancelFunc
}

type RichResult struct {
	Url        string `json:"Url"`
	Readable   bool   `json:"Readable"`
	Title      string `json:"Title"`
	Desc       string `json:"Desc"`
	ImageUrl   string `json:"ImageUrl"`
	FaviconUrl string `json:"FaviconUrl"`
}

// Init 初始化
// 1. 初始化超时等配置
// 2. 初始化 chromedp，在服务进程中，只需要调用一次，维持一个 chromedp 的实例
func (o *UrlRich) Init(opt ...Option) {
	for _, op := range opt {
		op(o)
	}

	if o.timeout == 0 {
		o.timeout = 10 * time.Second
	}

	if len(o.userAgent) == 0 {
		o.userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"
	}

	if o.useChrome && !o.chromeInit {
		if o.useRemoteChrome {
			o.UpdateRemoteChromeDebugURL(o.remoteChromeHTTP)
			o.connectRemoteChrome()
		} else {
			o.initLocalChromedpCtx()
		}
	}
}

func (o *UrlRich) ReConnectRemote(remoteChromeHTTP string) {
	chromedp.Cancel(*o.allocatorCtx)
	o.remoteChromeHTTP = remoteChromeHTTP
	o.UpdateRemoteChromeDebugURL(remoteChromeHTTP)
	o.connectRemoteChrome()
}

// Do 执行请求
func (o *UrlRich) Do(url string) (RichResult, error) {

	var html string
	var err error
	if o.useChrome {
		// html, err = requestByChromedp2(o.chromedpCtx, o.url, o.timeout)
		html, err = o.requestByChromedp(url)
	} else {
		html, err = o.requestByHTTP(url)
	}

	if err != nil {
		return RichResult{}, err
	}

	result := RichResult{}
	result.Readable = readability.IsReadable(strings.NewReader(html))

	article, err := readability.FromReader(strings.NewReader(html), url)
	if err != nil {
		log.Fatalf("failed to parse %s: %v\n", url, err)
		return RichResult{}, err
	}

	result.Url = url
	result.Title = article.Title
	result.Desc = article.Excerpt
	result.ImageUrl = article.Image
	result.FaviconUrl = article.Favicon

	return result, nil
}
