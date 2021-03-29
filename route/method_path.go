package route

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"unicode"
)

// 路径分割符
const (
	pathSep = "/"
)

// addPathPrefix 添加路径前缀
func addPathPrefix(path, prefix string) string {
	if prefix == "" {
		return path
	}
	return pathSep + strings.ToLower(prefix) + path
}

// getMethodPathFromFunc 通过f的名字获取method，path
func getMethodPathFromFunc(f HandlerFunc) (method, path string) {
	// 利用反射和运行时获取函数名
	refValue := reflect.ValueOf(f)
	fn := runtime.FuncForPC(refValue.Pointer())
	fullFuncName := fn.Name()

	return getMethodPath(fullFuncName)
}

// getMethodPath 获取methid, path
func getMethodPath(fullFuncName string) (method, path string) {
	const sep = "."

	upperFunc := func(r rune) bool {
		return unicode.IsUpper(r)
	}

	// 过滤函数名的包名部分
	lastDotIndex := strings.LastIndex(fullFuncName, sep)
	funcName := fullFuncName[lastDotIndex+1:]

	// 找到函数名里的首个大写字母，并以此作为依据将字符串分割
	firstUpperIndex := strings.IndexFunc(funcName, upperFunc)
	if firstUpperIndex == 0 {
		// 如果方法是可导出的，首字母就是大写，需要过滤掉
		firstUpperIndex = strings.IndexFunc(funcName[1:], upperFunc) + 1
	}
	// 如果还是0，说明方法名只有method，没有path
	if firstUpperIndex == 0 {
		method = methodMap(funcName)
		return
	}

	method = funcName[:firstUpperIndex]
	method = methodMap(method)

	// 如果剩下的路径部分还有大写字母，需要分为多段路径
	tmpPath := funcName[firstUpperIndex:]
	for {
		tmpPath = strings.ToLower(tmpPath[:1]) + tmpPath[1:]
		firstUpperIndex = strings.IndexFunc(tmpPath, upperFunc)
		if firstUpperIndex == -1 {
			path += pathSep + strings.ToLower(tmpPath)
			return
		}
		path += pathSep + strings.ToLower(tmpPath[:firstUpperIndex])

		tmpPath = tmpPath[firstUpperIndex:]
	}
}

// methodMap 方法映射
func methodMap(m string) (r string) {
	m = strings.ToUpper(m)
	switch m {
	case "GET":
		r = http.MethodGet
	case "ADD":
		r = http.MethodPost
	case "MOD":
		r = http.MethodPut
	case "DEL":
		r = http.MethodDelete
	default:
		r = m
	}
	return
}
