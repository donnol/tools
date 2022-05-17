package proxy

import (
	"log"
)

func A(ctx any, id int, args ...string) (string, error) {
	log.Printf("arg, ctx: %v, id: %v, args: %+v\n", ctx, id, args)
	return "A", nil
}
func C() {
	args := []string{"a", "b", "c", "d"}
	r1, err := AProxy(1, 1, args...)
	if err != nil {
		log.Printf("err: %v\n", err)
		return
	}
	log.Printf("r1: %v\n", r1)
}
