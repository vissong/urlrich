package urlrich

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/vissong/go-readability"
)

type UrlRich struct {
	chromeInit       bool
	localChromeRuned bool   // 本地的chrome是否已经执行（第一次RUN是否已经调用过）
	useChrome        bool   // 使用 chrome 获取数据，初始化了 chrome 之后，自动设置
	useRemoteChrome  bool   // 使用远程 chrome，初始化了远程 chrome 之后，自动设置
	remoteChromeHTTP string // 远程 chrome 的 ws 地址，比如 http://127.0.0.1:9222
	remoteChromeWS   string // 请求 http://127.0.0.1:9222/json 后，返回中有一个完整的 url，比如 ws://127.0.0.1:9222/devtools/page/6AAF75357FA5B76E36E50C2C7B3FC284
	httpProxy        string // http 请求代理

	userAgent   string        // ua 不配置则使用默认UA
	debug       bool          // 是否 debug，debug 进对本地 chrome 起作用，会开启 chrome 的debug日志
	timeout     time.Duration // 通过http get 或者 chrome 抓取网页数据的超时时间
	downgrading bool          // 当使用 chrome 超时的时候，是否退化到 http 请求，默认为 false

	allocatorCtx    context.Context
	allocatorCancel context.CancelFunc

	// 避免内网抓取配置
	blackDomains []string // 黑名单域名，在黑名单中的域名不去抓取
	blackIps     []string // 黑名单ip，在黑名单中的ip不去抓取
}

type RichResult struct {
	Url        string `json:"url"`
	Readable   bool   `json:"readable"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	ImageUrl   string `json:"imageurl"`
	FaviconUrl string `json:"faviconurl"`

	Downgrading bool `json:"downgrading"` // 是否是基于chrome降级到 http 拉取的页面结果

	IsBlack bool `json:"isblack"` // URL是否在黑名单，禁止抓取
}

var logger Logger

var ERR_BLACK_URL = errors.New("black url")

// Init 初始化
// 1. 初始化超时等配置
// 2. 初始化 chromedp，在服务进程中，只需要调用一次，维持一个 chromedp 的实例
// 3. Init 方法报错的时候，如果是连接远程chrome报错，需要重新取 chrome ip 进行重连
func New(opt ...Option) (*UrlRich, error) {

	logger = *NewLogger("/tmp/", "urlrich.log", 255)

	o := new(UrlRich)
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
			err := o.updateRemoteChromeWS(o.remoteChromeHTTP)
			if err != nil {
				return nil, err
			}
			err = o.connectRemoteChrome()
			if err != nil {
				return nil, err
			}
		} else {
			err := o.initLocalChromedpCtx()
			if err != nil {
				return nil, err
			}
		}
	}

	return o, nil
}

// Do 执行请求
func (o *UrlRich) Do(url string) (RichResult, error) {

	checkResult, errCheck := o.IsBlack(url)
	// 黑名单检测的错误，不阻断逻辑，做柔性，只打日志
	if errCheck != nil {
		log.Println(errCheck)
	}
	if checkResult {
		return RichResult{Url: url, IsBlack: checkResult}, nil
	}

	var html string
	var err error
	var hadDowngrading bool
	if o.useChrome {
		if o.chromeInit != true {
			log.Println("UrlRich is not inited, please check")
			return RichResult{}, errors.New("UrlRich is not inited, please check")
		}
		// html, err = requestByChromedp2(o.chromedpCtx, o.url, o.timeout)
		html, err, hadDowngrading = o.requestByChromedp(url)
	} else {
		html, err = o.requestByHTTP(url)
	}

	if err != nil {
		// 如果是在 http 请求跳转过程中命中的黑名单，会返回这个错误
		if err == ERR_BLACK_URL {
			return RichResult{Url: url, IsBlack: true}, nil
		}
		return RichResult{}, err
	}

	result := RichResult{}
	result.Readable = readability.IsReadable(strings.NewReader(html))

	article, err := readability.FromReader(strings.NewReader(html), url)
	if err != nil {
		log.Printf("failed to parse %s: %v\n", url, err)
		return RichResult{}, err
	}

	result.Url = url
	result.Title = article.Title
	result.Desc = article.Excerpt
	result.ImageUrl = article.Image
	result.FaviconUrl = article.Favicon

	if hadDowngrading {
		result.Downgrading = true
	}

	return result, nil
}

// IsUseRemoteChrome 实例是否设置了使用远程 chrome
func (o *UrlRich) IsUseRemoteChrome() bool {
	return o.useRemoteChrome
}
