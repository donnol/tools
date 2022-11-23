package basetype

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPassword_Encrypt(t *testing.T) {
	tests := []struct {
		name    string
		p       Password
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "123",
			p:       "jd123XXX",
			wantErr: false,
		},
		{
			name:    "t123@mgr",
			p:       "t123@mgr",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPp, err := tt.p.Encrypt()
			if (err != nil) != tt.wantErr {
				t.Errorf("Password.Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := tt.p.Compare(gotPp); err != nil {
				t.Fatalf("compare failed: %s is not hash value of %s", gotPp, tt.p)
			}
		})
	}
}

func TestPassword_String(t *testing.T) {
	tests := []struct {
		name string
		p    Password
		want string
	}{
		// TODO: Add test cases.
		{
			name: "p",
			p:    "123",
			want: "*",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("Password.String() = %v, want %v", got, tt.want)
			}

			buf := bytes.NewBuffer([]byte(""))
			_, err := fmt.Fprintf(buf, "%s", tt.p.String())
			if err != nil {
				t.Errorf("printf failed: %v", err)
			}
			if buf.String() != tt.want {
				t.Errorf("bad case: %s != %s", buf.String(), tt.want)
			}
		})
	}
}
