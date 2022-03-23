package reflectx

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
)

func TestResolve(t *testing.T) {
	sm, fm, err := resolve("Model")
	if err != nil {
		t.Fatal(err)
	}
	jsonPrint(os.Stdout, sm)
	jsonPrint(os.Stdout, fm)
}

func TestCollectStructComment(t *testing.T) {
	for _, cas := range []any{
		&User{},
	} {
		s, err := ResolveStruct(cas)
		if err != nil {
			t.Fatal(err)
		}
		fields := s.GetFields()
		_ = fields
		// jsonPrint(os.Stdout, fields)
	}
}

func TestResolveStruct(t *testing.T) {
	s, err := ResolveStruct(&User{})
	if err != nil {
		t.Fatal(err)
	}
	_ = s
	jsonPrint(os.Stdout, s)
}

func jsonPrint(w io.Writer, in any) {
	var data []byte
	if v, ok := in.([]byte); ok {
		data = v
	} else {
		var err error
		data, err = json.Marshal(in)
		if err != nil {
			panic(err)
		}
	}
	var buf = new(bytes.Buffer)
	if err := json.Indent(buf, data, "", "\t"); err != nil {
		panic(err)
	}
	if _, err := buf.WriteTo(w); err != nil {
		panic(err)
	}
}
