// tri-bifilar-coil-diagram renders a diagram for the coil.
package main

import (
	"flag"
	"log"
	"math"

	"github.com/fogleman/gg"
	. "github.com/gmlewis/go-fonts/fonts"
	_ "github.com/gmlewis/go-fonts/fonts/latoregular"
)

var (
	width  = flag.Int("width", 800, "Image width")
	height = flag.Int("height", 800, "Image height")
	out    = flag.String("out", "tri-bifilar-coil-diagram.png", "PNG output filename")

	cx, cy float64
	outerR float64
)

const (
	textFont   = "latoregular"
	nlayers    = 3
	angleDelta = math.Pi / nlayers
	innerR     = 120
	segment    = 20
	maxA       = 33.0 * math.Pi
)

func main() {
	flag.Parse()

	innerHole := map[string]int{
		"TR": 0, "TL": 1, "BR": 0, "BL": 1,
		"2R": 2, "2L": 2,
	}
	innerHoleRev := map[int][]string{}
	for k, v := range innerHole {
		innerHoleRev[v] = append(innerHoleRev[v], k)
	}
	innerConnection := map[string]string{}
	for _, v := range innerHoleRev {
		if len(v) != 2 {
			log.Fatalf("v=%v", v)
		}
		innerConnection[v[0]] = v[1]
		innerConnection[v[1]] = v[0]
		// log.Printf("Inner: %v <=> %v", v[0], v[1])
	}

	outerHole := map[string]int{
		"TR": 0, "TL": 1, "BR": 1, "BL": 2,
		"2R": 3, "2L": 2,
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
		// log.Printf("Outer: %v <=> %v", v[0], v[1])
	}

	labels := []string{"TR"}
	for i := 0; i < 2*nlayers; i++ {
		last := labels[len(labels)-1]
		next := innerConnection[last]
		// log.Printf("A next: %v", next)
		labels = append(labels, next)
		if outerHole[next] == nlayers {
			break
		}
		next = outerConnection[next]
		// log.Printf("B next: %v", next)
		labels = append(labels, next)
	}

	cx = float64(*width) * 0.5
	cy = float64(*height) * 0.5
	outerR = float64(*width) * 0.25
	innerTS := 4.0
	outerTS := 6.0

	dc := gg.NewContext(*width, *height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	for n := 0; n < 2*nlayers; n++ {
		num := float64(n)
		if n < len(labels) {
			dc.Stroke()
			label := labels[n]
			// log.Printf("labels[%v]=%v", n, label)
			tp := innerPt(num, segment)
			text, err := Text(tp.X, tp.Y, innerTS, innerTS, label, textFont, &Center)
			check(err)
			text.RenderToDC(dc, tp.X-2*innerTS, tp.Y+2*innerTS, innerTS, 0)
			tp = outerPt(num, 1.5*segment)
			text, err = Text(tp.X, tp.Y, outerTS, outerTS, label, textFont, &Center)
			check(err)
			text.RenderToDC(dc, tp.X-2*outerTS, tp.Y+2*outerTS, outerTS, 0)
		}
	}
	dc.SetRGB(0, 0, 0)
	for n := 0; n < 2*nlayers; n++ {
		drawCoil(dc, n)
		num := float64(n)
		if n%2 == 0 {
			ip1 := innerPt(num, 0)
			ip2 := innerPt(num+1.0, 0)
			dc.MoveTo(ip1.X, ip1.Y)
			dc.LineTo(ip2.X, ip2.Y)
			mid1 := gg.Point{X: 0.5 * (ip1.X + ip2.X), Y: 0.5 * (ip1.Y + ip2.Y)}
			dc.Stroke()
			dc.DrawCircle(mid1.X, mid1.Y, 0.2*segment)
			dc.Fill()
		} else if n != 2*nlayers-1 {
			num := float64(n)
			op1 := outerPt(num, 0)
			op2 := outerPt(num+1.0, 0)
			dc.MoveTo(op1.X, op1.Y)
			dc.LineTo(op2.X, op2.Y)
			mid1 := gg.Point{X: 0.5 * (op1.X + op2.X), Y: 0.5 * (op1.Y + op2.Y)}
			dc.Stroke()
			dc.DrawCircle(mid1.X, mid1.Y, 0.5*segment)
			dc.Fill()
		} else {
			num := float64(n)
			op1 := outerPt(num, 0)
			op2 := outerPt(num+1.0, 0)
			dc.Stroke()
			dc.DrawCircle(op1.X, op1.Y, 0.5*segment)
			dc.DrawCircle(op2.X, op2.Y, 0.5*segment)
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

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
