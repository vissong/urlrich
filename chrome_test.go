package urlrich

import (
	"fmt"
	"testing"
)

func TestUrlRich_UpdateRemoteChromeDebugURL(t *testing.T) {
	type args struct {
		remoteChromeHTTP string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"1", args{"http://127.0.0.1:9222"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &UrlRich{}
			if err := o.updateRemoteChromeWS(tt.args.remoteChromeHTTP); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRemoteChromeDebugURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(o.remoteChromeWS)
		})
	}
}
