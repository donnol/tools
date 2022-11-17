package comparable

// From https://github.com/golang/go/issues/56548#issuecomment-1317673963

// we want to ensure that T is strictly comparable
type T struct {
	x int
}

// define a helper function with a type parameter P constrained by T
// and use that type parameter with isComparable
// -- 定义一个带有被T约束的类型参数P的辅助函数以在编译时确保类型T是可比较的
func TisComparable[P T]() {
	_ = isComparable[P]
}

func isComparable[_ comparable]() {}
