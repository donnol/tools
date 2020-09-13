package reflectx

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"github.com/donnol/tools/importpath"
	"github.com/donnol/tools/parser"
)

// Field 字段
type Field struct {
	reflect.StructField        // 内嵌反射结构体字段类型
	Comment             string // 注释
	Struct              Struct // 字段的类型是其它结构体
}

// Struct 结构体
type Struct struct {
	Name        string       // 名字
	Comment     string       // 注释
	Description string       // 描述
	Type        reflect.Type // 反射类型
	Fields      []Field      // 结构体字段
}

// MakeStruct 新建结构体
func MakeStruct() Struct {
	return Struct{
		Fields: make([]Field, 0),
	}
}

// ResolveStruct 解析结构体
func ResolveStruct(value interface{}) (Struct, error) {
	s := MakeStruct()

	var refType reflect.Type
	if v, ok := value.(reflect.Type); ok {
		refType = v
	} else {
		refType = reflect.TypeOf(value)
	}
	s.Type = refType

	if refType == nil {
		return s, fmt.Errorf("nil refType")
	}

	if refType.Kind() == reflect.Ptr { // 指针
		refType = refType.Elem()
	}
	if refType.Kind() != reflect.Struct {
		return s, fmt.Errorf("bad value type , type is %v", refType.Kind())
	}
	structName := refType.PkgPath() + "." + refType.Name()
	s.Name = structName

	if refType.NumField() == 0 { // 空结构体
		return s, nil
	}
	if err := collectStructComment(refType, &s); err != nil {
		return s, err
	}

	return s, nil
}

// GetFields 返回结构体的所有field，包括匿名字段的field
func (s Struct) GetFields() []Field {
	return getFields(s)
}

func getFields(s Struct) []Field {
	var fields = make([]Field, 0)
	for _, f := range s.Fields {
		if f.Anonymous {
			fields = append(fields, getFields(f.Struct)...)
		} else {
			fields = append(fields, f)
		}
	}

	return fields
}

// collectStructComment 收集结构体的注释
func collectStructComment(refType reflect.Type, s *Struct) error {
	// 解析-获取结构体注释
	var r map[string]string
	var f map[string]string
	var err error
	if r, f, err = resolve(s.Name); err != nil {
		return fmt.Errorf("resolve output failed, error is %v", err)
	}
	s.Comment = r[commentKey]
	s.Description = r[descriptionKey]

	// 内嵌结构体
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)

		sf := Field{
			StructField: field,
			Comment:     f[field.Name],
		}

		fieldType := field.Type
		if field.Anonymous { // 匿名
			// 忽略匿名接口
			if fieldType.Kind() != reflect.Interface {
				sf.Struct, err = ResolveStruct(fieldType)
				if err != nil {
					return err
				}
			}
		}
		// 非匿名结构体类型
		if fieldType.Kind() == reflect.Ptr ||
			fieldType.Kind() == reflect.Slice ||
			fieldType.Kind() == reflect.Map ||
			fieldType.Kind() == reflect.Chan ||
			fieldType.Kind() == reflect.Array {
			fieldType = fieldType.Elem()
		}
		// 忽略time.Time
		if fieldType.Kind() == reflect.Struct && fieldType != reflect.TypeOf((*time.Time)(nil)).Elem() {
			sf.Struct, err = ResolveStruct(fieldType)
			if err != nil {
				return err
			}
		}

		s.Fields = append(s.Fields, sf)
	}

	return nil
}

const (
	structStart    = "type"
	structEnd      = "}"
	fieldSep       = " "
	commentSep     = "//"
	commentKey     = "comment"
	descriptionKey = "description"
)

func resolve(structName string) (map[string]string, map[string]string, error) {
	return resolveWithParser(structName)
}

func resolveWithParser(structName string) (map[string]string, map[string]string, error) {
	var structCommentMap = make(map[string]string)
	var fieldCommentMap = make(map[string]string)

	ip := &importpath.ImportPath{}
	path, err := ip.GetByCurrentDir()
	if err != nil {
		return structCommentMap, fieldCommentMap, err
	}

	name := structName
	dotIndex := strings.LastIndex(structName, ".")
	if dotIndex != -1 {
		if dotIndex != 0 {
			path = structName[:dotIndex]
		}
		name = structName[dotIndex+1:]
	}

	parser := parser.New(parser.Option{})
	structs, err := parser.ParseAST(path)
	if err != nil {
		return structCommentMap, fieldCommentMap, err
	}

	var exist bool
	for _, oneStruct := range structs {
		if oneStruct.Name != name {
			continue
		}
		exist = true
		structCommentMap[commentKey] = strings.TrimSpace(oneStruct.Comment)
		structCommentMap[descriptionKey] = strings.TrimSpace(oneStruct.Doc)
		for _, field := range oneStruct.Fields {
			fieldCommentMap[field.Name] = strings.TrimSpace(field.Comment)
		}
	}
	if !exist {
		fmt.Printf("Can't find comment info of %s", structName)
	}

	return structCommentMap, fieldCommentMap, nil
}

// 返回结构体注释，字段名注释映射和错误
func resolveWithGoDoc(structName string) (map[string]string, map[string]string, error) {
	var structCommentMap = make(map[string]string)
	var fieldCommentMap = make(map[string]string)

	// 运行go doc命令
	cmd := exec.Command("go", "doc", structName)
	output, err := cmd.Output()
	if err != nil {
		// FIXME go doc失败的结构体先忽略, 如: time.zoneTrans
		// log.Printf("go doc failed, struct name is %s, error is %v", structName, err)
		return structCommentMap, fieldCommentMap, nil
	}

	var isEnd bool
	var nowLineNo, commentLineNo, endNo int
	var commentLine string
	buf := bytes.NewBuffer(output)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return structCommentMap, fieldCommentMap, fmt.Errorf("go doc failed, struct name is %s, error is %v", structName, err)
		}
		nowLineNo++

		if strings.TrimSpace(line) == "" {
			continue
		}

		if strings.TrimSpace(line) == structEnd {
			isEnd = true
			continue
		}

		if isCommentLine(line) {
			if commentLineNo == 0 ||
				commentLineNo == nowLineNo-1 {
				commentLineNo = nowLineNo
				commentLine += fmt.Sprintf("<%s>", strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(line), commentSep)))
			}
		}

		var comment string
		pieceList := strings.Split(line, commentSep)
		if !isEnd {
			keyList := strings.Split(strings.TrimSpace(pieceList[0]), fieldSep)
			if len(keyList) == 1 { // 匿名结构体
				continue
			}
			key := keyList[0]
			if key == structStart {
				continue
			}
			if len(pieceList) == 2 {
				comment = strings.TrimSpace(pieceList[1])
			}

			fieldCommentMap[key] = comment + commentLine
			commentLineNo = 0
			commentLine = ""
		} else {
			endNo++

			if endNo == 2 {
				if _, ok := structCommentMap[commentKey]; !ok {
					// 不是func
					if strings.Index(line, "func") != 0 {
						structCommentMap[commentKey] = strings.TrimSpace(line)
					}
				}
			}
			if endNo > 2 {
				if _, ok := structCommentMap[descriptionKey]; !ok {
					// 不是func
					if strings.Index(line, "func") != 0 {
						structCommentMap[descriptionKey] += strings.TrimSpace(line)
					}
				}
			}
		}
	}

	return structCommentMap, fieldCommentMap, nil
}

func isCommentLine(line string) bool {
	line = strings.TrimSpace(line)
	if strings.Index(line, commentSep) == 0 {
		return true
	}
	return false
}
