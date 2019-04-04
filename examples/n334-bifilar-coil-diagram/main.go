// n334-bifilar-coil-diagram renders a diagram for the coil.
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
	width  = flag.Int("width", 5400, "Image width")
	height = flag.Int("height", 5400, "Image height")
	out    = flag.String("out", "n334-bifilar-coil-diagram.png", "PNG output filename")

	cx, cy float64
	outerR float64
)

const (
	textFont   = "latoregular"
	nlayers    = 334
	angleDelta = math.Pi / nlayers
	innerR     = 2400
	segment    = 20
	maxA       = 33.0 * math.Pi
)

func main() {
	flag.Parse()

	solution := findSolution()
	log.Printf("starting point: %v", solution)
	if solution == nil {
		log.Fatal("Unable to find acceptable starting point.")
	}

	labels := []string{solution.startLabel}
	for i := 0; i < 2*nlayers; i++ {
		last := labels[len(labels)-1]
		next := solution.innerConnection[last]
		// log.Printf("A next: %v", next)
		labels = append(labels, next)
		if solution.outerHole[next] == nlayers {
			break
		}
		next = solution.outerConnection[next]
		// log.Printf("B next: %v", next)
		labels = append(labels, next)
	}

	cx = float64(*width) * 0.5
	cy = float64(*height) * 0.5
	outerR = float64(*width) * 0.25
	innerTS := 3.0
	outerTS := 4.0

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
			tp := innerPt(num, 0.5*segment)
			text, err := Text(tp.X, tp.Y, innerTS, innerTS, label, textFont, &Center)
			check(err)
			text.RenderToDC(dc, tp.X-1.75*innerTS, tp.Y+1.5*innerTS, innerTS, 0)
			tp = outerPt(num, segment)
			text, err = Text(tp.X, tp.Y, outerTS, outerTS, label, textFont, &Center)
			check(err)
			text.RenderToDC(dc, tp.X-1.75*outerTS, tp.Y+1.5*outerTS, outerTS, 0)
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

type solution struct {
	innerTR         int
	outerTR         int
	innerHole       map[string]int
	outerHole       map[string]int
	innerConnection map[string]string
	outerConnection map[string]string
	startLabel      string
	endLabel        string
}

func findSolution() *solution {
	opposite, onAxisMirror := genMaps(nlayers, true)
	_, offAxisMirror := genMaps(nlayers, false)
	var result []string
	for innerTR := 0; innerTR < nlayers; innerTR++ {
		_, _, innerHole, innerConnection := wiring(innerTR, nlayers, opposite, onAxisMirror)
		for outerTR := 0; outerTR < nlayers; outerTR++ {
			startLabel, endLabel, outerHole, outerConnection := wiring(outerTR, nlayers, opposite, offAxisMirror)
			log.Printf("attempt: innerTR=%v, outerTR=%v", innerTR, outerTR)

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
				log.Printf("SUCCESS! %v: innerTR=%v, outerTR=%v, result: %v", len(result), innerTR, outerTR, strings.Join(result, " "))
				// if innerTR == 17 {
				fmt.Printf("\ninnerHole := map[string]int{\n")
				fmt.Printf("  %q: %v, %q: %v, %q: %v, %q: %v,\n",
					"TR", innerHole["TR"],
					"TL", innerHole["TL"],
					"BR", innerHole["BR"],
					"BL", innerHole["BL"],
				)
				for n := 2; n < nlayers; n += 2 {
					nr := fmt.Sprintf("%vR", n)
					nl := fmt.Sprintf("%vL", n)
					np1r := fmt.Sprintf("%vR", n+1)
					np1l := fmt.Sprintf("%vL", n+1)
					fmt.Printf("  %q: %v, %q: %v, %q: %v, %q: %v,\n",
						nr, innerHole[nr],
						nl, innerHole[nl],
						np1r, innerHole[np1r],
						np1l, innerHole[np1l],
					)
				}
				fmt.Printf("}\n\nouterHole := map[string]int{\n")
				fmt.Printf("  %q: %v, %q: %v, %q: %v, %q: %v,\n",
					"TR", outerHole["TR"],
					"TL", outerHole["TL"],
					"BR", outerHole["BR"],
					"BL", outerHole["BL"],
				)
				for n := 2; n < nlayers; n += 2 {
					nr := fmt.Sprintf("%vR", n)
					nl := fmt.Sprintf("%vL", n)
					np1r := fmt.Sprintf("%vR", n+1)
					np1l := fmt.Sprintf("%vL", n+1)
					fmt.Printf("  %q: %v, %q: %v, %q: %v, %q: %v,\n",
						nr, outerHole[nr],
						nl, outerHole[nl],
						np1r, outerHole[np1r],
						np1l, outerHole[np1l],
					)
				}
				fmt.Printf("}\n\n")
				return &solution{
					innerTR:         innerTR,
					outerTR:         outerTR,
					innerHole:       innerHole,
					outerHole:       outerHole,
					innerConnection: innerConnection,
					outerConnection: outerConnection,
					startLabel:      startLabel,
					endLabel:        endLabel,
				}
				// }
			}
		}
	}
	return nil
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
