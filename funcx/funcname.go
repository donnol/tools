package funcx

import "runtime"

func FuncName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	fun := runtime.FuncForPC(pc)
	return fun.Name()
}
