package graphviz

import (
	"io"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type Graph struct {
	Name    string
	Edges   map[string]map[string]Edge
	Entitys []Entity
}

func GraphViz(w io.Writer, graphs []Graph) error {
	if err := wrapGraph(w, graphviz.PNG, func(graph *cgraph.Graph) error {
		for _, gra := range graphs {
			subGraph := graph.SubGraph(gra.Name, 1)
			if err := viz(subGraph, gra.Edges, gra.Entitys...); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// Viz 将v画图png输出到w
// edges是从`from entity name`到`to entity name`的边
func Viz(w io.Writer, edges map[string]map[string]Edge, vs ...Entity) error {
	if err := wrapGraph(w, graphviz.PNG, func(graph *cgraph.Graph) error {
		return viz(graph, edges, vs...)
	}); err != nil {
		return err
	}

	return nil
}

func viz(graph *cgraph.Graph, edges map[string]map[string]Edge, vs ...Entity) error {
	graph.SetConcentrate(true)

	nodes := make(map[string]*cgraph.Node, len(vs))
	for _, v := range vs {
		n, err := graph.CreateNode(v.Name())
		if err != nil {
			return err
		}
		n.SetMargin(0)
		n.SetShape(cgraph.BoxShape)
		n.SetLabel(graph.StrdupHTML(v.Label()))

		nodes[v.Name()] = n
	}

	for from, m := range edges {
		for to, edge := range m {
			fromNode := nodes[from]
			toNode := nodes[to]
			e, err := graph.CreateEdge(edge.Name, fromNode, toNode)
			if err != nil {
				return err
			}

			// 从箭头的尾部到头部
			if edge.TailPort != "" {
				e.SetTailPort(edge.TailPort)
			}
			if edge.HeadPort != "" {
				e.SetHeadPort(edge.HeadPort)
			}
		}
	}

	return nil
}

func wrapGraph(w io.Writer, format graphviz.Format, viz func(graph *cgraph.Graph) error) error {
	g := graphviz.New()
	g.SetLayout(graphviz.SFDP)
	graph, err := g.Graph()
	if err != nil {
		return err
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Printf("close graph failed: %v\n", err)
			return
		}
		if err := g.Close(); err != nil {
			log.Printf("close graphviz failed: %v\n", err)
			return
		}
	}()

	if err := viz(graph); err != nil {
		return err
	}

	if err := g.Render(graph, format, w); err != nil {
		return err
	}

	return nil
}
