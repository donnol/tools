package dbdoc

import (
	"fmt"
	"io"
	"strings"

	"github.com/donnol/tools/reflectx"
)

// Table 表
type Table struct {
	Name        string  // 名字
	Comment     string  // 注释
	Description string  // 描述
	FieldList   []Field // 字段列表
	IndexList   []Index // 索引列表

	tableMapper Mapper // 表名映射函数
	mapper      Mapper // 字段映射函数
	typeMapper  Mapper // 字段类型映射函数

	doc []byte

	//  er图
	Subgraph string // 子图
	Edge     string // 边
}

// Mapper 映射函数类型
type Mapper func(string) string

// Field 字段
type Field struct {
	Name        string // 名字
	Type        string // 类型
	Nullable    bool   // 是否可null
	Primary     bool   // 是否主键
	Description string // 描述
	Index              // 索引
	Relation           // 关系
}

// Index 索引
type Index struct {
	Name      string  // 名字
	Unique    bool    // 是否唯一索引
	FieldList []Field // 涉及字段
}

// Relation 关系
type Relation struct {
	TableName  string // 表名字
	TableField string // 表字段
}

// NewTable 新建表
func NewTable() *Table {
	return &Table{
		FieldList:   make([]Field, 0),
		IndexList:   make([]Index, 0),
		tableMapper: tableMapper,
		mapper:      fieldMapper,
		typeMapper:  fieldTypeMapper,
	}
}

// New 新建
func (t *Table) New() *Table {
	return NewTable()
}

// Resolve 解析
func (t *Table) Resolve(v interface{}) *Table {
	const relationTagName = "rel"
	const relationTagSep = "."

	var err error
	var vstruct reflectx.Struct
	vstruct, err = reflectx.ResolveStruct(v)
	must(err)

	vstructName := vstruct.Name
	nameList := strings.Split(vstructName, ".")
	tableName := nameList[len(nameList)-1]
	if tableName == "Entity" {
		// 获取包名
		pkgNameList := strings.Split(nameList[len(nameList)-2], "/")
		tableName = pkgNameList[len(pkgNameList)-1]
		t.Name = tablePrefix + tableName
	} else {
		t.Name = t.tableMapper(tableName)
	}
	t.Comment = vstruct.Comment
	t.Description = vstruct.Description

	var tf Field
	for _, sf := range vstruct.Fields {
		// 忽略匿名结构体/接口
		if sf.Anonymous {
			continue
		}

		// 主键
		if sf.Name == "ID" {
			tf = Field{
				Name:        "id",
				Type:        t.typeMapper(sf.Type.Kind().String()),
				Primary:     true,
				Description: sf.Comment,
			}
		} else {
			tf = FieldByTag(sf.Tag)
			if tf.Name == "" {
				tf.Name = t.mapper(sf.Name)
			}
			if tf.Type == "" {
				tf.Type = t.typeMapper(sf.Type.Kind().String())
			}
			tf.Description = sf.Comment

			// 关系
			relTagValue, ok := sf.Tag.Lookup(relationTagName)
			if ok {
				relTagValues := strings.Split(relTagValue, relationTagSep)
				if len(relTagValues) != 2 {
					must(fmt.Errorf(`请指定rel标签的表和字段，如: 'rel:"user.id"'`))
				}
				tf.Relation.TableName = relTagValues[0]
				tf.Relation.TableField = relTagValues[1]
			}
		}

		t.FieldList = append(t.FieldList, tf)
	}

	return t
}

// SetComment 设置注释
func (t *Table) SetComment(comment string) *Table {
	t.Comment = comment
	return t
}

// SetDescription 设置描述
func (t *Table) SetDescription(description string) *Table {
	t.Description = description
	return t
}

// Write 写入f
func (t *Table) Write(w io.Writer) *Table {
	t = t.makeDoc()

	_, err := w.Write(t.doc)
	must(err)

	return t
}

// SetMapper 设置字段名映射方法
func (t *Table) SetMapper(f Mapper) *Table {
	t.mapper = f
	return t
}

// SetTypeMapper 设置类型映射方法
func (t *Table) SetTypeMapper(f Mapper) *Table {
	t.typeMapper = f
	return t
}

