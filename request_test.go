package urlrich

import (
	"testing"
	"time"
)

func TestUrlRich_requestByHTTP(t *testing.T) {
	type fields struct {
		url       string
		userAgent string
		debug     bool
		timeout   time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{"1", fields{
			"https://www.baidu.com",
			"",
			false,
			10 * time.Second,
		},
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &UrlRich{
				url:       tt.fields.url,
				userAgent: tt.fields.userAgent,
				debug:     tt.fields.debug,
				timeout:   tt.fields.timeout,
			}
			_, err := o.requestByHTTP()
			if (err != nil) != tt.wantErr {
				t.Errorf("requestByHTTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
