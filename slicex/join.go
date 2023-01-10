package slicex

import (
	"sort"

	"github.com/samber/lo"
)

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

var (
	joinBNLBufSize = 128
)

// setJoinBNLBufSize size must greater than 0, otherwise use default value.
func setJoinBNLBufSize(size int) {
	if size <= 0 {
		return
	}
	joinBNLBufSize = size
}

// joinBNL bnl是为了减少磁盘读写次数，纯内存操作时意义不大，不建议使用
func joinBNL[LE, RE, R any](
	left []LE,
	right []RE,
	match func(lele LE, rele RE) bool,
	mapper func(lele LE, rele RE) R,
) []R {
	var result = make([]R, 0, len(left))

	joinBuf := make([]LE, 0, joinBNLBufSize)
	for i, leftEle := range left {
		joinBuf = append(joinBuf, leftEle)
		if len(joinBuf) < joinBNLBufSize && i != len(left)-1 {
			continue
		}

		for _, rightEle := range right {
			for _, lele := range joinBuf {
				if !match(lele, rightEle) {
					continue
				}
				result = append(result, mapper(lele, rightEle))
			}
		}
		joinBuf = make([]LE, 0, joinBNLBufSize)
	}

	return result
}

var (
	_ sort.Interface = sortSlice[int]{}
)

type sortSlice[E any] struct {
	eles []E
	l    int
	less func(i, j int) bool
}

func (ss sortSlice[E]) Len() int {
	return ss.l
}
func (ss sortSlice[E]) Less(i, j int) bool {
	return ss.less(i, j)
}
func (ss sortSlice[E]) Swap(i, j int) {
	ss.eles[i], ss.eles[j] = ss.eles[j], ss.eles[i]
}

// joinSortMerge sort merge join
func joinSortMerge[LE, RE, R any](
	left []LE,
	right []RE,
	leftLess func(i, j int) bool,
	rightLess func(i, j int) bool,
	compare func(lele LE, rele RE) int,
	mapper func(lele LE, rele RE) R,
) []R {
	var result = make([]R, 0, len(left))

	var (
		c  int
		ll = len(left)
		rl = len(right)
	)

	lss := sortSlice[LE]{
		eles: left,
		l:    ll,
		less: leftLess,
	}
	sort.Sort(lss)
	rss := sortSlice[RE]{
		eles: right,
		l:    rl,
		less: rightLess,
	}
	sort.Sort(rss)

	lj := 0
	for i := 0; i < ll; i++ {
		leftEle := left[i]

		for j := lj; j < rl; j++ {
			rightEle := right[j]

			c = compare(leftEle, rightEle)
			if c == -1 {
				break
			} else if c == 1 {
				continue
			}

			lj = j
			result = append(result, mapper(leftEle, rightEle))
			break
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
