package urlrich

import (
	"log"
	"regexp"
	"strings"

	"github.com/goware/urlx" // urlx 可以解析不带协议的 url,也可以直接解析ip
)

// SetBlackDomains 设置黑名单域名与IP列表，设置后，会根据列表，生成正则
func (o *UrlRich) SetBlack(domains []string, ips []string) error {

	// 域名正则
	// *.ab.com to .*\.ab.com
	for _, v := range domains {
		regx := strings.Replace(v, "*.", `.*\.`, 1)
		o.blackDomains = append(o.blackDomains, regx)
	}

	// ip 正则
	// 10.*.*.* to 10\.\d+\.\d+\.\d+
	for _, v := range ips {
		regx := strings.Replace(v, ".*", `\.\d+`, -1)
		o.blackIps = append(o.blackIps, regx)
	}

	return nil
}

// IsBlack 判断域名是否命中黑名单(先判断域名，后判断ip，ip黑名单需要先解析才行，所以最好先通过域名配置好黑名单)
// 支持完整配 www.inner.com, 10.14.87.167
// 支持通配如：*.inner.com, 10.*.*.*
func (o *UrlRich) IsBlack(rawUrl string) (bool, error) {

	if len(o.blackDomains) == 0 && len(o.blackIps) == 0 {
		return false, nil
	}

	url, err := urlx.Parse(rawUrl)
	if err != nil {
		return false, err
	}
	domain := url.Hostname()

	// 域名判断
	for _, v := range o.blackDomains {
		matched, err := regexp.MatchString(v, domain)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}

	ip, err := urlx.Resolve(url)
	// 解析ip失败，当做命中黑名单处理
	if err != nil {
		log.Println(err)
		return true, nil
	}

	for _, v := range o.blackIps {
		matched, err := regexp.MatchString(v, ip.String())
		// 正则失败，说明可能是规则问题，返回未命中
		if err != nil {
			log.Println(err)
			return false, nil
		}
		if matched {
			return true, nil
		}
	}

	return false, nil
}
