package proxy

import (
	"log"
	"time"
)

// === TODO: 任意函数的Around ===

// 用户编写的代码
func A(ctx any, id int, args ...string) (string, error) {
	log.Printf("arg, ctx: %v, id: %v, args: %+v\n", ctx, id, args) // call1
	return "A", nil
}

// 正常函数
func C() {
	r1, err := A(1, 1, "a", "b") // call2
	if err != nil {
		log.Printf("err is not nil: %v\n", err) // call3
		return
	}
	log.Printf("r1: %v\n", r1) // call4
}

// 生成Proxy之后的函数
func C2() {
	// 编译之前，通过重写ast，改为调用B（B内再调用A）
	// 1 遍历源码，找到函数调用（可配置规则，以过滤出想要改变的函数）- *ast.CallExpr

	// 2 生成一个对应的附加了额外逻辑的函数B（B内调用A）- *ast.FuncDecl
	// 生成出来的代码
	// 有没有办法生成一个B，使得C调B，跟C调A一样，但是又可以在B里添加额外的逻辑呢？
	B := func(ctx any, id int, args ...string) (string, error) {
		// 为了要支持添加额外的逻辑，显然不能直接返回
		// return A(ctx, id, args...)

		// 添加逻辑
		begin := time.Now() // call5

		// 根据签名，不难生成出以下返回值定义
		var r1 string
		var r2 error

		r1, r2 = A(ctx, id, args...) // call6

		// 添加逻辑
		log.Printf("used time: %v\n", time.Since(begin)) // call7, call8

		return r1, r2
	}

	// 3 将此处对A的调用替换为对B的调用 -
	// r1, err := A(1, 1, "a", "b") // call2
	r1, err := B(1, 1, "a", "b") // call2
	if err != nil {
		log.Printf("err is not nil: %v\n", err) // call3
		return
	}
	log.Printf("r1: %v\n", r1) // call4
}
