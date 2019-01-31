package fonts

import (
	"errors"
	"log"

	"github.com/gmlewis/go3d/float64/bezier2"
	"github.com/gmlewis/go3d/float64/qbezier2"
	"github.com/gmlewis/go3d/float64/vec2"
)

const (
	resolution = 0.01 // ems
	minSteps   = 4
	// If more resolution is needed in the rendered polygons,
	// MaxSteps could be increased.
	MaxSteps = 100
)

// Render represents a collection of polygons and includes
// the minimum bounding box of their union.
type Render struct {
	Xmin, Ymin float64
	Xmax, Ymax float64
	Polygons   []*Polygon
}

// Polygon represents a dark or clear polygon.
type Polygon struct {
	Dark bool
	Pts  []Pt
}

// Area calculates the area of the polygon by iterating
// over all its points. Therefore, don't call this
// function in a loop (such as a sort, for example).
func (p *Polygon) Area() float64 {
	var xmin, xmax, ymin, ymax float64
	for i, pt := range p.Pts {
		if i == 0 || pt.X < xmin {
			xmin = pt.X
		}
		if i == 0 || pt.X > xmax {
			xmax = pt.X
		}
		if i == 0 || pt.Y < ymin {
			ymin = pt.Y
		}
		if i == 0 || pt.Y > ymax {
			ymax = pt.Y
		}
	}
	return (xmax - xmin) * (ymax - ymin)
}

// Pt represents a 2D Point.
type Pt struct {
	X, Y float64
}

// Text returns a Render representing the rendered text.
// All dimensions are in "em"s, the width of the character "M" in the
// desired font.
//
// xScale and yScale are provided to convert the font to any scale desired.
func Text(xPos, yPos, xScale, yScale float64, message, fontName string) (*Render, error) {
	if len(Fonts) == 0 {
		return nil, errors.New("No fonts available")
	}

	font, ok := Fonts[fontName]
	if !ok {
		var name string
		for name, font = range Fonts { // Use the first (random) font found.
			break
		}
		log.Printf("Could not find font %q: using %q instead", fontName, name)
	}

	if font.HorizAdvX == 0 {
		font.HorizAdvX = font.UnitsPerEm
	}

	x, y := xPos, yPos
	fsf := 1.0 / font.UnitsPerEm
	xScale *= fsf
	yScale *= fsf

	result := &Render{}
	for _, c := range message {
		if c == rune('\n') {
			x, y = xPos, y-yScale*(font.Ascent-font.Descent)
			continue
		}
		if c == rune('\t') {
			x += 2.0 * xScale * font.HorizAdvX
			continue
		}
		g, ok := font.Glyphs[string(c)]
		if !ok {
			log.Printf("Warning: missing glyph %+q: skipping", c)
			x += xScale * font.HorizAdvX
			continue
		}
		dx, r := g.Render(x, y, xScale, yScale)
		if len(result.Polygons) == 0 || r.Xmin < result.Xmin {
			result.Xmin = r.Xmin
		}
		if len(result.Polygons) == 0 || r.Xmax > result.Xmax {
			result.Xmax = r.Xmax
		}
		if len(result.Polygons) == 0 || r.Ymin < result.Ymin {
			result.Ymin = r.Ymin
		}
		if len(result.Polygons) == 0 || r.Ymax > result.Ymax {
			result.Ymax = r.Ymax
		}
		result.Polygons = append(result.Polygons, r.Polygons...)
		if dx == 0 {
			dx = font.HorizAdvX
		}
		x += dx * xScale
	}
	return result, nil
}

