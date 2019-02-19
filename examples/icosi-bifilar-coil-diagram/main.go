// icosi-bifilar-coil-diagram renders a diagram for the coil.
package main

import (
	"flag"
	"log"
	"math"

	"github.com/fogleman/gg"
)

var (
	width  = flag.Int("width", 800, "Image width")
	height = flag.Int("height", 800, "Image height")
	out    = flag.String("out", "icosi-bifilar-coil-diagram.png", "PNG output filename")

	cx, cy float64
	outerR float64
)

const (
	nlayers    = 20
	angleDelta = math.Pi / nlayers
	innerR     = 120
	segment    = 20
	maxA       = 33.0 * math.Pi
)

func main() {
	flag.Parse()

	innerHole := map[string]int{
		"TR": 17, "TL": 7, "BR": 13, "BL": 3,
		"2R": 18, "2L": 8, "3R": 12, "3L": 2,
		"4R": 16, "4L": 6, "5R": 14, "5L": 4,
		"6R": 19, "6L": 9, "7R": 11, "7L": 1,
		"8R": 15, "8L": 5, "9R": 15, "9L": 5,
		"10R": 0, "10L": 10, "11R": 10, "11L": 0,
		"12R": 14, "12L": 4, "13R": 16, "13L": 6,
		"14R": 1, "14L": 11, "15R": 9, "15L": 19,
		"16R": 13, "16L": 3, "17R": 17, "17L": 7,
		"18R": 12, "18L": 2, "19R": 18, "19L": 8,
	}
	innerHoleRev := map[int][]string{}
	for k, v := range innerHole {
		innerHoleRev[v] = append(innerHoleRev[v], k)
	}
	innerConnection := map[string]string{}
	for _, v := range innerHoleRev {
		if len(v) != 2 {
			log.Fatalf("len(v)=%v, want 2", len(v))
		}
		innerConnection[v[0]] = v[1]
		innerConnection[v[1]] = v[0]
	}

	outerHole := map[string]int{
		"TR": 1, "TL": 11, "BR": 9, "BL": 19,
		"2R": 2, "2L": 12, "3R": 8, "3L": 18,
		"4R": 0, "4L": 10, "5R": 10, "5L": 20,
		"6R": 3, "6L": 13, "7R": 7, "7L": 17,
		"8R": 19, "8L": 9, "9R": 11, "9L": 1,
		"10R": 4, "10L": 14, "11R": 6, "11L": 16,
		"12R": 18, "12L": 8, "13R": 12, "13L": 2,
		"14R": 5, "14L": 15, "15R": 5, "15L": 15,
		"16R": 17, "16L": 7, "17R": 13, "17L": 3,
		"18R": 16, "18L": 6, "19R": 14, "19L": 4,
	}
	outerHoleRev := map[int][]string{}
	for k, v := range outerHole {
		outerHoleRev[v] = append(outerHoleRev[v], k)
	}
	outerConnection := map[string]string{}
	for _, v := range outerHoleRev {
		if len(v) != 2 {
			continue
		}
		outerConnection[v[0]] = v[1]
		outerConnection[v[1]] = v[0]
	}

	labels := []string{"4R"}
	for i := 0; i < 2*nlayers; i++ {
		last := labels[len(labels)-1]
		next := innerConnection[last]
		log.Printf("A next: %v", next)
		labels = append(labels, next)
		if outerHole[next] == 20 {
			break
		}
		next = outerConnection[next]
		log.Printf("B next: %v", next)
		labels = append(labels, next)
	}

	cx = float64(*width) * 0.5
	cy = float64(*height) * 0.5
	outerR = float64(*width) * 0.25

	dc := gg.NewContext(*width, *height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	for n := 0; n < 2*nlayers; n++ {
		drawCoil(dc, n)
		if n%2 == 0 {
			num := float64(n)
			ip1 := innerPt(num, 0)
			ip2 := innerPt(num+1.0, 0)
			dc.MoveTo(ip1.X, ip1.Y)
			dc.LineTo(ip2.X, ip2.Y)
			mid1 := innerPt(num+0.5, 0)
			dc.Stroke()
			dc.DrawCircle(mid1.X, mid1.Y, 0.2*segment)
			dc.Fill()
			if n < len(labels) {
				label := labels[n]
				tp := innerPt(num, segment)
				dc.DrawStringAnchored(label, tp.X, tp.Y, 0.5, 0.5)
				tp = outerPt(num, segment)
				dc.DrawStringAnchored(label, tp.X, tp.Y, 0.5, 0.5)
			}
		} else if n != 39 {
			num := float64(n)
			op1 := outerPt(num, 0)
			op2 := outerPt(num+1.0, 0)
			dc.MoveTo(op1.X, op1.Y)
			dc.LineTo(op2.X, op2.Y)
			mid1 := outerPt(num+0.5, 0)
			dc.Stroke()
			dc.DrawCircle(mid1.X, mid1.Y, 0.2*segment)
			dc.Fill()
			if n < len(labels) {
				label := labels[n]
				tp := innerPt(num, segment)
				dc.DrawStringAnchored(label, tp.X, tp.Y, 0.5, 0.5)
				tp = outerPt(num, segment)
				dc.DrawStringAnchored(label, tp.X, tp.Y, 0.5, 0.5)
			}
		} else {
			num := float64(n)
			op1 := outerPt(num, 0)
			op2 := outerPt(num+1.0, 0)
			dc.Stroke()
			dc.DrawCircle(op1.X, op1.Y, 0.2*segment)
			dc.DrawCircle(op2.X, op2.Y, 0.2*segment)
			dc.Fill()
		}
	}
	dc.Stroke()
	dc.SavePNG(*out)
}

func innerPt(num float64, dr float64) gg.Point {
	cos := math.Cos(num * angleDelta)
	sin := math.Sin(num * angleDelta)
	x1 := cx + (innerR-dr)*cos
	y1 := cy + (innerR-dr)*sin
	return gg.Point{X: x1, Y: y1}
}

func outerPt(num float64, dr float64) gg.Point {
	cos := math.Cos(num * angleDelta)
	sin := math.Sin(num * angleDelta)
	x2 := cx + (innerR+4*segment+maxA+dr)*cos
	y2 := cy + (innerR+4*segment+maxA+dr)*sin
	return gg.Point{X: x2, Y: y2}
}

func drawCoil(dc *gg.Context, n int) {
	num := float64(n)
	cos := math.Cos(num * angleDelta)
	sin := math.Sin(num * angleDelta)
	ip := innerPt(num, 0)
	x2 := cx + (innerR+segment)*cos
	y2 := cy + (innerR+segment)*sin
	dc.MoveTo(ip.X, ip.Y)
	dc.LineTo(x2, y2)
	theta := math.Pi
	for a := 0.1; a <= maxA; a += 0.1 {
		angle := num*angleDelta + theta + a
		x2 := cx + (innerR+2*segment+a)*cos
		y2 := cy + (innerR+2*segment+a)*sin
		x := x2 + (0.5*segment)*math.Cos(angle)
		y := y2 + (0.5*segment)*math.Sin(angle)
		dc.LineTo(x, y)
	}
	op := outerPt(num, 0)
	dc.LineTo(op.X, op.Y)
}
