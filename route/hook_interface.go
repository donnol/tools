package route

import "context"

// Newer 新建
type Newer interface {
	New() any
}

// Checker 检查接口
type Checker interface {
	Check(context.Context) error
}

// Filter 过滤器
type Filter interface {
	Filter() any
}