// Render renders a glyph to polygons.
func (g *Glyph) Render(x, y, xScale, yScale float64) (float64, *Render) {
	oX, oY := x, y         // origin for this glyph
	var pts []Pt           // Current polygon
	currentPolarity := "d" // d=dark, c=clear
	var curveNum int

	result := &Render{}
	dumpPoly := func() {
		if g.GerberLP != "" && curveNum < len(g.GerberLP) {
			polarity := g.GerberLP[curveNum : curveNum+1]
			if polarity != currentPolarity {
				currentPolarity = polarity
			}
		}

		r := &Render{}
		for i, pt := range pts {
			if i == 0 || pt.X < r.Xmin {
				r.Xmin = pt.X
			}
			if i == 0 || pt.X > r.Xmax {
				r.Xmax = pt.X
			}
			if i == 0 || pt.Y < r.Ymin {
				r.Ymin = pt.Y
			}
			if i == 0 || pt.Y > r.Ymax {
				r.Ymax = pt.Y
			}
		}
		if len(result.Polygons) == 0 || r.Xmin < result.Xmin {
			result.Xmin = r.Xmin
		}
		if len(result.Polygons) == 0 || r.Xmax > result.Xmax {
			result.Xmax = r.Xmax
		}
		if len(result.Polygons) == 0 || r.Ymin < result.Ymin {
			result.Ymin = r.Ymin
		}
		if len(result.Polygons) == 0 || r.Ymax > result.Ymax {
			result.Ymax = r.Ymax
		}

		result.Polygons = append(result.Polygons, &Polygon{
			Dark: currentPolarity == "d",
			Pts:  pts,
		})
		pts = []Pt{}
	}

	var lastC *bezier2.T
	var lastQ *qbezier2.T
	var lastCommand byte
	for _, ps := range g.PathSteps {
		switch ps.C {
		case 'M':
			if len(pts) > 0 {
				dumpPoly()
			}
			x, y = oX+xScale*ps.P[0], oY+yScale*ps.P[1]
			pts = []Pt{{X: x, Y: y}}
		case 'm':
			if len(pts) > 0 {
				dumpPoly()
			}
			x, y = x+xScale*ps.P[0], y+yScale*ps.P[1]
			pts = []Pt{{X: x, Y: y}}
		case 'L':
			for i := 0; i < len(ps.P); i += 2 {
				x, y = oX+xScale*ps.P[i], oY+yScale*ps.P[i+1]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case 'l':
			for i := 0; i < len(ps.P); i += 2 {
				x, y = x+xScale*ps.P[i], y+yScale*ps.P[i+1]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case 'H':
			for i := 0; i < len(ps.P); i++ {
				x = oX + xScale*ps.P[i]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case 'h':
			for i := 0; i < len(ps.P); i++ {
				x += xScale * ps.P[i]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case 'V':
			for i := 0; i < len(ps.P); i++ {
				y = oY + yScale*ps.P[i]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case 'v':
			for i := 0; i < len(ps.P); i++ {
				y += yScale * ps.P[i]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case 'C':
			for i := 0; i < len(ps.P); i += 6 {
				x1, y1, x2, y2, ex, ey := oX+xScale*ps.P[i], oY+yScale*ps.P[i+1], oX+xScale*ps.P[i+2], oY+yScale*ps.P[i+3], oX+xScale*ps.P[i+4], oY+yScale*ps.P[i+5]
				b := &bezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x1, y1},
					P2: vec2.T{x2, y2},
					P3: vec2.T{ex, ey},
				}
				lastC = b
				length := b.Length(1)
				steps := int(0.5 + length/resolution)
				if steps < minSteps {
					steps = minSteps
				}
				if steps > MaxSteps {
					steps = MaxSteps
				}
				for j := 1; j <= steps; j++ {
					t := float64(j) / float64(steps)
					p := b.Point(t)
					pts = append(pts, Pt{X: p[0], Y: p[1]})
				}
				x, y = ex, ey
			}
		case 'c':
			for i := 0; i < len(ps.P); i += 6 {
				dx1, dy1, dx2, dy2, dx, dy := xScale*ps.P[i], yScale*ps.P[i+1], xScale*ps.P[i+2], yScale*ps.P[i+3], xScale*ps.P[i+4], yScale*ps.P[i+5]
				b := &bezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x + dx1, y + dy1},
					P2: vec2.T{x + dx2, y + dy2},
					P3: vec2.T{x + dx, y + dy},
				}
				lastC = b
				length := b.Length(1)
				steps := int(0.5 + length/resolution)
				if steps < minSteps {
					steps = minSteps
				}
				if steps > MaxSteps {
					steps = MaxSteps
				}
				for j := 1; j <= steps; j++ {
					t := float64(j) / float64(steps)
					p := b.Point(t)
					pts = append(pts, Pt{X: p[0], Y: p[1]})
				}
				x, y = x+dx, y+dy
			}
			// case 'S':
		case 's':
			for i := 0; i < len(ps.P); i += 4 {
				dx2, dy2, dx, dy := xScale*ps.P[i], yScale*ps.P[i+1], xScale*ps.P[i+2], yScale*ps.P[i+3]
				dx1, dy1 := 0.0, 0.0
				if lastC != nil && (lastCommand == 'c' || lastCommand == 's') {
					dx1, dy1 = lastC.P3[0]-lastC.P2[0], lastC.P3[1]-lastC.P2[1]
				}
				b := &bezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x + dx1, y + dy1},
					P2: vec2.T{x + dx2, y + dy2},
					P3: vec2.T{x + dx, y + dy},
				}
				lastC = b
				length := b.Length(1)
				steps := int(0.5 + length/resolution)
				if steps < minSteps {
					steps = minSteps
				}
				if steps > MaxSteps {
					steps = MaxSteps
				}
				for j := 1; j <= steps; j++ {
					t := float64(j) / float64(steps)
					p := b.Point(t)
					pts = append(pts, Pt{X: p[0], Y: p[1]})
				}
				x, y = x+dx, y+dy
			}
		// case 'Q':
		case 'q':
			for i := 0; i < len(ps.P); i += 4 {
				dx1, dy1, dx, dy := xScale*ps.P[i], yScale*ps.P[i+1], xScale*ps.P[i+2], yScale*ps.P[i+3]
				b := &qbezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x + dx1, y + dy1},
					P2: vec2.T{x + dx, y + dy},
				}
				lastQ = b
				length := b.Length(1)
				steps := int(0.5 + length/resolution)
				if steps < minSteps {
					steps = minSteps
				}
				if steps > MaxSteps {
					steps = MaxSteps
				}
				for j := 1; j <= steps; j++ {
					t := float64(j) / float64(steps)
					p := b.Point(t)
					pts = append(pts, Pt{X: p[0], Y: p[1]})
				}
				x, y = x+dx, y+dy
			}
		// case 'T':
		case 't':
			for i := 0; i < len(ps.P); i += 2 {
				dx, dy := xScale*ps.P[i], yScale*ps.P[i+1]
				dx1, dy1 := 0.0, 0.0
				if lastQ != nil && (lastCommand == 'q' || lastCommand == 't') {
					dx1, dy1 = lastQ.P2[0]-lastQ.P1[0], lastQ.P2[1]-lastQ.P1[1]
				}
				lastQ = &qbezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x + dx1, y + dy1},
					P2: vec2.T{x + dx, y + dy},
				}
				lastCommand = ps.C
				length := lastQ.Length(1)
				steps := int(0.5 + length/resolution)
				if steps < minSteps {
					steps = minSteps
				}
				if steps > MaxSteps {
					steps = MaxSteps
				}
				for j := 1; j <= steps; j++ {
					t := float64(j) / float64(steps)
					p := lastQ.Point(t)
					pts = append(pts, Pt{X: p[0], Y: p[1]})
				}
				x, y = x+dx, y+dy
			}
		// case 'A':
		// case 'a':
		case 'Z', 'z':
			if len(pts) > 0 {
				pts = append(pts, pts[0]) // Close the path.
				dumpPoly()
			}
			curveNum++
		default:
			log.Fatalf("Unsupported path command %q", ps.C)
		}
		lastCommand = ps.C
	}
	if len(pts) > 0 {
		dumpPoly()
	}

	return g.HorizAdvX, result
}
