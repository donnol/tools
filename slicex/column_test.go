package slicex

import (
	"reflect"
	"testing"
)

func TestColumn(t *testing.T) {
	type args struct {
		slice  []int
		column func(item int) int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				slice: []int{1, 2, 3},
				column: func(item int) int {
					return item + 1
				},
			},
			want: []int{2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Column(tt.args.slice, tt.args.column); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Column() = %v, want %v", got, tt.want)
			}
		})
	}
}
