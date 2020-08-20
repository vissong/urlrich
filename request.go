package urlrich

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chromedp/chromedp"
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

	return string(html), nil
}

// requestByChromedp 使用 chromedp 抓取网页数据
func (o *UrlRich) requestByChromedp(url string) (string, error) {

	// 先打开一个新tab
	chromeCtx, chromeCancel := o.NewChromeTab()
	defer chromeCancel()

	// 超时控制
	taskCtx, taskCancel := context.WithTimeout(chromeCtx, o.timeout)
	defer taskCancel()

	// 载入页面
	var body string
	err := chromedp.Run(taskCtx,
		chromedp.Navigate(url),
		chromedp.WaitReady("document"),
		// chromedp.WaitVisible("body"),
		chromedp.OuterHTML(`document.querySelector("html")`, &body, chromedp.ByJSPath),
	)
	if err != nil {
		log.Fatalf("failed to request %s: %v\n", url, err)
		return "", err
	}

	return body, nil
}
