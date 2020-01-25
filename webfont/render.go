package webfont

import (
	"log"
	"sort"
	"strings"

	"github.com/fogleman/gg"
	"github.com/gmlewis/go-fonts/fonts"
	"github.com/gmlewis/go3d/float64/bezier2"
	"github.com/gmlewis/go3d/float64/qbezier2"
	"github.com/gmlewis/go3d/float64/vec2"
)

type polyInfoT struct {
	starts int
	area   float64
}

// GenGerberLP renders a glyph to figure out the curve polarity
// and populate the GerberLP field. It also uses heuristics to
// determine the proper rendering order of the path subcommands.
func (g *Glyph) GenGerberLP(ff *FontFace) {
	if g == nil || g.Unicode == nil || len(g.PathSteps) == 0 {
		return
	}

	// First, get the bounding boxes of each polygon.
	gl := &fonts.Glyph{}
	var polyInfo []polyInfoT
	for i, ps := range g.PathSteps {
		switch ps.C {
		case "M", "m":
			polyInfo = append(polyInfo, polyInfoT{starts: i})
		}
		gl.PathSteps = append(gl.PathSteps, &fonts.PathStep{C: ps.C[0], P: ps.P})
	}
	_, render := gl.Render(0, 0, 1, 1)
	if len(polyInfo) > 1 {
		for i, poly := range render.Polygons {
			polyInfo[i].area = poly.Area()
		}
		g.ReorderByArea(polyInfo)
	}
	g.MBB = render.MBB
	// log.Printf("GenGerberLP: Glyph %+q: mbb=%v", *g.Unicode, g.MBB)

	width := int(0.5 + render.MBB.Max[0] - render.MBB.Min[0])
	height := int(0.5 + render.MBB.Max[1] - render.MBB.Min[1])

	oX, oY := -render.MBB.Min[0], -render.MBB.Min[1]
	x, y := oX, oY
	var lastC *bezier2.T
	var lastQ *qbezier2.T
	var lastCommand string

	var result []string
	dc := gg.NewContext(width, height)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	// for Processor:
	fillToX, fillToY := oX, oY
	for _, ps := range g.PathSteps {
		switch ps.C {
		case "M":
			x, y = oX+ps.P[0], oY+ps.P[1]
			dc.MoveTo(x, y)
			if g.Processor != nil {
				g.Processor.MoveTo(x, y)
			}
			fillToX, fillToY = x, y
			if len(result) == 0 {
				result = append(result, "d")
			} else {
				c := dc.Image().At(int(0.5+x), int(0.5+y))
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
			x, y = x+ps.P[0], y+ps.P[1]
			dc.MoveTo(x, y)
			if g.Processor != nil {
				g.Processor.MoveTo(x, y)
			}
			fillToX, fillToY = x, y
			if len(result) == 0 {
				result = append(result, "d")
			} else {
				c := dc.Image().At(int(0.5+x), int(0.5+y))
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
				x, y = oX+ps.P[i], oY+ps.P[i+1]
				dc.LineTo(x, y)
				if g.Processor != nil {
					g.Processor.LineTo(x, y)
				}
			}
		case "l":
			for i := 0; i < len(ps.P); i += 2 {
				x, y = x+ps.P[i], y+ps.P[i+1]
				dc.LineTo(x, y)
				if g.Processor != nil {
					g.Processor.LineTo(x, y)
				}
			}
		case "H":
			for i := 0; i < len(ps.P); i++ {
				x = oX + ps.P[i]
				dc.LineTo(x, y)
				if g.Processor != nil {
					g.Processor.LineTo(x, y)
				}
			}
		case "h":
			for i := 0; i < len(ps.P); i++ {
				x += ps.P[i]
				dc.LineTo(x, y)
				if g.Processor != nil {
					g.Processor.LineTo(x, y)
				}
			}
		case "V":
			for i := 0; i < len(ps.P); i++ {
				y = oY + ps.P[i]
				dc.LineTo(x, y)
				if g.Processor != nil {
					g.Processor.LineTo(x, y)
				}
			}
		case "v":
			for i := 0; i < len(ps.P); i++ {
				y += ps.P[i]
				dc.LineTo(x, y)
				if g.Processor != nil {
					g.Processor.LineTo(x, y)
				}
			}
		case "C":
			for i := 0; i < len(ps.P); i += 6 {
				x1, y1, x2, y2, ex, ey := oX+ps.P[i], oY+ps.P[i+1], oX+ps.P[i+2], oY+ps.P[i+3], oX+ps.P[i+4], oY+ps.P[i+5]
				b := &bezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x1, y1},
					P2: vec2.T{x2, y2},
					P3: vec2.T{ex, ey},
				}
				lastC = b
				dc.CubicTo(x1, y1, x2, y2, ex, ey)
				if g.Processor != nil {
					g.Processor.CubicTo(x1, y1, x2, y2, ex, ey)
				}
				x, y = ex, ey
			}
		case "c":
			for i := 0; i < len(ps.P); i += 6 {
				dx1, dy1, dx2, dy2, dx, dy := ps.P[i], ps.P[i+1], ps.P[i+2], ps.P[i+3], ps.P[i+4], ps.P[i+5]
				b := &bezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x + dx1, y + dy1},
					P2: vec2.T{x + dx2, y + dy2},
					P3: vec2.T{x + dx, y + dy},
				}
				lastC = b
				dc.CubicTo(x+dx1, y+dy1, x+dx2, y+dy2, x+dx, y+dy)
				if g.Processor != nil {
					g.Processor.CubicTo(x+dx1, y+dy1, x+dx2, y+dy2, x+dx, y+dy)
				}
				x, y = x+dx, y+dy
			}
		// case "S":
		case "s":
			for i := 0; i < len(ps.P); i += 4 {
				dx2, dy2, dx, dy := ps.P[i], ps.P[i+1], ps.P[i+2], ps.P[i+3]
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
				dc.CubicTo(x+dx1, y+dy1, x+dx2, y+dy2, x+dx, y+dy)
				if g.Processor != nil {
					g.Processor.CubicTo(x+dx1, y+dy1, x+dx2, y+dy2, x+dx, y+dy)
				}
				lastC = b
				x, y = x+dx, y+dy
			}
		// case "Q":
		case "q":
			for i := 0; i < len(ps.P); i += 4 {
				dx1, dy1, dx, dy := ps.P[i], ps.P[i+1], ps.P[i+2], ps.P[i+3]
				b := &qbezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x + dx1, y + dy1},
					P2: vec2.T{x + dx, y + dy},
				}
				dc.QuadraticTo(x+dx1, y+dy1, x+dx, y+dy)
				if g.Processor != nil {
					g.Processor.QuadraticTo(x+dx1, y+dy1, x+dx, y+dy)
				}
				lastQ = b
				x, y = x+dx, y+dy
			}
		// case "T":
		case "t":
			for i := 0; i < len(ps.P); i += 2 {
				dx, dy := ps.P[i], ps.P[i+1]
				dx1, dy1 := 0.0, 0.0
				if lastQ != nil && (lastCommand == "q" || lastCommand == "t") {
					dx1, dy1 = lastQ.P2[0]-lastQ.P1[0], lastQ.P2[1]-lastQ.P1[1]
				}
				lastQ = &qbezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x + dx1, y + dy1},
					P2: vec2.T{x + dx, y + dy},
				}
				dc.QuadraticTo(x+dx1, y+dy1, x+dx, y+dy)
				if g.Processor != nil {
					g.Processor.QuadraticTo(x+dx1, y+dy1, x+dx, y+dy)
				}
				x, y = x+dx, y+dy
			}
			// case "A":
			// case "a":
		case "Z", "z":
			if x != fillToX || y != fillToY {
				if g.Processor != nil {
					g.Processor.LineTo(fillToX, fillToY)
				}
			}
			dc.Fill()
		default:
			log.Fatalf("Unsupported path command %q in glyph %+q: %v", ps.C, *g.Unicode, *g.D)
		}
		lastCommand = ps.C
	}
	// if *g.Unicode == "j" {
	// 	dc.SavePNG(fmt.Sprintf("glyph-%v.png", *g.Unicode))
	// }

	s := strings.Join(result, "")
	g.GerberLP = &s
}

func (g *Glyph) ReorderByArea(polyInfo []polyInfoT) {
	var ascending bool
	for i, pi := range polyInfo {
		if i == 0 {
			continue
		}
		if pi.area > polyInfo[i-1].area {
			ascending = true
			break
		}
	}
	if !ascending { // Nothing to do.
		return
	}
	sort.Slice(polyInfo, func(a, b int) bool { return polyInfo[a].area > polyInfo[b].area })
	// log.Printf("reordering glyph %+q: polyInfo=%+v", *g.Unicode, polyInfo)

	var newPS []*PathStep
	for _, pi := range polyInfo {
		newPS = append(newPS, g.PathSteps[pi.starts])
		for i := pi.starts + 1; i < len(g.PathSteps); i++ {
			if g.PathSteps[i].C == "M" || g.PathSteps[i].C == "m" {
				break
			}
			newPS = append(newPS, g.PathSteps[i])
		}
	}
	if len(newPS) != len(g.PathSteps) {
		log.Fatalf("newPS(%v) != g.PathSteps(%v)", len(newPS), len(g.PathSteps))
	}
	g.PathSteps = newPS
}
