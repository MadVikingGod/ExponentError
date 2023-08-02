package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"text/tabwriter"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
)

func main() {
	f, err := os.Create("output.csv")
	if err != nil {
		fmt.Println("Could not open file output.csv")
		return
	}
	defer f.Close()
	c := csv.NewWriter(f)

	c.Write([]string{"value", "scale", "MapToIndexScale0", "MapToIndexFrexp", "MapToIndex"})

	cfg := config{
		start:    0x1.fffffffffffe1p34,
		stop:     0x1.0000000000010p35,
		title:    "MapToIndex for Values near 2^35",
		FileName: "MapToIndex_35.png",
	}

	datapoints, err := iterate(cfg)
	if err != nil {
		fmt.Println("iteration Failed")
	}
	c.WriteAll(datapoints)

	cfg = config{
		start:    0x1.fffffffffff00p299,
		stop:     0x1.0000000000010p300,
		title:    "MapToIndex for Values near 2^300",
		FileName: "MapToIndex_300.png",
	}

	datapoints, err = iterate(cfg)
	if err != nil {
		fmt.Println("iteration Failed")
	}
	c.WriteAll(datapoints)

	cfg = config{
		start:    0x1.ffffffffffff0p49,
		stop:     0x1.0000000000010p50,
		title:    "MapToIndex for Values near 2^50",
		FileName: "MapToIndex_50.png",
	}

	datapoints, err = iterate(cfg)
	if err != nil {
		fmt.Println("iteration Failed")
	}
	c.WriteAll(datapoints)

	cfg = config{
		start:    0x1.aaaaec907a705p+300,
		stop:     0x1.aaaaec907a725p+300,
		title:    "MapToIndex for Values near 2^300",
		FileName: "MapToIndex_mid_300.png",
		scale:    15,
		debug:    true,
	}

	datapoints, err = iterate(cfg)
	if err != nil {
		fmt.Println("iteration Failed")
	}
	c.WriteAll(datapoints)

	fmt.Printf("%x\n", 3.3950679616459195429964713197e90)

	genError()
}

type config struct {
	start, stop float64
	title       string
	FileName    string
	scale       int
	debug       bool
}

func iterate(cfg config) ([][]string, error) {
	p := plot.New()
	p.Title.Text = cfg.title
	p.X.Tick.Label.Rotation = math.Pi / 2
	p.X.Tick.Label.XAlign = text.XRight
	p.X.Padding = 1 * vg.Centimeter
	p.X.Label.Text = "Value"
	p.Y.Label.Text = "Bucket Index"
	xys := plotter.XYs{}
	labels := []string{}
	x := 0.0

	if cfg.start > cfg.stop {
		cfg.start, cfg.stop = cfg.stop, cfg.start
	}

	dataPoints := [][]string{}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	for v := cfg.start; v < cfg.stop; v = math.Nextafter(v, math.MaxFloat64) {
		mtis := MapToIndexScale0(v)
		mtif := MapToIndexFrexp(v)
		mti := MapToIndex(v, cfg.scale)
		dataPoints = append(dataPoints, []string{fmt.Sprintf("%x", v), fmt.Sprintf("%d", cfg.scale), fmt.Sprintf("%v", mtis), fmt.Sprintf("%v", mtif), fmt.Sprintf("%v", mti)})

		if cfg.debug {
			fmt.Fprintf(w, "%x\t%v\t%v\t%v\n", v, mti, MapExperimental(v, cfg.scale), MapExperimental2(v, cfg.scale))
		}
		xys = append(xys, plotter.XY{X: x, Y: float64(mti)})
		labels = append(labels, fmt.Sprintf("%x", v))
		x += 1.0
	}

	w.Flush()
	line, scatter, err := plotter.NewLinePoints(xys)
	if err != nil {
		fmt.Println("Could not create labels")
		return dataPoints, err
	}
	p.NominalX(labels...)
	p.Add(line)
	p.Add(scatter)
	err = p.Save(30*vg.Centimeter, 10*vg.Centimeter, cfg.FileName)

	return dataPoints, err
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

func genError() {
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
