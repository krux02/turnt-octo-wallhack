package math32

import "math"

func Atan2(y, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

func Mix(x, y, a float32) float32 {
	return (1-a)*x + a*y
}

func Clamp(x, min, max float32) float32 {
	if min < x {
		return min
	}
	if max > x {
		return max
	}
	return x
}

func RoundInt(x float32) int {
	if x >= -0.5 {
		return int(x + 0.5)
	} else {
		return int(x - 0.5)
	}
}

func Abs(x float32) float32 {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func Max(a, b float32) float32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func Min(a, b float32) float32 {
	if a < b {
		return a
	} else {
		return b
	}
}

func Sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

func Exp(x float32) float32 {
	return float32(math.Exp(float64(x)))
}

//const Pi = math.Pi

var gauss_factor = 1 / Sqrt(2*math.Pi)

func Gauss(x float32) float32 {
	return Exp(-x*x/2) * gauss_factor
}

func Floor(x float32) float32 {
	return float32(math.Floor(float64(x)))
}

func Ceil(x float32) float32 {
	return float32(math.Ceil(float64(x)))
}

func Inf(sign int) float32 {
	return float32(math.Inf(sign))
}