func (t *Table) makeDoc() *Table {

	leftAngle := "<"
	rightAngle := ">"
	format := "## %s\n\n%s\n\n%s\n%s索引：\n\n%s\n"
	fieldFormat := "| %s | %s | %v | %v | %s |\n"
	header := "| Field | Type | Nullable | Primary | Description |\n| :-: | :-: | :-: | :-: | :-: |\n"
	indexFormat := "* %s(%s: %s)\n"

	// 字段
	var field, index, fieldEnum string
	for _, tf := range t.FieldList {
		var nullableString, primaryString, description string
		if tf.Nullable {
			nullableString = "*"
		}
		if tf.Primary {
			primaryString = "*"
		}
		description = tf.Description
		if strings.Contains(description, leftAngle) &&
			strings.Contains(description, rightAngle) &&
			strings.Index(description, leftAngle) < strings.Index(description, rightAngle) {

			if fieldEnum != "" {
				fieldEnum += "\n"
			}
			fieldEnumLine := description[strings.Index(description, leftAngle)+1 : strings.Index(description, rightAngle)-1]
			fieldEnumList := strings.Split(fieldEnumLine, ":")
			for fei, fe := range fieldEnumList {
				if fei == 0 {
					fieldEnum += fe + ":\n\n"
				}
				if fei > 0 {
					fieldEnumValueList := strings.Split(fe, ";")
					for _, fev := range fieldEnumValueList {
						fieldEnum += fmt.Sprintf("* %s\n", strings.TrimRight(fev, ";"))
					}
				}
			}

			description = description[:strings.Index(description, leftAngle)]
		}
		field += fmt.Sprintf(fieldFormat, tf.Name, tf.Type, nullableString, primaryString, description)

		// 索引
		var indexName = tf.Index.Name
		var uniqueString string
		if tf.Index.Unique {
			uniqueString = "UNIQUE"
			if indexName == "" {
				indexName = tf.Name
			}
		}
		if indexName != "" {
			index += fmt.Sprintf(indexFormat, uniqueString, indexName, tf.Name)
		}
	}

	description := t.Description + "\n\n" + fieldEnum
	if description != "" {
		description += "\n\n"
	}
	content := fmt.Sprintf(format, t.Name, t.Comment, header+field, description, index)
	t.doc = []byte(content)

	return t
}

// dot脚本
var (
	// %s分别是多个子图，边样式和边
	GraphFormat = `digraph "Database Structure" {
		label = "ER Diagram";
		labelloc = t;
		compound = true;
		node [ shape = record ];
		fontname = "Helvetica";
		ranksep = 1.25;
		ratio = 0.7;
		rankdir = LR;
		%s
		%s
		%s
	}`
	// %s分别是表名，表名，标签列表
	subgraphFormat = `
		subgraph "table_%s" {
			node [ shape = "plaintext" ]
			"%s" [ %s 
			]
		}
			`
	// %s分别是表名和列的列表
	subgraphLabelFormat = `label=<
				<TABLE BORDER="0" CELLSPACING="0" CELLBORDER="1">
				<TR><TD COLSPAN="3" BGCOLOR="#DDDDDD">%s</TD></TR>
				%s
				</TABLE>>
	`
	// %s分别是字段名和字段类型
	subgraphLabelRowFormat = `
				<TR><TD COLSPAN="3" PORT="%s">%s:<FONT FACE="Helvetica-Oblique" POINT-SIZE="10">%s</FONT></TD></TR>
	`
	EdgeStyleFormat = `edge [ arrowtail=normal, style=dashed, color="#444444" ]
	`
	// %s分别是源表名，源字段名，目标表名，目标字段名
	edgeFormat = `
		%s:%s -> %s:%s
	`
)

// MakeGraph 生成图
func (t *Table) MakeGraph() *Table {
	// 字段
	var label, edge string
	for _, tf := range t.FieldList {
		label += fmt.Sprintf(subgraphLabelRowFormat, tf.Name, tf.Name, tf.Type)

		if tf.Relation.TableName != "" && tf.Relation.TableField != "" {
			edge += fmt.Sprintf(edgeFormat, t.Name, tf.Name, tf.Relation.TableName, tf.Relation.TableField)
		}
	}

	subgraphLabel := fmt.Sprintf(subgraphLabelFormat, t.Name, label)
	subgraph := fmt.Sprintf(subgraphFormat, t.Name, t.Name, subgraphLabel)

	t.Subgraph = subgraph
	t.Edge = edge

	return t
}
