package field

import "testing"

func Test_checkerImpl_CheckField(t *testing.T) {
	type fields struct {
		pkgs []string
	}
	type args struct {
		pkg string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "", fields: fields{pkgs: []string{"github.com/fishedee/tools/query"}}, args: args{pkg: "github.com/donnol/tools/staticcheck/test_data"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &checkerImpl{
				pkgs: tt.fields.pkgs,
			}
			if err := impl.CheckField(tt.args.pkg); (err != nil) != tt.wantErr {
				t.Errorf("checkerImpl.CheckField() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
