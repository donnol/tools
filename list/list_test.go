package list

import (
	"reflect"
	"testing"
)

func TestListFilter(t *testing.T) {
	for _, cas := range []struct {
		in     []string
		factor string
		want   []string
	}{
		{[]string{"a", "b", "c"}, "a", []string{"b", "c"}},
		{[]string{"a", "b", "c"}, "b", []string{"a", "c"}},
		{[]string{"a", "b", "c"}, "c", []string{"a", "b"}},
	} {
		r := Filter(cas.in, cas.factor)
		if !reflect.DeepEqual(r, cas.want) {
			t.Fatalf("Bad result: %+v != %+v\n", r, cas.want)
		}
	}
}
