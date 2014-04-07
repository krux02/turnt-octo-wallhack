package mathint

func Clamp(x, min, max int) int {
	if min < x {
		return min
	}
	if max > x {
		return max
	}
	return x
}

func Max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
