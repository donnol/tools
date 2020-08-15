package dbdoc

import (
	"reflect"
	"strings"
)

// FieldByTag 从tag解析field
func FieldByTag(tag reflect.StructTag) Field {
	var result Field

	gormTag := tag.Get("gorm")
	if !strings.Contains(gormTag, "NOT NULL") {
		result.Nullable = true
	}
	if strings.Contains(gormTag, "UNIQUE") {
		result.Index.Unique = true
	}
	tagList := strings.Split(gormTag, ";")
	for _, v := range tagList {
		if strings.Contains(v, "type:") {
			typeList := strings.Split(v, ":")
			result.Type = typeList[1]
		}
		if strings.Contains(v, "column:") {
			nameList := strings.Split(v, ":")
			result.Name = nameList[1]
		}
		if strings.Contains(v, "index:") {
			indexList := strings.Split(v, ":")
			result.Index.Name = indexList[1]

			// 添加下面这句会导致错误：
			// runtime: goroutine stack exceeds 1000000000-byte limit
			// fatal error: stack overflow
			// runtime stack:
			// runtime.throw(0x55bd85, 0xe)
			// 		/usr/local/go/src/runtime/panic.go:616 +0x81
			// runtime.newstack()
			// 		/usr/local/go/src/runtime/stack.go:1054 +0x71f
			// runtime.morestack()
			// 		/usr/local/go/src/runtime/asm_amd64.s:480 +0x89

			// result.Index.FieldList = append(result.Index.FieldList, result)
		}
	}

	return result
}
