package apitest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/donnol/tools/reflectx"
)

func init() {
	// 随机数种子
	rand.Seed(time.Now().Unix())
}

// JSONIndent json格式化后输出
func JSONIndent(w io.Writer, v any) {
	var b []byte
	if vb, ok := v.([]byte); ok {
		b = vb
	} else if buf, ok1 := v.(*bytes.Buffer); ok1 {
		b = buf.Bytes()
	} else {
		var err error
		b, err = json.Marshal(v)
		if err != nil {
			panic(err)
		}
	}
	var out bytes.Buffer
	if err := json.Indent(&out, b, "", "    "); err != nil {
		panic(err)
	}
	if _, err := out.WriteTo(w); err != nil {
		panic(err)
	}
}

// CookieMapToSlice map转为slice
func CookieMapToSlice(cm map[string]string) []*http.Cookie {
	var cookies []*http.Cookie
	for k, v := range cm {
		tmp := &http.Cookie{
			Name:     k,
			Value:    v,
			HttpOnly: true,
			Path:     "/",
		}
		cookies = append(cookies, tmp)
	}

	return cookies
}

// OpenFile 打开文件
func OpenFile(file, title string) (*os.File, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return f, err
	}
	log.Printf("Open %s file for write api document\n", file)

	if _, err := f.WriteString("# " + title + "\n\n"); err != nil {
		return f, err
	}

	return f, nil
}

func structRandomValue(v any) (any, error) {
	var r any

	vtype, _, err := structTypeValue(v)
	if err != nil {
		return r, err
	}

	r = makeStructWithValue(vtype)

	return r, nil
}

func makeStructWithValue(vtype reflect.Type) any {
	r := compositeStructValue(vtype).Interface()
	return r
}

func compositeStructValue(vtype reflect.Type) reflect.Value {
	var v = reflect.New(vtype)

	fieldKind := vtype.Kind()
	if fieldKind == reflect.Ptr { // 指针
		vtype = vtype.Elem()
		fieldKind = vtype.Kind()
	}
	if fieldKind != reflect.Struct {
		panic("Please input struct or struct pointer")
	}

	// 导出字段
	var sf = make([]reflect.StructField, 0)
	for i := 0; i < vtype.NumField(); i++ {
		field := vtype.Field(i)

		if field.PkgPath != "" { // 忽略非导出字段
			continue
		}

		sf = append(sf, field)
	}
	if len(sf) == 0 {
		return v
	}

	// 以导出字段为基础，新建结构体
	v = structOf(sf)

	return v
}

func structOf(sf []reflect.StructField) reflect.Value {
	var v reflect.Value

	if len(sf) == 0 {
		return v
	}

	typ := reflect.StructOf(sf)
	v = reflect.New(typ).Elem()

	// 为新结构体赋值
	for i := 0; i < v.NumField(); i++ {
		field := sf[i]

		fieldType := field.Type
		originKind := fieldType.Kind()
		originType := fieldType
		if fieldType.Kind() == reflect.Ptr { // 指针
			fieldType = fieldType.Elem()
		}

		if isSimpleKind(fieldType.Kind()) {

			value := randomValue(fieldType.Kind(), 10)
			rvalue := reflect.ValueOf(value)
			if originKind == reflect.Ptr {
				tv := reflect.New(originType.Elem()).Elem() // 新建原始类型指针
				tv.Set(rvalue)                              // 设置值
				v.Field(i).Set(tv.Addr())                   // 将地址赋给field
			} else {
				v.Field(i).Set(rvalue)
			}

		} else { // 非简单类型

			switch fieldType.Kind() {
			case reflect.Struct:
				value := compositeStructValue(fieldType)
				if originKind == reflect.Ptr { // 指针
					tv := reflect.New(originType.Elem()).Elem() // 新建原始类型指针
					tv.Set(value)                               // 设置值
					v.Field(i).Set(tv.Addr())                   // 将地址赋给field
				} else {
					if value.Kind() == reflect.Ptr {
						value = value.Elem()
					}
					v.Field(i).Set(value)
				}
			case reflect.Slice:
				var value reflect.Value
				if isSimpleKind(fieldType.Elem().Kind()) {
					rvalue := randomValue(fieldType.Elem().Kind(), 10)
					value = reflect.ValueOf(rvalue)
				} else if fieldType.Elem().Kind() == reflect.Struct {
					value = compositeStructValue(fieldType.Elem())
				}

				sliceValue := reflect.MakeSlice(fieldType, 0, 0)
				sliceValue = reflect.Append(sliceValue, value)
				v.Field(i).Set(sliceValue)
			default:
				panic("Not support type: " + fieldType.Kind().String())
			}

		}
	}

	return v
}

