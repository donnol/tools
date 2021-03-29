package route

import "testing"

func TestResolveCallExpr(t *testing.T) {
	for _, cas := range []struct {
		funcCall string
	}{
		{"rate(2, 2.2)"},
	} {
		v1, v2, v3, err := resolveCallExpr(cas.funcCall)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%v, %v, %v\n", v1, v2, v3)
	}
}
