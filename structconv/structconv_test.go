package structconv

import (
	"reflect"
	"testing"
)

// 虽然通过名字来赋值，省下了功夫，但现实中字段对应往往不是简单的通过名字就能实现，还有其它因素，所以实用性并不太高
type (
	UserReq struct {
		Phone string
	}

	UserResp struct {
		Name string
		Age  uint
	}

	UserTable struct {
		Id    uint
		Name  string
		Age   uint
		Phone string
	}

	Article struct {
		Name     string
		UserName string
	}
)

func TestConvByFieldName(t *testing.T) {
	from := UserReq{
		Phone: "12345678901",
	}
	to := &UserTable{}
	ConvByFieldName(from, to)

	if to.Phone != from.Phone {
		t.Fatalf("converse failed: %s != %s\n", to.Phone, from.Phone)
	}

	to.Id = 1
	to.Name = "jd"
	to.Age = 18

	to2 := &UserResp{}
	ConvByFieldName(to, to2)

	if to2.Name != to.Name {
		t.Fatalf("converse failed: %s != %s\n", to2.Name, to.Name)
	}
	if to2.Age != to.Age {
		t.Fatalf("converse failed: %d != %d\n", to2.Age, to.Age)
	}
}

func TestConvListByFieldName(t *testing.T) {
	type args struct {
		from []UserReq
		to   []*UserTable
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				from: []UserReq{
					{
						Phone: "12345678901",
					},
					{
						Phone: "234",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.to = MakeSliceWithNew[UserReq, UserTable](tt.args.from)

			ConvSliceByFieldName(tt.args.from, tt.args.to)

			for i := range tt.args.from {
				if tt.args.to[i].Phone == "" || tt.args.to[i].Phone != tt.args.from[i].Phone {
					t.Fatalf("converse failed: %s != %s\n", tt.args.to[i].Phone, tt.args.from[i].Phone)
				}
			}
		})
	}
}

func TestConvSliceByFunc(t *testing.T) {
	type args struct {
		from []UserReq
		conv func(UserReq) UserTable
	}
	tests := []struct {
		name string
		args args
		want []UserTable
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				from: []UserReq{
					{
						Phone: "12345678901",
					},
					{
						Phone: "234",
					},
				},
				conv: func(ur UserReq) UserTable {
					ut := &UserTable{}
					ConvByFieldName(ur, ut)
					return *ut
				},
			},
			want: []UserTable{
				{Phone: "12345678901"},
				{Phone: "234"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvSliceByFunc(tt.args.from, tt.args.conv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvSliceByFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