type SimpleKind interface {
	~int | ~int16 | ~int32 | ~int64 | ~int8 |
		~uint | ~uint16 | ~uint32 | ~uint64 | ~uint8 |
		~string |
		~bool |
		~float32 | ~float64
}

var (
	_ = handleSimpleKind[int]
	_ = isSimpleKind2[int]
	_ = isSimpleKind2[int]
	_ = isSimpleKind3
)

// handleSimpleKind use interface SimpleKind to make sure t is always SimpleKind.
func handleSimpleKind[T SimpleKind](t T) {

	// interface{}可以断言回具体类型，那泛型约束可以断言回具体类型吗？
	// invalid operation: cannot use type assertion on type parameter value t (variable of type T constrained by SimpleKind)
	// _, ok := t.(int)
	// if ok {

	// }

	// type switch?
	switch tt := interface{}(t).(type) { // 需要先转为interface{}，有没有可能不需要呢？-- switch tt := t.(type)
	case int:
		fmt.Printf("[int] value: %d", tt)
	}
}

func isSimpleKind2[T any](t T) bool {
	// invalid operation: cannot use type assertion on type parameter value t (variable of type T constrained by any)
	// _, ok := t.(SimpleKind)
	// if ok {
	// 	return true
	// }

	return false
}

func isSimpleKind3(t any) bool {
	// interface contains type constraints
	// _, ok := t.(SimpleKind)
	// if ok {
	// 	return false
	// }

	return false
}

func isSimpleKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.String, reflect.Bool, reflect.Float32, reflect.Float64:
		return true
	}

	return false
}

func randomValueByTag(tag string) any {
	var v any

	splitList := strings.Split(tag, "=")
	if len(splitList) < 2 {
		return v
	}
	name := splitList[0]     // 名称
	funcCall := splitList[1] // 值方法

	funcName, vs, t := resolveCallExpr(funcCall)
	// log.Printf("funcName: %s, vs: %+v, t: %v\n", funcName, vs, t)
	l := len(vs)
	if l == 0 {
		return v
	}

	switch name {
	case "range":
		d := l % 2
		if d != 0 {
			log.Printf("Please Input Double Args: %v", vs)
			return v
		}
		var rl []any
		for i := 0; i < l; i += 2 {
			switch t {
			case token.INT:
				b := vs[i].(int)
				e := vs[i+1].(int)
				rv := rand.Intn(e-b) + b
				rl = append(rl, rv)
			case token.FLOAT:
				b := vs[i].(float64)
				e := vs[i+1].(float64)
				bi := int(b)
				ei := int(e)
				rv := rand.Intn(ei-bi) + bi
				rfv := float64(rv) + rand.Float64()
				rl = append(rl, rfv)
			default:
				log.Printf("Not support token.Kind for range now: %s", t)
				return v
			}
		}
		i := rand.Intn(len(rl))
		v = rl[i]
	case "enum":
		switch funcName {
		case "one":
			i := rand.Intn(l)
			v = vs[i]
		case "many":
			al := rand.Intn(l + 1)
			var m = make(map[int]bool)
			var vl = make([]any, 0)
			for j := 0; j < al; j++ {
				index := rand.Intn(l)
				if _, ok := m[index]; ok {
					continue
				}
				m[index] = true
				vl = append(vl, vs[index])
			}
			return vl
		default:
			log.Printf("Not support method: %s\n", funcName)
		}
	case "call":
		switch funcName {
		case "year":
			if l != 1 {
				log.Printf("Bad call var number: %d\n", l)
				return v
			}
			v = year(vs[0].(int))
		case "month":
			if l != 2 {
				log.Printf("Bad call var number: %d\n", l)
				return v
			}
			v = month(vs[0].(int), vs[1].(int))
		case "day":
			if l != 3 {
				log.Printf("Bad call var number: %d\n", l)
				return v
			}
			v = day(vs[0].(int), vs[1].(int), vs[2].(int))
		default:
			log.Printf("Not support method now: %s\n", name)
		}
	case "db": // TODO
	case "regexp": // TODO
	default:
		log.Printf("Not support type: %s\n", name)
	}

	return v
}

func year(year int) int64 {
	return time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
}

func month(year, month int) int64 {
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Unix()
}

func day(year, month, day int) int64 {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Unix()
}

// From https://golang.org/pkg/reflect/#FuncOf
func makeFunc() {
	swap := func(in []reflect.Value) []reflect.Value { // 这里定义实际运行的函数体
		if len(in) != 2 {
			panic("in length is not two")
		}
		return []reflect.Value{in[1], in[0]}
	}

	makeSwap := func(fptr any) { // 这里将生产一个符合swap的函数
		fn := reflect.ValueOf(fptr).Elem()

		v := reflect.MakeFunc(fn.Type(), swap)

		fn.Set(v)
	}

	var intSwap func(int, int) (int, int) // 我们最后将使用的变量
	makeSwap(&intSwap)
	fmt.Println(intSwap(0, 1)) // 运行
}

