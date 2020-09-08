package urlrich

import (
	"reflect"
	"testing"
)

func TestUrlRich_Do(t *testing.T) {

	u1, _ := New(WithUseChrome(true))
	u2, _ := New()

	type args struct {
		url string
	}
	tests := []struct {
		name    string
		u       *UrlRich
		args    args
		want    RichResult
		wantErr bool
	}{
		{"baidu with chrome", u1, args{url: "https://www.baidu.com"}, RichResult{Desc: "全球最大的中文搜索引擎、致力于让网民更便捷地获取信息，找到所求。百度超过千亿的中文网页数据库，可以瞬间找到相关的搜索结果。"}, false},
		// {"baidu with http", u2, args{url: "https://www.baidu.com"}, RichResult{Desc: "百度一下，你就知道 全球最大的中文搜索引擎、致力于让网民更便捷地获取信息，找到所求。百度超过千亿的中文网页数据库，可以瞬间找到相关的搜索结果。"}, false},
		// {"baidu with http err", u2, args{url: "https//www.baidu.com"}, RichResult{Desc: ""}, true},
		{"youdao with chrome", u1, args{url: "https://xue.youdao.com/sw/m/1946563?keyfrom=dict2.index"}, RichResult{Desc: "现在的考生写作文，大部分人都套着模板。"}, false},
		{"youdao with http", u2, args{url: "https://xue.youdao.com/sw/m/1946563?keyfrom=dict2.index"}, RichResult{Desc: ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := tt.u
			got, err := o.Do(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Desc, tt.want.Desc) {
				t.Errorf("Do() got = %v, want %v", got, tt.want)
			}
		})
	}
}
