package errors

import "testing"

func TestError(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		err := NewNormal(0, "错误")
		if e, ok := err.(Error); ok && e.IsNormal() {
			t.Log(err)
		} else {
			t.Fatal("Bad Error")
		}
	})

	t.Run("Fatal", func(t *testing.T) {
		err := NewFatal(0, "错误")
		if e, ok := err.(Error); ok && e.IsFatal() {
			t.Log(err)
		} else {
			t.Fatal("Bad Error")
		}
	})
}