func resolveCallExpr(funcCall string) (name string, v []any, t token.Token) {

	expr, err := parser.ParseExpr(funcCall)
	if err != nil {
		log.Printf("parse vfunc failed, error: %+v\n", err)
		return
	}

	// 函数调用表达式
	callExpr, ok := expr.(*ast.CallExpr)
	if !ok {
		log.Printf("'%s' is not func call\n", funcCall)
		return
	}
	name = callExpr.Fun.(*ast.Ident).Name

	// 参数值和类型
	for _, arg := range callExpr.Args {
		switch argValue := arg.(type) {
		case *ast.BasicLit:
			lit := arg.(*ast.BasicLit)
			t = lit.Kind
			var lv any
			switch t {
			case token.INT:
				lv, _ = strconv.Atoi(lit.Value)
			case token.FLOAT:
				lv, _ = strconv.ParseFloat(lit.Value, 64)
			}
			v = append(v, lv)
		case *ast.Ident:
			ident := arg.(*ast.Ident)
			v = append(v, ident.Name)
			t = token.STRING
		default:
			_ = argValue
		}
	}

	return
}

func randomValue(kind reflect.Kind, l int) any {
	if l <= 0 {
		l = 1
	}
	var v any

	switch kind {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		for i := 0; i < l; i++ {
			v = rand.Int()
		}
	case reflect.String:
		var s = make([]byte, 0)
		for i := 0; i < l; i++ {
			// 可见字符32~126
			// 大写字母65~90
			// 小写字母97~122
			// 数字48~57
			n := rand.Intn(95) + 32
			s = append(s, byte(n))
		}
		v = string(s)
	case reflect.Bool:
		var b bool
		for i := 0; i < l; i++ {
			n := rand.Intn(2)
			if n == 1 {
				b = true
			} else {
				b = false
			}
		}
		v = b
	case reflect.Float32, reflect.Float64:
		for i := 0; i < l; i++ {
			v = rand.Float64()
		}
	case reflect.Complex64, reflect.Complex128:
		fallthrough
	case reflect.Uintptr:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.Interface:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Ptr:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Struct:
		fallthrough
	case reflect.UnsafePointer:
		fallthrough
	case reflect.Invalid:
		fallthrough
	default:
		panic(fmt.Errorf("not support kind: %s", kind))
	}

	return v
}

func collectStructField(vtype reflect.Type) []reflect.StructField {
	var sf = make([]reflect.StructField, 0)

	for i := 0; i < vtype.NumField(); i++ {
		field := vtype.Field(i)

		if field.PkgPath != "" { // 忽略非导出字段
			continue
		}

		if field.Anonymous { // 匿名字段
			fieldType := field.Type
			if field.Type.Kind() == reflect.Ptr { // 指针
				fieldType = field.Type.Elem()
			}
			fieldList := collectStructField(fieldType)
			sf = append(sf, fieldList...)
			continue
		}

		sf = append(sf, field)
	}

	return sf
}

var (
	_ = collectStructField
)

func structToMap(v any) (map[string]any, error) {
	var m = make(map[string]any)

	vtype, vvalue, err := structTypeValue(v)
	if err != nil {
		return m, err
	}
	toMap(vtype, vvalue, m)

	return m, nil
}

func structTypeValue(v any) (vtype reflect.Type, vvalue reflect.Value, err error) {
	if v == nil {
		err = fmt.Errorf("Input is nil")
		return
	}
	vtype = reflect.TypeOf(v)
	vvalue = reflect.ValueOf(v)
	if vtype.Kind() == reflect.Ptr {
		vtype = vtype.Elem()
		vvalue = vvalue.Elem()
	}
	if vtype.Kind() != reflect.Struct {
		err = fmt.Errorf("Input is not struct or struct pointer")
		return
	}

	return
}

func toMap(vtype reflect.Type, vvalue reflect.Value, m map[string]any) {
	// 将field转为map，json tag的值作为键，字段值为值
	for i := 0; i < vtype.NumField(); i++ {
		field := vtype.Field(i)
		value := vvalue.FieldByName(field.Name)

		if field.PkgPath != "" { // 忽略非导出字段
			continue
		}

		// 匿名结构体
		if field.Anonymous {
			fieldType := field.Type
			if field.Type.Kind() == reflect.Ptr { // 指针
				fieldType = field.Type.Elem()
				value = value.Elem()
			}
			toMap(fieldType, value, m)
		} else {
			jsonName := field.Tag.Get("json")
			if jsonName == "-" { // 忽略字段
				continue
			}

			if jsonName == "" { // 使用默认名
				jsonName = strings.ToLower(string(field.Name[0])) // 字段名首字母小写
				if len(field.Name) > 1 {
					jsonName += field.Name[1:]
				}
			} else { // 多个部分时，使用第一个部分
				split := strings.Split(jsonName, ",")
				jsonName = strings.TrimSpace(split[0])
			}

			// 字段类型是指针时，如果值是nil则忽略，如果字段值非nil则取值
			if field.Type.Kind() == reflect.Ptr {
				if value.IsNil() {
					continue
				}
				value = value.Elem()
			}

			m[jsonName] = value.Interface()
		}
	}
}

