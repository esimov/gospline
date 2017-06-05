package spline

import (
	"image/draw"
	"image/color"
	"math"
)

type Canvas struct {
	draw.Image
	Matrix
}

func (canvas *Canvas) DrawLine(x1, y1, x2, y2 float64, col color.Color, antialiased bool) *Canvas {
	if antialiased {
		xiaolinWuLine(canvas, x1, y1, x2, y2, col.(color.NRGBA))
	} else {
		bresenhamLine(canvas, x1, y1, x2, y2, col.(color.NRGBA))
	}
	return canvas
}

// Bresenham's line algorithm
// http://en.wikipedia.org/wiki/Bresenham's_line_algorithm
func bresenhamLine(canvas *Canvas, x1, y1, x2, y2 float64, col color.NRGBA) *Canvas {
	x1, y1 = canvas.TransformPoint(x1, y1)
	x2, y2 = canvas.TransformPoint(x2, y2)

	abs := func(i float64) float64 {
		if i < 0 {
			return -i
		}
		return i
	}
	steep := abs(y1-y2) > abs(x1-x2)
	if steep {
		x1, y1 = y1, x1
		x2, y2 = y2, x2
	}
	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	dx := x2 - x1
	dy := abs(y2 - y1)
	err := dx / 2
	y := y1

	var ystep float64 = -1
	if y1 < y2 {
		ystep = 1
	}

	for x := x1; x <= x2; x++ {
		if steep {
			canvas.Set(int(y), int(x), color.NRGBA{col.R, col.G, col.B, col.A})
		} else {
			canvas.Set(int(x), int(y), color.NRGBA{col.R, col.G, col.B, col.A})
		}
		err -= dy
		if err < 0 {
			y += ystep
			err += dx
		}
	}
	return canvas
}

// Xialin Wu's antialiased line algorithm
// https://en.wikipedia.org/wiki/Xiaolin_Wu's_line_algorithm
func xiaolinWuLine(canvas *Canvas, x1, y1, x2, y2 float64, col color.NRGBA) *Canvas {
	ipart := func(x float64) float64 {
		return math.Floor(x)
	}
	round := func(x float64) float64 {
		return ipart(x + 0.5)
	}
	fpart := func(x float64) float64 {
		if x < 0 {
			return 1 - (x - ipart(x))
		}
		return x - ipart(x)
	}
	rfpart := func(x float64) float64 {
		return 1 - fpart(x)
	}
	x1, y1 = canvas.TransformPoint(x1, y1)
	x2, y2 = canvas.TransformPoint(x2, y2)
	dx := x2 - x1
	dy := y2 - y1
	ax := dx
	if ax < 0 {
		ax = -ax
	}
	ay := dy
	if ay < 0 {
		ay = -ay
	}

	var plot func(int, int, float64)

	if ax < ay {
		x1, y1 = y1, x1
		x2, y2 = y2, x2
		dx, dy = dy, dx
		plot = func(x, y int, c float64) {
			canvas.Set(y, x, color.NRGBA{col.R, col.G, col.B, uint8(255 * c)})
		}
	} else {
		plot = func(x, y int, c float64) {
			canvas.Set(x, y, color.NRGBA{col.R, col.G, col.B, uint8(255 * c)})
		}
	}
	if x2 < x1 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	gradient := dy / dx

	// handle first endpoint
	xend := round(x1)
	yend := y1 + gradient * (xend-x1)
	xgap := rfpart(x1 + 0.5) * 1.5
	xpxl1 := int(xend)
	ypxl1 := int(ipart(yend))

	plot(xpxl1, ypxl1, rfpart(yend) * xgap)
	plot(xpxl1, ypxl1+1, fpart(yend) * xgap)
	intery := yend + gradient

	// handle second endpoint
	xend = round(x2)
	yend = y2 + gradient * (xend-x2)
	xgap = fpart(x2 + 0.8)
	xpxl2 := int(xend)
	ypxl2 := int(ipart(yend))

	plot(xpxl2, ypxl2, rfpart(yend) * xgap)
	plot(xpxl2, ypxl2+1, fpart(yend) * xgap)

	for x := xpxl1 + 1; x <= xpxl2-1; x++ {
		plot(x, int(ipart(intery)), rfpart(intery))
		plot(x, int(ipart(intery))+1, fpart(intery))
		intery = intery + gradient
	}
	return canvas
}
