package drawing

import "math"

type Point struct {
	X, Y int
}

func lerp(t float64, x1, y1, x2, y2 int) (x, y int) {
	if t < 0 || t > 1 {
		panic("Cannot lerp with t outside [0, 1]")
	}

	xf := float64(x1)*(1-t) + float64(x2)*t
	yf := float64(y1)*(1-t) + float64(y2)*t

	x = int(math.Round(xf))
	y = int(math.Round(yf))

	return x, y
}
