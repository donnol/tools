package sqlparser

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
)

func ParseCreateSQL(sql string) *Struct {
	s := &Struct{}

	node, err := parse(sql)
	if err != nil {
		log.Fatal(err)
	}
	(*node).Accept(s)

	return s
}

func parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()

	stmtNodes, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		return nil, err
	}

	return &stmtNodes, nil
}

type colX struct {
	colNames []string
}

func (v *colX) Enter(in ast.Node) (ast.Node, bool) {
	if name, ok := in.(*ast.ColumnName); ok {
		v.colNames = append(v.colNames, name.Name.O)
	}
	return in, false
}

func (v *colX) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func extract(rootNode *ast.StmtNode) []string {
	v := &colX{}
	(*rootNode).Accept(v)
	return v.colNames
}

const (
	structHeadTmpl = `
	// {{.StructName}} {{.StructComment}}
	type {{.StructName}} struct {
	`
	structFieldTmpl = `	{{.FieldName}} {{.FieldType}} {{.FieldTag}} // {{.FieldComment}}
	`
	structFootTmpl      = `}`
	structTableNameTmpl = `
	func ({{.StructName}}) TableName() string {
		return "{{.TableName}}"
	}`
)

type Struct struct {
	Name    string
	Comment string
	Fields  []Field
}
type Field struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

func (v *Struct) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.CreateTableStmt:
		v.Name = node.Table.Name.O
		for _, opt := range node.Options {
			if opt.Tp == ast.TableOptionComment {
				v.Comment = opt.StrValue
			}
		}
		for _, col := range node.Cols {
			field := Field{
				Name: col.Name.Name.O,
			}

			field.Type = col.Tp.InfoSchemaStr()
			field.Type = processFieldType(field.Type)

			for _, opt := range col.Options {
				switch opt.Tp {
				case ast.ColumnOptionPrimaryKey:
				case ast.ColumnOptionNotNull:
				case ast.ColumnOptionAutoIncrement:
				case ast.ColumnOptionDefaultValue:
				case ast.ColumnOptionUniqKey:
				case ast.ColumnOptionNull:
				case ast.ColumnOptionOnUpdate:
				case ast.ColumnOptionFulltext:
				case ast.ColumnOptionComment:
					field.Comment = opt.Expr.(ast.ValueExpr).GetDatumString()
				case ast.ColumnOptionGenerated:
				case ast.ColumnOptionReference:
				case ast.ColumnOptionCollate:
				case ast.ColumnOptionCheck:
				case ast.ColumnOptionColumnFormat:
				case ast.ColumnOptionStorage:
				case ast.ColumnOptionAutoRandom:
				}
			}
			v.Fields = append(v.Fields, field)
		}
	}
	return in, false
}

func (v *Struct) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func processFieldType(fieldType string) string {
	// 去掉括号及其内部内容
	re := regexp.MustCompile(`([(\d)])`)
	fieldType = re.ReplaceAllString(fieldType, "")
	return fieldType
}

type Option struct {
	StructNameMapper func(string) string         // 名称映射
	IgnoreField      []string                    // 忽略字段
	FieldNameMapper  func(string) string         // 字段名称映射
	FieldTypeMapper  func(string) string         // 字段类型映射
	FieldTagMapper   func(string, string) string // 可根据名称和类型自行决定字段tag
}

