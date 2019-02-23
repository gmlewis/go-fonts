// icosi-bifilar-coil-diagram renders a diagram for the coil.
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/fogleman/gg"
	. "github.com/gmlewis/go-fonts/fonts"
	_ "github.com/gmlewis/go-fonts/fonts/latoregular"
)

var (
	width  = flag.Int("width", 800, "Image width")
	height = flag.Int("height", 800, "Image height")
	out    = flag.String("out", "icosi-bifilar-coil-diagram.png", "PNG output filename")

	cx, cy float64
	outerR float64
)

const (
	textFont   = "latoregular"
	nlayers    = 6
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
		// log.Printf("Inner: %v <=> %v", v[0], v[1])
	}

	startingPoint := findStartingPoint()
	log.Printf("starting point: %v", startingPoint)
	if startingPoint == "" {
		log.Fatal("Unable to find acceptable starting point.")
	}

	outerHole := map[string]int{
		"TR": 18, "TL": 8, "BR": 12, "BL": 2,
		"2R": 19, "2L": 9, "3R": 11, "3L": 1,
		"4R": 17, "4L": 7, "5R": 13, "5L": 3,
		"6R": 0, "6L": 10, "7R": 10, "7L": 20,
		"8R": 16, "8L": 6, "9R": 14, "9L": 4,
		"10R": 1, "10L": 11, "11R": 9, "11L": 19,
		"12R": 15, "12L": 5, "13R": 15, "13L": 5,
		"14R": 2, "14L": 12, "15R": 8, "15L": 18,
		"16R": 14, "16L": 4, "17R": 16, "17L": 6,
		"18R": 13, "18L": 3, "19R": 17, "19L": 7,
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
		log.Printf("Outer: %v <=> %v", v[0], v[1])
	}

	labels := []string{"6R"}
	for i := 0; i < 2*nlayers; i++ {
		last := labels[len(labels)-1]
		next := innerConnection[last]
		log.Printf("A next: %v", next)
		labels = append(labels, next)
		if outerHole[next] == nlayers {
			break
		}
		next = outerConnection[next]
		log.Printf("B next: %v", next)
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
			log.Printf("labels[%v]=%v", n, label)
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
			dc.DrawCircle(mid1.X, mid1.Y, 0.2*segment)
			dc.Fill()
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

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func findStartingPoint() string {
	opposite, mirror := genMaps(nlayers, false)
	var result []string
	for innerTR := 0; innerTR < nlayers; innerTR++ {
		_, _, _, innerConnection := wiring(innerTR, nlayers, opposite, mirror)
		for startTR := 0; startTR < nlayers; startTR++ {
			startLabel, _, _, outerConnection := wiring(startTR, nlayers, opposite, mirror)
			log.Printf("attempt: innerTR=%v, startTR=%v", innerTR, startTR)

			result = []string{startLabel}
			seen := map[string]bool{startLabel: true}
			for i := 0; i < 2*nlayers; i++ {
				last := result[len(result)-1]
				next := innerConnection[last]
				if seen[next] {
					break
				}
				seen[next] = true
				result = append(result, next)
				next = outerConnection[next]
				if seen[next] {
					break
				}
				seen[next] = true
				result = append(result, next)
			}
			if len(result) == 2*nlayers {
				log.Fatalf("SUCCESS! %v: innerTR=%v, startTR=%v, result: %v", len(result), innerTR, startTR, strings.Join(result, " "))
			} else {
				log.Printf("failed: %v", len(result))
			}
		}
	}
	return ""
}

// wiring generates the startLabel (at point 0), the wiring map, and connection map
// for the given starting point of "TR" at position startTR.
func wiring(startTR, total int, opposite, mirror map[int]int) (string, string, map[string]int, map[string]string) {
	result := map[string]int{
		"TR": startTR,
		"TL": opposite[startTR],
		"BR": mirror[startTR],
		"BL": opposite[mirror[startTR]],
	}

	for n := 2; n < total; n += 2 {
		var offset int
		if (n+2)%4 == 0 {
			offset = (n + 2) / 4
		} else {
			offset = -(n / 4)
		}
		nr := (startTR + offset + total) % total
		result[fmt.Sprintf("%vR", n)] = nr
		result[fmt.Sprintf("%vL", n)] = opposite[nr]
		result[fmt.Sprintf("%vR", n+1)] = mirror[nr]
		result[fmt.Sprintf("%vL", n+1)] = opposite[mirror[nr]]
	}

	startLabel, endLabel := "", ""
	rev := map[int][]string{}
	for k, v := range result {
		if v == 0 {
			if strings.HasSuffix(k, "R") {
				startLabel = k
			} else {
				endLabel = k
			}
		}
		rev[v] = append(rev[v], k)
	}
	connection := map[string]string{}
	for _, v := range rev {
		if len(v) != 2 {
			continue
		}
		connection[v[0]] = v[1]
		connection[v[1]] = v[0]
	}

	return startLabel, endLabel, result, connection
}

func genMaps(total int, onAxis bool) (opposite, mirror map[int]int) {
	opposite, mirror = map[int]int{}, map[int]int{}

	for n := 0; n < total/2; n++ {
		o := n + total/2
		opposite[n] = o
		opposite[o] = n
	}
	if onAxis {
		for n := 0; n <= total/4; n++ {
			o := total/2 - n
			mirror[n] = o
			mirror[o] = n
			if n > 0 {
				m := total - n
				o = total/2 + n
				mirror[m] = o
				mirror[o] = m
			}
		}
	} else {
		for n := 0; n <= total/4+1; n++ {
			o := total/2 - n - 1
			mirror[n] = o
			mirror[o] = n
			m := total - n - 1
			o = total/2 + n
			mirror[m] = o
			mirror[o] = m
		}
	}

	return opposite, mirror
}
