package reflectx

import "testing"

// M M
type M struct {
	Name string
}

func TestIsStructPointer(t *testing.T) {
	var i int
	for _, cas := range []struct {
		In   any
		Want bool
	}{
		{&M{}, true},
		{M{}, false},
		{&i, false},
	} {
		r := IsStructPointer(cas.In)
		if r != cas.Want {
			t.Fatalf("%v != %v\n", r, cas.Want)
		}
	}
}
