package urlrich

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chromedp/chromedp"
)

// requestByHTTP 使用http请求网页
func (o *UrlRich) requestByHTTP() (string, error) {
	client := http.Client{
		Timeout: o.timeout,
	}
	req, err := http.NewRequest("GET", o.url, nil)
	if err != nil {
		log.Fatalf("new http request error %s: %v\n", o.url, err)
		return "", err
	}
	req.Header.Add("User-Agent", o.userAgent)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to download %s: %v\n", o.url, err)
		return "", err
	}
	defer resp.Body.Close()

	html, _ := ioutil.ReadAll(resp.Body)

	return string(html), nil
}

// requestByChromedp 使用 chromedp 抓取网页数据
func (o *UrlRich) requestByChromedp() (string, error) {

	taskCtx, cancel := context.WithTimeout(o.chromedpCtx, o.timeout)
	defer cancel()

	var body string
	err := chromedp.Run(taskCtx,
		chromedp.Navigate(o.url),
		chromedp.WaitReady("document"),
		// chromedp.WaitVisible("body"),
		chromedp.OuterHTML(`document.querySelector("html")`, &body, chromedp.ByJSPath),
	)
	if err != nil {
		log.Fatalf("failed to request %s: %v\n", o.url, err)
		return "", err
	}

	return body, nil
}
