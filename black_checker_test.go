package urlrich

import (
	"testing"
)

func TestUrlRich_IsBlackDomain(t *testing.T) {

	blackDomains := []string{
		"*.inner.com",
		"123.inner2.com",
	}
	blackIPs := []string{
		"10.*.*.*",
	}

	type fields struct {
		blackDomains []string
		blackIps     []string
	}
	type args struct {
		rawurl string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"inner.com", args{"https://www.inner.com"}, true, false},
		{"inner2.com", args{"123.inner2.com"}, true, false},
		{"a.inner2.com", args{"a.inner2.com"}, false, false},
		{"baidu", args{"baidu.com"}, false, false},
		{"ip", args{"https://10.1.2.3"}, true, false},
		{"ip2", args{"10.1.2.3"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &UrlRich{}
			o.SetBlack(blackDomains, blackIPs)

			got, err := o.IsBlack(tt.args.rawurl)
			// fmt.Println(err)

			if (err != nil) != tt.wantErr {
				t.Errorf("IsBlackDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsBlackDomain() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrlRich_SetBlack(t *testing.T) {

	blackDomains := []string{
		"*.inner.com",
		"123.inner2.com",
	}
	blackIPs := []string{
		"10.*.*.*",
	}

	type args struct {
		domains []string
		ips     []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{blackDomains, blackIPs}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &UrlRich{}
			if err := o.SetBlack(tt.args.domains, tt.args.ips); (err != nil) != tt.wantErr {
				t.Errorf("SetBlack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
