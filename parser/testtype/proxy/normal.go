package proxy

import (
	"log"
)

// === TODO: 任意函数的Around ===

// 用户编写的代码
func A(ctx any, id int, args ...string) (string, error) {
	log.Printf("arg, ctx: %v, id: %v, args: %+v\n", ctx, id, args)
	return "A", nil
}

func C() {
	r1, err := A(1, 1, "a", "b")
	if err != nil {
		log.Printf("err: %v\n", err)
		return
	}
	log.Printf("r1: %v\n", r1)
}
