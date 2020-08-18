package urlrich

import (
	"time"
)

// Option urlrich option
type Option func(o *UrlRich)

// WithUrl 设置需要取抓取的 url，正常应该在 Do 的时候传入，而不是在这里设置
func WithUrl(url string) Option {
	return func(o *UrlRich) {
		o.url = url
	}
}

// WithDebug 设置是否打开debug，会影响 chromedp 的debuglog
func WithDebug(debug bool) Option {
	return func(o *UrlRich) {
		o.debug = debug
	}
}

// WithUserAgent 设置UA头
// 默认为：Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36
func WithUserAgent(ua string) Option {
	return func(o *UrlRich) {
		o.userAgent = ua
	}
}

// WithUseChrome 设置是否使用 chromedp，如果是不是则使用http直接抓取，不支持使用js渲染的页面
// 但是使用 chrome 抓取，会比直接 http 抓取更慢
func WithUseChrome(useChrome bool) Option {
	return func(o *UrlRich) {
		o.useChrome = useChrome
	}
}

// WithTimeout 设置抓取时候超时时间（http为网络，chrome为网络+dom渲染）
func WithTimeout(timeout time.Duration) Option {
	return func(o *UrlRich) {
		o.timeout = timeout
	}
}
