package dbdoc

import (
	"fmt"
	"io"
)

// Resolve 解析多个结构体，并将它们写到w
func Resolve(w io.Writer, v ...any) error {
	var table = NewTable()

	for _, s := range v {
		table.New().Resolve(s).Write(w)
	}

	return nil
}

// ResolveGraph 解析多个结构体，生成图脚本，并将它们写到w
func ResolveGraph(w io.Writer, v ...any) error {
	var table = NewTable()

	var content, subgraph, edge string
	for _, s := range v {
		t := table.New().Resolve(s).MakeGraph()

		subgraph += t.Subgraph
		edge += t.Edge
	}
	content = fmt.Sprintf(GraphFormat, subgraph, EdgeStyleFormat, edge)

	_, err := w.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}
