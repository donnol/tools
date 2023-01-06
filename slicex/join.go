package slicex

import "github.com/samber/lo"

// Join nested loop，no limit about condition
func Join[J, K, R any](
	left []J,
	right []K,
	match func(J, K) bool,
	mapper func(J, K) R,
) []R {
	var r = make([]R, 0, len(left))

	for _, j := range left {
		for _, k := range right {
			if !match(j, k) {
				continue
			}
			r = append(r, mapper(j, k))
		}
	}

	return r
}

// JoinByKey hash join, must be '=' condition
func JoinByKey[K comparable, LE, RE, R any](
	left []LE,
	right []RE,
	lk func(item LE) K,
	rk func(item RE) K,
	mapper func(LE, RE) R,
) []R {
	var r = make([]R, 0, len(left))

	rm := lo.KeyBy(right, rk)

	for _, j := range left {
		k := lk(j)
		re := rm[k]
		r = append(r, mapper(j, re))
	}

	return r
}