var (
	doption = Option{
		StructNameMapper: func(k string) string {
			// snake to camel
			return strcase.ToCamel(k)
		},
		FieldNameMapper: func(k string) string {
			return strcase.ToCamel(k)
		},
		FieldTypeMapper: func(k string) string {
			var (
				// type mapper
				m = map[string]string{
					"BIT":               "bool",
					"TEXT":              "string",
					"BLOB":              "[]byte",
					"DATE":              "time.Time",
					"DATETIME":          "time.Time",
					"DECIMAL":           "float64",
					"DOUBLE":            "float64",
					"ENUM":              "",
					"FLOAT":             "float32",
					"GEOMETRY":          "",
					"MEDIUMINT":         "int",
					"JSON":              "json.RawMessage",
					"UNSIGNED INT":      "uint",
					"INT UNSIGNED":      "uint",
					"INT":               "int",
					"LONGTEXT":          "string",
					"LONGBLOB":          "[]byte",
					"UNSIGNED BIGINT":   "uint64",
					"BIGINT UNSIGNED":   "uint64",
					"BIGINT":            "int64",
					"MEDIUMTEXT":        "string",
					"MEDIUMBLOB":        "[]byte",
					"NULL":              "",
					"SET":               "",
					"UNSIGNED SMALLINT": "uint16",
					"SMALLINT UNSIGNED": "uint16",
					"SMALLINT":          "int16",
					"BINARY":            "[]byte",
					"CHAR":              "string",
					"TIME":              "time.Time",
					"TIMESTAMP":         "time.Time",
					"UNSIGNED TINYINT":  "uint8",
					"TINYINT UNSIGNED":  "uint8",
					"TINYINT":           "int8",
					"TINYTEXT":          "string",
					"TINYBLOB":          "[]byte",
					"VARBINARY":         "[]byte",
					"VARCHAR":           "string",
					"YEAR":              "time.Time",
				}
			)
			if v, ok := m[strings.ToUpper(k)]; ok {
				return v
			}
			return k
		},
		FieldTagMapper: func(name string, typ string) string {
			if name == "" {
				return ""
			}
			jname := strcase.ToCamel(name)
			jname = strings.ToLower(string(jname[0])) + jname[1:]
			return fmt.Sprintf("`json:\"%s\" db:\"%s\"`", jname, name)
		},
	}
)

func (opt *Option) fillByDefault() {
	if opt.StructNameMapper == nil {
		opt.StructNameMapper = doption.StructNameMapper
	}
	if opt.FieldNameMapper == nil {
		opt.FieldNameMapper = doption.FieldNameMapper
	}
	if opt.FieldTypeMapper == nil {
		opt.FieldTypeMapper = doption.FieldTypeMapper
	}
	if opt.FieldTagMapper == nil {
		opt.FieldTagMapper = doption.FieldTagMapper
	}
}

func (s *Struct) Gen(w io.Writer, opt Option) error {
	(&opt).fillByDefault()

	name := s.Name
	if opt.StructNameMapper != nil {
		name = opt.StructNameMapper(name)
	}
	{
		temp, err := template.New("structHead").Parse(structHeadTmpl)
		if err != nil {
			return err
		}
		if err := temp.Execute(w, map[string]any{
			"StructName":    name,
			"StructComment": s.Comment,
		}); err != nil {
			return err
		}
	}

	{
		for _, field := range s.Fields {
			fieldName := field.Name
			if opt.FieldNameMapper != nil {
				fieldName = opt.FieldNameMapper(fieldName)
			}
			fieldType := field.Type
			if opt.FieldTypeMapper != nil {
				fieldType = opt.FieldTypeMapper(fieldType)
			}
			fieldTag := field.Tag
			if opt.FieldTagMapper != nil {
				fieldTag = opt.FieldTagMapper(field.Name, field.Type)
			}
			temp, err := template.New("structField").Parse(structFieldTmpl)
			if err != nil {
				return err
			}
			if err := temp.Execute(w, map[string]any{
				"FieldName":    fieldName,
				"FieldType":    fieldType,
				"FieldTag":     fieldTag,
				"FieldComment": field.Comment,
			}); err != nil {
				return err
			}
		}
	}

	{
		if _, err := w.Write([]byte(structFootTmpl)); err != nil {
			return err
		}
	}

	{
		temp, err := template.New("structTableName").Parse(structTableNameTmpl)
		if err != nil {
			return err
		}
		if err := temp.Execute(w, map[string]any{
			"StructName": name,
			"TableName":  s.Name,
		}); err != nil {
			return err
		}
	}
	return nil
}
