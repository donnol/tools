package httpreq

import (
	"io"
	"net/http"
	"reflect"
	"testing"
)

type (
	result struct {
		BgQuality int `json:"BgQuality"`
	}
)

func TestDoHTTPRequest(t *testing.T) {
	type args[R any] struct {
		client        *http.Client
		method        string
		link          string
		body          io.Reader
		codeChecker   CodeChecker
		extractResult ResultExtractor[result]
	}
	type testCase[R any] struct {
		name    string
		args    args[R]
		want    result
		wantErr bool
	}
	tests := []testCase[result]{
		// TODO: Add test cases.
		{
			name: "1",
			args: args[result]{
				method:        http.MethodGet,
				link:          "https://www.bing.com/hp/api/model",
				body:          nil,
				codeChecker:   CodeIs200,
				extractResult: JSONExtractor[result],
			},
			want: result{BgQuality: 50},
		},
		{
			name: "nil",
			args: args[result]{
				method:        http.MethodGet,
				link:          "https://www.bing.com/hp/api/model",
				body:          nil,
				codeChecker:   nil,
				extractResult: nil,
			},
			want: result{BgQuality: 50},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Do(tt.args.client, tt.args.method, tt.args.link, tt.args.body, tt.args.codeChecker, tt.args.extractResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoHTTPRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoHTTPRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
