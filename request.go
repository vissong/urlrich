package urlrich

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/chromedp/chromedp"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// requestByHTTP 使用http请求网页
func (o *UrlRich) requestByHTTP(url string) (string, error) {
	client := http.Client{
		Timeout: o.timeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("new http request error %s: %v\n", url, err)
		return "", err
	}
	req.Header.Add("User-Agent", o.userAgent)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to download %s: %v\n", url, err)
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
	}

	return body, nil
}

// requestByChromedp 使用 chromedp 抓取网页数据
// 第三个参数标示本次请求是否是 chrome 降级到 http 拉取
func (o *UrlRich) requestByChromedp(url string) (string, error, bool) {

	// 先打开一个新tab
	chromeCtx, chromeCancel := o.NewChromeTab()
	defer chromeCancel()

	// 超时控制
	taskCtx, taskCancel := context.WithTimeout(chromeCtx, o.timeout)
	defer taskCancel()

	// 载入页面
	var body string
	// var charset, contentType string
	// ch := make(chan error, 1)
	err := chromedp.Run(taskCtx,
		// network.Enable(),
		chromedp.Navigate(url),
		chromedp.WaitReady("document"),
		// chromedp.WaitVisible(`html`),
		chromedp.OuterHTML(`document.querySelector("html")`, &body, chromedp.ByJSPath),
		// chromedp.OuterHTML(`document.querySelector("meta[charset]")`, &charset, chromedp.ByJSPath),
		// chromedp.OuterHTML(`document.querySelector("meta[http-equiv='Content-Type']")`, &contentType, chromedp.ByJSPath),
	)

	if err == context.DeadlineExceeded {
		body, err := o.chromeTimeout(url)
		return body, err, true
	}

	// var matched []string
	// reg, _ := regexp.Compile("(utf-8|gbk|gb2312)")
	// if len(charset) > 0 {
	// 	matched = reg.FindAllString(strings.ToLower(charset), -1)
	// }
	// if len(contentType) > 0 {
	// 	matched = reg.FindAllString(strings.ToLower(contentType), -1)
	// }
	// // 如果返回的数据是 gbk 的，则需要进行转码
	// if matched[0] == "gbk" || matched[0] == "gb2312" {
	// 	// utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(body)) // 将gbk转换为utf-8
	// 	// body = string(utf8Data)
	// }

	return body, nil, false
}

// chromeTimeout timeout 处理，如果有降级，则使用http拉取
func (o *UrlRich) chromeTimeout(url string) (string, error) {
	if o.downgrading {
		return o.requestByHTTP(url)
	}

	return "", errors.New("chrome timeout")
}
