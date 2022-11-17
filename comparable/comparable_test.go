package comparable

import "testing"

var (
	_ = T{x: 1}
)

func TestTisComparable(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TisComparable()
		})
	}
}
