package urlrich

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	urllib "net/url"
	"regexp"
	"strings"

	"github.com/chromedp/chromedp"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// requestByHTTP 使用http请求网页
func (o *UrlRich) requestByHTTP(url string) (string, error) {

	var proxy *urllib.URL
	var err error
	if o.httpProxy != "" {
		proxy, err = urllib.Parse(o.httpProxy)
		if err != nil {
			return "", err
		}
	}

	transport := &http.Transport{
		Proxy:               http.ProxyURL(proxy),
		MaxIdleConnsPerHost: 50,
	}
	client := http.Client{
		Timeout:   o.timeout,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 跳转黑名单检查
			if isBlack, _ := o.IsBlack(req.URL.String()); isBlack == true {
				return ERR_BLACK_URL
			}
			return nil
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("new http request error %s: %v\n", url, err)
		return "", err
	}
	req.Header.Add("User-Agent", o.userAgent)
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("failed to download %s: %v\n", url, err)
		if strings.Contains(err.Error(), "black url") {
			return "", ERR_BLACK_URL
		}
		return "", err
	}
	defer resp.Body.Close()

	html, _ := ioutil.ReadAll(resp.Body)

	var body string

	// 判断是否是 utf-8 页面
	reg, _ := regexp.Compile("(utf-8|gbk|gb2312)")
	matched := reg.FindAllString(strings.ToLower(resp.Header.Get("content-type")), -1)
	// 如果返回的数据是 gbk 的，则需要进行转码
	if len(matched) > 0 {
		if matched[0] == "gbk" || matched[0] == "gb2312" {
			utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(html) //将gbk再转换为utf-8
			body = string(utf8Data)
		} else {
			body = string(html)
		}
	} else {
		body = string(html)
	}

	return body, nil
}

// requestByChromedp 使用 chromedp 抓取网页数据
// 第三个参数标示本次请求是否是 chrome 降级到 http 拉取
func (o *UrlRich) requestByChromedp(url string) (string, error, bool) {

	logger.Debug("begin chrome request")

	// 先打开一个新tab
	chromeCtx, chromeCancel := o.NewChromeTab()
	defer chromeCancel()

	// 超时控制
	taskCtx, taskCancel := context.WithTimeout(chromeCtx, o.timeout)
	defer taskCancel()

	// 载入页面
	var body string
	// var charset, contentType string
	err := chromedp.Run(taskCtx,
		// network.Enable(),
		chromedp.Navigate(url),
		chromedp.WaitReady("document"),
		// chromedp.WaitVisible(`html`),
		chromedp.OuterHTML(`document.querySelector("html")`, &body, chromedp.ByJSPath),
		// chromedp.OuterHTML(`document.querySelector("meta[charset]")`, &charset, chromedp.ByJSPath),
		// chromedp.OuterHTML(`document.querySelector("meta[http-equiv='Content-Type']")`, &contentType, chromedp.ByJSPath),
	)

	// 超时处理
	// context.DeadlineExceeded 为chrome渲染执行超时
	// 其他错误的情况，可能是连接到远程 chrome 超时
	// if err == context.DeadlineExceeded {
	if err != nil {
		logger.Error("chrome err: %v", err)
		return o.chromeTimeout(url)
	}

	return body, nil, false
}

// chromeTimeout timeout 处理，如果有降级，则使用http拉取
func (o *UrlRich) chromeTimeout(url string) (string, error, bool) {
	if o.downgrading {
		body, err := o.requestByHTTP(url)
		return body, err, true
	}

	return "", errors.New("request from chrome timeout"), false
}
