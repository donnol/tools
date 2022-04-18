package domain

import "testing"

func TestUser_ModName(t *testing.T) {
	type fields struct {
		id   uint
		name string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "modName", fields: fields{id: 1, name: "jd"}, args: args{name: "jk"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				id:   tt.fields.id,
				name: tt.fields.name,
			}
			u.ModName(tt.args.name)
			if u.name != tt.args.name {
				t.Fatalf("bad result: %s != %s", u.name, tt.args.name)
			}
		})
	}
}
