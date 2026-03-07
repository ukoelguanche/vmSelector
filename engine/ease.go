package engine

import "math"

func EeaseLinear(t float64) float64 {
	return t
}

func EaseInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - math.Pow(-2*t+2, 3)/2
}

func EaseInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return -1 + (4-2*t)*t
}

func EaseOutQuad(t float64) float64 {
	return t * (2 - t)
}

func EaseInQuad(t float64) float64 {
	return t * t
}

func EaseOutCubic(t float64) float64 {
	t--
	return t*t*t + 1
}
