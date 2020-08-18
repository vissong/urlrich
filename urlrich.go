package urlrich

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/go-shiori/go-readability"
)

type UrlRich struct {
	url            string
	userAgent      string
	debug          bool
	timeout        time.Duration
	useChrome      bool
	chromeInit     bool
	chromedpCtx    context.Context
	chromedpCancel context.CancelFunc
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
func (o *UrlRich) Init(opt ...Option) *UrlRich {
	for _, op := range opt {
		op(o)
	}

	if o.useChrome && !o.chromeInit {
		o.initChromedpCtx()
	}

	if o.timeout == 0 {
		o.timeout = 10 * time.Second
	}

	return o
}

// Do 执行请求
func (o *UrlRich) Do(url string) (RichResult, error) {

	o.url = url

	var html string
	var err error
	if o.useChrome {
		html, err = o.requestByChromedp()
	} else {
		html, err = o.requestByHTTP()
	}

	if err != nil {
		return RichResult{}, err
	}

	result := RichResult{}
	result.Readable = readability.IsReadable(strings.NewReader(html))

	article, err := readability.FromReader(strings.NewReader(html), o.url)
	if err != nil {
		log.Fatalf("failed to parse %s: %v\n", o.url, err)
	}

	result.Url = o.url
	result.Title = article.Title
	result.Desc = article.Excerpt
	result.ImageUrl = article.Image
	result.FaviconUrl = article.Favicon

	return result, nil
}

// initChromedpCtx 初始化 chromedp context & 初始化 chromedp
func (o *UrlRich) initChromedpCtx() {
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
	o.chromedpCancel = cancel
}
