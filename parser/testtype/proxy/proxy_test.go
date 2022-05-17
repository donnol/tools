package proxy

import (
	"testing"
)

func TestC(t *testing.T) {
	t.Run("普通调用", func(t *testing.T) {
		r1, err := A(1, 1, "a", "b")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("r1: %v\n", r1)
	})

	// t.Run("Proxy调用", func(t *testing.T) {
	// 	r1, err := AProxy(1, 1, "a", "b") // 在对A调用的前后添加耗时统计
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	t.Logf("r1: %v\n", r1)
	// })
	t.Run("C", func(t *testing.T) {
		C()
	})
}
