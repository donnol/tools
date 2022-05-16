package proxy

import (
	"log"
	"time"
)

func AProxy(ctx any, id int, args ...string) (string, error) {
	begin := time.Now()

	var r0 string
	var r1 error

	r0, r1 = A(ctx, id, args...)

	log.Printf("used time: %v\n", time.Since(begin))

	return r0, r1
}
