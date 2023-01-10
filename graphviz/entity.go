package graphviz

import (
	"fmt"
	"strings"

	"github.com/donnol/tools/reflectx"
	"github.com/samber/lo"
)

type Entity interface {
	Name() string
	Label() string
	Ports() []string
}

type entity struct {
	name  string
	label string

	ports []string
}

func (ent *entity) Name() string {
	return ent.name
}

func (ent *entity) Label() string {
	return ent.label
}

func (ent *entity) Ports() []string {
	return ent.ports
}

type Edge struct {
	Name     string
	TailPort string
	HeadPort string
}

func ToEntity(v any) Entity {
	var (
		labelHead = `<table border="0" cellborder="1" cellspacing="0" cellpadding="4">
	`
		labelTitle = `<tr><td port="title_name_%s" bgcolor="lightblue">%s</td></tr>
	`
		labelField = `<tr><td port="%s" align="left">%s: %s</td></tr>
`
		labelFoot = `</table>`
	)

	var label string
	label += labelHead

	stru, err := reflectx.ResolveStruct(v)
	if err != nil {
		panic(err)
	}
	name := stru.Name
	index := strings.LastIndex(name, ".")
	if index != -1 {
		name = name[index+1:]
	}
	label += fmt.Sprintf(labelTitle, strings.ToLower(name), name)

	fields := stru.GetFields()
	ports := make([]string, 0, len(fields))
	for _, field := range fields {
		name := field.StructField.Name
		port := name

		ports = append(ports, port)
		label += fmt.Sprintf(labelField, port, name, field.Type)
	}

	label += labelFoot

	return &entity{
		name:  name,
		label: label,
		ports: ports,
	}
}

func ToEntityBatch(vs ...any) []Entity {
	r := make([]Entity, 0, len(vs))
	for _, v := range vs {
		r = append(r, ToEntity(v))
	}
	return r
}

func MakeEdge(from, to Entity, tail, head string) Edge {
	findex := lo.IndexOf(from.Ports(), tail)
	if findex == -1 {
		panic(fmt.Errorf("tail %s not exist in %+v", tail, from.Ports()))
	}
	tindex := lo.IndexOf(to.Ports(), head)
	if tindex == -1 {
		panic(fmt.Errorf("head %s not exist in %+v", head, to.Ports()))
	}

	return Edge{
		Name:     fmt.Sprintf("%s:%s->%s:%s", from.Name(), to.Name(), tail, head),
		TailPort: tail,
		HeadPort: head,
	}
}
