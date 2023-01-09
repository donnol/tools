package slicex

import "github.com/samber/lo"

// Join nested loop，no limit about condition
func Join[LE, RE, R any](
	left []LE,
	right []RE,
	match func(lele LE, rele RE) bool,
	mapper func(lele LE, rele RE) R,
) []R {
	var result = make([]R, 0, len(left))

	for _, leftEle := range left {
		for _, rightEle := range right {
			if !match(leftEle, rightEle) {
				continue
			}
			result = append(result, mapper(leftEle, rightEle))
		}
	}

	return result
}

// JoinByKey hash join, must be '=' condition
func JoinByKey[K comparable, LE, RE, R any](
	left []LE,
	right []RE,
	leftKey func(item LE) K,
	rightKey func(item RE) K,
	mapper func(lele LE, rele RE) R,
) []R {
	var result = make([]R, 0, len(left))

	rightMap := lo.KeyBy(right, rightKey)

	for _, leftEle := range left {
		exchangeKey := leftKey(leftEle)
		rightEle := rightMap[exchangeKey]
		result = append(result, mapper(leftEle, rightEle))
	}

	return result
}