// func copyResponseBody2(dst []byte, resp *http.Response) (int64, error) {
// 	// 读取body
// 	buf := new(bytes.Buffer)
// 	n, err := io.Copy(buf, resp.Body)
// 	if err != nil {
// 		return n, err
// 	}

// 	// 复制到dst
// 	copy(dst, buf.Bytes())

// 	// 重置resp.Body
// 	resp.Body = ioutil.NopCloser(buf)

// 	return int64(buf.Len()), nil
// }

func copyResponseBody(resp *http.Response) ([]byte, int64, error) {
	if resp == nil {
		return []byte{}, 0, fmt.Errorf("nil response")
	}

	// 读取body
	buf := new(bytes.Buffer)
	n, err := io.Copy(buf, resp.Body)
	if err != nil {
		return []byte{}, n, err
	}

	// 重置resp.Body
	resp.Body = ioutil.NopCloser(buf)

	return buf.Bytes(), int64(buf.Len()), nil
}

// apiKey 获取key
func apiKey(path, method string) string {
	return fmt.Sprintf("%s %s", method, path)
}

func structToList(name string, data ...ate) (string, error) {
	var list string

	list += name + "\n\n"
	var lines string
	for _, d := range data {
		lines += fmt.Sprintf("* `%d` %s\n", d.Code, d.Msg)
	}
	list += lines + "\n"

	return list, nil
}

func structToBlock(name string, data any) (string, error) {
	var block string

	dataStruct, err := reflectx.ResolveStruct(data)
	if err != nil {
		return block, err
	}

	block += name + "\n\n"
	fields := dataStruct.GetFields()
	var level int
	lines := fieldsToLine(level, fields)
	block += lines + "\n"

	return block, nil
}

func fieldsToLine(level int, fields []reflectx.Field) string {
	var lines string
	for _, field := range fields {
		var fieldName, fieldTypeName, fieldComment string

		// 是否内嵌结构体
		isEmbed := field.StructField.Anonymous

		// 字段名
		fieldName = field.StructField.Tag.Get("json")
		if fieldName == "-" {
			continue
		}
		if fieldName == "" {
			fieldName = field.StructField.Name
		}

		// 字段类型
		fieldType := field.StructField.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		switch fieldType.Kind() {
		case reflect.Struct:
			fieldTypeName = "object"
		case reflect.Slice:
			sliceType := fieldType.Elem()
			if sliceType.Kind() == reflect.Struct {
				fieldTypeName = "object"
			} else {
				fieldTypeName = sliceType.String()
			}
			fieldTypeName += " list"
		default:
			fieldTypeName = fieldType.Kind().String()
		}

		// 字段注释
		fieldComment = field.Comment

		// 添加一行
		if !isEmbed { // 如果是内嵌结构体，不需要添加该行
			line := fmt.Sprintf("%s %s (*%s*) %s\n", linePrefix(level), fieldName, fieldTypeName, fieldComment)
			lines += line
		}

		// 结构体，切片等复合结构，需要继续遍历，并且在写入时向内缩进
		switch fieldType.Kind() {
		case reflect.Struct, reflect.Slice:
			newLevel := level
			if !isEmbed { // 如果是内嵌结构体，不需要向内缩进
				newLevel = level + 1
			}
			innerLines := fieldsToLine(newLevel, field.Struct.Fields)
			lines += innerLines
		}
	}

	return lines
}

func linePrefix(level int) string {
	var empty = " "
	var prefix = "*"
	if level == 0 {
		return prefix
	}

	var r string
	for i := 0; i < level*4; i++ {
		r += empty
	}

	return r + prefix
}

func dataToSummary(name string, data []byte, isJSON bool) string {
	var summary string

	summary += `<details>
<summary>` + name + `</summary>` + "\n\n```json\n"
	if isJSON {
		var buf = new(bytes.Buffer)
		if data != nil {
			JSONIndent(buf, data)
		}
		summary += buf.String()
	} else {
		summary += string(data)
	}
	summary += "\n```\n\n" + `</details>` + "\n\n"

	return summary
}
