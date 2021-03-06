package urlrich

import (
	"time"
)

// Option urlrich option
type Option func(o *UrlRich)

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

// WithDowngrading 设置降级
func WithDowngrading(downgrading bool) Option {
	return func(o *UrlRich) {
		o.downgrading = downgrading
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

// WithRemoteChrome 设置远程 chrome 的ws地址
func WithRemoteChrome(httpUrl string) Option {
	return func(o *UrlRich) {
		o.useRemoteChrome = true
		o.useChrome = true
		o.remoteChromeHTTP = httpUrl
	}
}

// WitchHttpProxy 设置http代理
func WithHttpProxy(proxyUrl string) Option {
	return func(o *UrlRich) {
		o.httpProxy = proxyUrl
	}
}
