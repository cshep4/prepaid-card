package underscore

// Range is 生成区间数据
func Range(start, stop, step int) []int {
	arr := make([]int, 0)
	if step == 0 {
		return arr
	}

	if start == stop {
		return arr
	} else if start < stop {
		if step < 0 {
			return arr
		}
	} else {
		if step > 0 {
			return arr
		}
	}

	for ; start < stop; start = start + step {
		arr = append(arr, start)
	}
	return arr
}

// Range is Queryer's method
func (q *Query) Range(start, stop, step int) Queryer {
	q.source = Range(start, stop, step)
	return q
}
