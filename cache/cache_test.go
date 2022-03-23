package cache

import (
	"reflect"
	"testing"
	"time"
)

func TestMemImpl(t *testing.T) {
	opt := Option{
		Type: TypeMem,
	}
	c := New(opt)

	for _, cas := range []struct {
		key   string
		value any
	}{
		{"key1", "value1"},
	} {
		r := c.Set(cas.key, cas.value)
		if !r {
			t.Fatalf("Set failed: %v\n", r)
		}

		v := c.Get(cas.key)
		if !reflect.DeepEqual(v, cas.value) {
			t.Fatalf("Get failed: %v != %v\n", v, cas.value)
		}

		v, ok := c.Lookup(cas.key)
		if !ok {
			t.Fatalf("Lookup failed: %v\n", ok)
		}
		if !reflect.DeepEqual(v, cas.value) {
			t.Fatalf("Lookup failed: %v != %v\n", v, cas.value)
		}

		n := time.Duration(5)
		expire := time.Second * n
		r = c.SetNX(cas.key, cas.value, expire)
		if !r {
			t.Fatalf("Setnx failed: %v\n", r)
		}
		// 5秒内获取
		v = c.Get(cas.key)
		if !reflect.DeepEqual(v, cas.value) {
			t.Fatalf("Get failed: %v != %v\n", v, cas.value)
		}

		// 5秒后获取
		time.Sleep(time.Second * (n + 1))
		v = c.Get(cas.key)
		if v != nil {
			t.Fatalf("Get failed: %v != %v\n", v, cas.value)
		}
	}
}

func BenchmarkMemImpl(b *testing.B) {
	opt := Option{
		Type: TypeMem,
	}
	c := New(opt)

	cas := struct {
		key   string
		value any
	}{"key1", "value1"}
	n := time.Duration(5)
	expire := time.Second * n
	for i := 0; i < b.N; i++ {
		c.Set(cas.key, cas.value)

		c.Get(cas.key)

		c.Lookup(cas.key)

		c.SetNX(cas.key, cas.value, expire)
	}
}
