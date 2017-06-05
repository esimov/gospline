// Ported from https://github.com/Tagussan/BSpline
package spline

import (
	"math"
)

type bspline struct {
	points [][]float64
	degree int
	copy bool
}

var (
	pts [][]float64
	degree int
	dimension int
	baseFunc func(x float64)float64
	baseFuncRangeInt int
)

func NewBSpline(points [][]float64, degree int, copy bool)*bspline {
	return &bspline{
		points,
		degree,
		copy,
	}
}

func (self *bspline) Init() {
	if self.copy {
		pts = make([][]float64, len(self.points))
		for i := 0; i < len(self.points); i++ {
			pts = append(pts, self.points[i])
		}
	} else {
		pts = self.points
	}
	degree = self.degree
	dimension = len(self.points[0])

	if degree == 2 {
		baseFunc = self.baseDeg2
		baseFuncRangeInt = 2
	} else if degree == 3 {
		baseFunc = self.baseDeg3
		baseFuncRangeInt = 2
	} else if degree == 4 {
		baseFunc = self.baseDeg4
		baseFuncRangeInt = 3
	} else if degree == 5 {
		baseFunc = self.baseDeg5
		baseFuncRangeInt = 3
	}
}

func (self *bspline) seqAt(dim int) func(int)float64 {
	margin := self.degree + 1

	return func(n int)float64 {
		if n < margin {
			return pts[0][dim]
		} else if len(pts) + margin <= n {
			return pts[len(pts)-1][dim]
		} else {
			return pts[n-margin][dim]
		}

	}
}

func (self *bspline) baseDeg2(x float64)float64 {
	if x >= -0.5 && x < 0.5 {
		return 0.75 - x * x
	} else if x >= 0.5 && x <= 1.5 {
		return 1.125 + (-1.5 + x/2.0) * x
	} else if x >= -1.5 && x < -0.5 {
		return 1.125 + (1.5 + x/2.0) * x
	} else {
		return 0.0
	}
}

func (self *bspline) baseDeg3(x float64)float64 {
	if x >= -1.0 && x < 0 {
		return 2.0/3.0 + (-1.0 - x/2.0) * x*x
	} else if x >= 1 && x <= 2 {
		return 4.0/3.0 + x * (-2.0 + (1.0 - x/6.0) * x)
	} else if x >= -2.0 && x < -1 {
		return 4.0/3.0 + x * (2.0 + (1.0 + x/6.0) * x)
	} else if x >= 0 && x < 1 {
		return 2.0/3.0 + (-1.0 + x/2.0) * x*x
	} else {
		return 0.0
	}
}

func (self *bspline) baseDeg4(x float64)float64 {
	if x >= -1.5 && x < -0.5 {
		return 55.0/96.0 + x * (-(5.0/24.0) + x * (-(5.0/4.0) + (-(5.0/6.0) - x/6.0) * x))
	} else if x >= 0.5 && x < 1.5 {
		return 55.0/96.0 + x * (5.0/24.0 + x * (-(5.0/4.0) + (5.0/6.0 - x/6.0)*x))
	} else if x >= 1.5 && x <= 2.5 {
		return 625.0/384.0 + x * (-(125.0/48.0) + x * (25.0/16.0 + (-(5.0/12.0) + x/24.0) * x))
	} else if x >= -2.5 && x <= -1.5 {
		return 625.0/384.0 + x * (125.0/48.0 + x * (25.0/16.0 + (5.0/12.0 + x/24.0) * x))
	} else if x >= -1.5 && x < 1.5 {
		return 115.0/192.0 + x*x * (-(5.0/8.0) + x*x/4.0)
	} else {
		return 0.0
	}
}

func (self *bspline) baseDeg5(x float64)float64 {
	if x >= -2.0 && x < -1 {
		return 17.0/40.0 + x*(-(5.0/8.0) + x*(-(7.0/4.0) + x*(-(5.0/4.0) + (-(3.0/8.0) - x/24.0)*x)))
	} else if x >= 0 && x < 1 {
		return 11.0/20.0 + x*x*(-(1.0/2.0) + (1.0/4.0 - x/12.0)*x*x)
	} else if x >= 2 && x <= 3 {
		return 81.0/40.0 + x*(-(27.0/8.0) + x*(9.0/4.0 + x*(-(3.0/4.0) + (1.0/8.0 - x/120.0)*x)))
	} else if x >= -3 && x < -2 {
		return 81.0/40.0 + x*(27.0/8.0 + x*(9.0/4.0 + x*(3.0/4.0 + (1.0/8.0 + x/120.0)*x)))
	} else if x >= 1 && x < 2 {
		return 17.0/40.0 + x*(5.0/8.0 + x*(-(7.0/4.0) + x*(5.0/4.0 + (-(3.0/8.0) + x/24.0)*x)))
	} else if x >= -1 && x < 0 {
		return 11.0/20.0 + x*x*(-(1.0/2.0) + (1.0/4.0 + x/12.0)*x*x)
	} else {
		return 0.0
	}
}

func (self *bspline) getInterpol(seq func(int)float64, t float64)float64 {
	tInt := int(math.Floor(t))
	var result float64
	for i := tInt - baseFuncRangeInt; i <= tInt + baseFuncRangeInt; i++ {
		result += seq(i) * baseFunc(t-float64(i))
	}
	return result
}

func (self *bspline) Interpolate(t, roundOn float64) []float64 {
	t = t * ((float64(degree) + 1.0) * 2.0 + float64(len(pts))) // t must be between 0...1

	if dimension == 2 {
		return []float64{
			round(self.getInterpol(self.seqAt(0), t), roundOn, 0),
			round(self.getInterpol(self.seqAt(1), t), roundOn, 0),
		}
	} else if dimension == 3 {
		return []float64{
			round(self.getInterpol(self.seqAt(0), t), roundOn, 0),
			round(self.getInterpol(self.seqAt(1), t), roundOn, 0),
			round(self.getInterpol(self.seqAt(2), t), roundOn, 0),
		}
	} else {
		result := []float64{}
		for i := 0; i < dimension; i++ {
			result = append(result, round(self.getInterpol(self.seqAt(i), t), roundOn, 0))
		}
		return result
	}
}

// Round float number with precision. Working with negative numbers
func round(val float64, roundOn float64, places int ) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	_div := math.Copysign(div, val)
	_roundOn := math.Copysign(roundOn, val)
	if _div >= _roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}