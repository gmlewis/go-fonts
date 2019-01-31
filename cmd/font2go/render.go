package main

import (
	"log"
	"strings"

	"github.com/fogleman/gg"
	"github.com/gmlewis/go3d/float64/bezier2"
	"github.com/gmlewis/go3d/float64/qbezier2"
	"github.com/gmlewis/go3d/float64/vec2"
)

// GenGerberLP renders a glyph to figure out the curve polarity
// and populate the GerberLP field.
func (g *Glyph) GenGerberLP(ff *FontFace) {
	if g == nil || len(g.PathSteps) == 0 || (g.GerberLP != nil && len(*g.GerberLP) > 0) {
		return
	}

	offY := -ff.Descent
	x, y, oX, oY := 0.0, 0.0, 0.0, 0.0
	xScale := 1.0
	var lastC *bezier2.T
	var lastQ *qbezier2.T
	var lastCommand string

	var result []string
	dc := gg.NewContext(2048, 2048+int(offY))
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	for _, ps := range g.PathSteps {
		switch ps.C {
		case "M":
			x, y = oX+xScale*ps.P[0], oY+ps.P[1]
			dc.MoveTo(x, y+offY)
			if len(result) == 0 {
				result = append(result, "d")
			} else {
				c := dc.Image().At(int(0.5+x), int(0.5+y+offY))
				// log.Printf("c=%#v", c)
				r, _, _, _ := c.RGBA()
				if r == 0 {
					result = append(result, "d")
					dc.SetRGB(1, 1, 1)
				} else {
					result = append(result, "c")
					dc.SetRGB(0, 0, 0)
				}
			}
		case "m":
			x, y = x+xScale*ps.P[0], y+ps.P[1]
			dc.MoveTo(x, y+offY)
			if len(result) == 0 {
				result = append(result, "d")
			} else {
				c := dc.Image().At(int(0.5+x), int(0.5+y+offY))
				// log.Printf("c=%#v", c)
				r, _, _, _ := c.RGBA()
				if r == 0 {
					result = append(result, "d")
					dc.SetRGB(1, 1, 1)
				} else {
					result = append(result, "c")
					dc.SetRGB(0, 0, 0)
				}
			}
		case "L":
			for i := 0; i < len(ps.P); i += 2 {
				x, y = oX+xScale*ps.P[i], oY+ps.P[i+1]
				dc.LineTo(x, y+offY)
			}
		case "l":
			for i := 0; i < len(ps.P); i += 2 {
				x, y = x+xScale*ps.P[i], y+ps.P[i+1]
				dc.LineTo(x, y+offY)
			}
		case "H":
			for i := 0; i < len(ps.P); i++ {
				x = oX + xScale*ps.P[i]
				dc.LineTo(x, y+offY)
			}
		case "h":
			for i := 0; i < len(ps.P); i++ {
				x += xScale * ps.P[i]
				dc.LineTo(x, y+offY)
			}
		case "V":
			for i := 0; i < len(ps.P); i++ {
				y = oY + ps.P[i]
				dc.LineTo(x, y+offY)
			}
		case "v":
			for i := 0; i < len(ps.P); i++ {
				y += ps.P[i]
				dc.LineTo(x, y+offY)
			}
		case "C":
			for i := 0; i < len(ps.P); i += 6 {
				x1, y1, x2, y2, ex, ey := oX+xScale*ps.P[i], oY+ps.P[i+1], oX+xScale*ps.P[i+2], oY+ps.P[i+3], oX+xScale*ps.P[i+4], oY+ps.P[i+5]
				b := &bezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x1, y1},
					P2: vec2.T{x2, y2},
					P3: vec2.T{ex, ey},
				}
				lastC = b
				dc.CubicTo(x1, y1+offY, x2, y2+offY, ex, ey+offY)
				x, y = ex, ey
			}
		case "c":
			for i := 0; i < len(ps.P); i += 6 {
				dx1, dy1, dx2, dy2, dx, dy := xScale*ps.P[i], ps.P[i+1], xScale*ps.P[i+2], ps.P[i+3], xScale*ps.P[i+4], ps.P[i+5]
				b := &bezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x + dx1, y + dy1},
					P2: vec2.T{x + dx2, y + dy2},
					P3: vec2.T{x + dx, y + dy},
				}
				lastC = b
				dc.CubicTo(x+dx1, y+dy1+offY, x+dx2, y+dy2+offY, x+dx, y+dy+offY)
				x, y = x+dx, y+dy
			}
		// case "S":
		case "s":
			for i := 0; i < len(ps.P); i += 4 {
				dx2, dy2, dx, dy := xScale*ps.P[i], ps.P[i+1], xScale*ps.P[i+2], ps.P[i+3]
				dx1, dy1 := 0.0, 0.0
				if lastC != nil && (lastCommand == "c" || lastCommand == "s") {
					dx1, dy1 = lastC.P3[0]-lastC.P2[0], lastC.P3[1]-lastC.P2[1]
				}
				b := &bezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x + dx1, y + dy1},
					P2: vec2.T{x + dx2, y + dy2},
					P3: vec2.T{x + dx, y + dy},
				}
				dc.CubicTo(x+dx1, y+dy1+offY, x+dx2, y+dy2+offY, x+dx, y+dy+offY)
				lastC = b
				x, y = x+dx, y+dy
			}
		// case "Q":
		case "q":
			for i := 0; i < len(ps.P); i += 4 {
				dx1, dy1, dx, dy := xScale*ps.P[i], ps.P[i+1], xScale*ps.P[i+2], ps.P[i+3]
				b := &qbezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x + dx1, y + dy1},
					P2: vec2.T{x + dx, y + dy},
				}
				dc.QuadraticTo(x+dx1, y+dy1+offY, x+dx, y+dy+offY)
				lastQ = b
				x, y = x+dx, y+dy
			}
		// case "T":
		case "t":
			for i := 0; i < len(ps.P); i += 2 {
				dx, dy := xScale*ps.P[i], ps.P[i+1]
				dx1, dy1 := 0.0, 0.0
				if lastQ != nil && (lastCommand == "q" || lastCommand == "t") {
					dx1, dy1 = lastQ.P2[0]-lastQ.P1[0], lastQ.P2[1]-lastQ.P1[1]
				}
				lastQ = &qbezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x + dx1, y + dy1},
					P2: vec2.T{x + dx, y + dy},
				}
				dc.QuadraticTo(x+dx1, y+dy1+offY, x+dx, y+dy+offY)
				lastCommand = ps.C
				x, y = x+dx, y+dy
			}
			// case "A":
			// case "a":
		case "Z", "z":
			dc.Fill()
		default:
			log.Fatalf("Unsupported path command %q in glyph %+q: %v", ps.C, *g.Unicode, *g.D)
		}
	}
	s := strings.Join(result, "")
	g.GerberLP = &s
	log.Printf("Setting GerberLP=%q", *g.GerberLP)
}
