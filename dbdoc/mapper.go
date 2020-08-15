package dbdoc

import (
	"strings"
	"unicode"
)

const (
	tablePrefix = "t_"
)

// tableMapper 表名映射
func tableMapper(name string) string {
	name = fieldMapper(name)
	return tablePrefix + name
}

// fieldTypeMapper 将结构体字段类型映射为数据库表字段类型
func fieldTypeMapper(fieldType string) string {
	var result = fieldType
	switch fieldType {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		result = "integer"
	case "string":
		result = "text"
	case "float32", "float64":
		result = "numeric"
	case "Time":
		result = "timestamp with time zone"
	case "Jsonb":
		result = "jsonb"
	case "UUID":
		result = "uuid"
	case "string[]":
		result = "text[]"
	}
	return result
}

// fieldMapper 将结构体字段映射为数据库表字段
func fieldMapper(structFieldName string) string {
	var result string

	// 找出每个部分首字母的位置
	var fieldIndexs []int
	var isUpper bool
	for i := 0; i < len(structFieldName); i++ {
		// 1. 首先判断字符是否是大写
		// 2. 如果是大写，则判断前面的字符是否是小写并且后面的字符是否是大写
		if !unicode.IsLower(rune(structFieldName[i])) { // 大写
			if !isUpper {
				fieldIndexs = append(fieldIndexs, i)
				isUpper = true
			}
		} else {
			if isUpper {
				fieldIndexs = append(fieldIndexs, i-1)
			}
			isUpper = false
		}
	}

	// 分割
	var offset int
	var fields []string
	for _, index := range fieldIndexs {
		tmpIndex := index - offset
		tmpField := structFieldName[:tmpIndex]
		if strings.TrimSpace(tmpField) == "" {
			continue
		}
		fields = append(fields, tmpField)
		structFieldName = structFieldName[tmpIndex:]
		offset = index
	}
	fields = append(fields, structFieldName) // 最后剩余的部分

	// 分割后各个部分转为小写
	for i, field := range fields {
		fields[i] = strings.ToLower(field)
	}

	// 使用下划线拼接到一起
	result = strings.Trim(strings.Join(fields, "_"), "_")

	return result
}
