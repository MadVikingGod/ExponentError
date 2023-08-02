package main

import (
	"fmt"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {
	scale := 20
	pts := errorPts(scale, MapToIndex)
	pts2 := errorPts(scale, MapExperimental2)

	p := plot.New()
	p.Title.Text = fmt.Sprintf("Error in MapToIndex, Scale %d", scale)
	p.X.Label.Text = "Power of 2"
	p.X.Scale = plot.LinearScale{}
	p.Y.Label.Text = "Count of values with error"

	line, scatter, _ := plotter.NewLinePoints(pts)
	line.Color = plotutil.Color(1)
	p.Add(line, scatter, plotter.NewGrid())

	line2, scatter2, _ := plotter.NewLinePoints(pts2)
	line2.Color = plotutil.Color(2)
	scatter2.GlyphStyle.Shape = plotutil.Shape(2)
	p.Add(line2, scatter2, plotter.NewGrid())

	p.Legend.Add("Current", line, scatter)
	p.Legend.Add("Proposed", line2, scatter2)
	p.Legend.Top = true

	p.Save(120*vg.Centimeter, 20*vg.Centimeter, "MapToIndex_Error.svg")
}

// MapToIndex for any scale.
func MapToIndex(value float64, scale int) int {
	// Special case for power-of-two values.
	scaleFactor := math.Ldexp(math.Log2E, scale)
	// Note: math.Floor(value) equals math.Ceil(value)-1 when value
	// is not a power of two, which is checked above.
	return int(math.Floor(math.Log(value) * scaleFactor))
}

func MapExperimental(value float64, scale int) int {
	// This splits the value into a fraction and exponent, but because of the
	// choice of fraction the exp is 1 higher then we want.
	frac, exp := math.Frexp(value)

	return exp<<scale + int(math.Log2(frac)*float64(int(1)<<scale)) - 1
}
func MapExperimental2(value float64, scale int) int {
	frac, exp := math.Frexp(value)

	return exp<<scale + int(math.Log(frac)*scaleFactors[scale]) - 1
}

var scaleFactors = [21]float64{
	math.Ldexp(math.Log2E, 0),
	math.Ldexp(math.Log2E, 1),
	math.Ldexp(math.Log2E, 2),
	math.Ldexp(math.Log2E, 3),
	math.Ldexp(math.Log2E, 4),
	math.Ldexp(math.Log2E, 5),
	math.Ldexp(math.Log2E, 6),
	math.Ldexp(math.Log2E, 7),
	math.Ldexp(math.Log2E, 8),
	math.Ldexp(math.Log2E, 9),
	math.Ldexp(math.Log2E, 10),
	math.Ldexp(math.Log2E, 11),
	math.Ldexp(math.Log2E, 12),
	math.Ldexp(math.Log2E, 13),
	math.Ldexp(math.Log2E, 14),
	math.Ldexp(math.Log2E, 15),
	math.Ldexp(math.Log2E, 16),
	math.Ldexp(math.Log2E, 17),
	math.Ldexp(math.Log2E, 18),
	math.Ldexp(math.Log2E, 19),
	math.Ldexp(math.Log2E, 20),
}

func errorPts(scale int, mapF func(float64, int) int) plotter.XYs {
	pts := plotter.XYs{}

	for x := -1022; x <= 1022; x++ {
		i := 0
		v := math.Exp2(float64(x))
		for !inBound(v, scale, mapF) {
			i++
			v = math.Nextafter(v, 0)
		}
		pts = append(pts, plotter.XY{X: float64(x), Y: float64(i)})
	}
	return pts
}

func inBound(value float64, scale int, mapF func(float64, int) int) bool {
	idx := mapF(value, scale)
	return value > lowerBound(idx, scale) && value <= lowerBound(idx+1, scale)
}

func lowerBound(index int, scale int) float64 {
	// The lowerBound of the index of Math.SmallestNonzeroFloat64 at any scale
	// is always rounded down to 0.0.
	// For example lowerBound(getBin(Math.SmallestNonzeroFloat64, 7), 7) == 0.0
	// 2 ^ (index * 2 ^ (-scale))
	return math.Exp2(math.Ldexp(float64(index), -scale))
}
