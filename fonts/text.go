package fonts

import (
	"errors"
	"image/color"
	"log"
	"math"

	"github.com/gmlewis/go3d/float64/bezier2"
	"github.com/gmlewis/go3d/float64/qbezier2"
	"github.com/gmlewis/go3d/float64/vec2"
)

const (
	resolution = 0.1 // ems
	minSteps   = 4
	// If more resolution is needed in the rendered polygons,
	// MaxSteps could be increased.
	MaxSteps = 100

	// These are convenience constants for aligning text.
	XLeft   = 0
	XCenter = 0.5
	XRight  = 1
	YBottom = 0
	YCenter = 0.5
	YTop    = 1
)

var (
	// Convenience options for aligning the text.
	BottomLeft   = TextOpts{YAlign: YBottom, XAlign: XLeft}
	BottomCenter = TextOpts{YAlign: YBottom, XAlign: XCenter}
	BottomRight  = TextOpts{YAlign: YBottom, XAlign: XRight}
	CenterLeft   = TextOpts{YAlign: YCenter, XAlign: XLeft}
	Center       = TextOpts{YAlign: YCenter, XAlign: XCenter}
	CenterRight  = TextOpts{YAlign: YCenter, XAlign: XRight}
	TopLeft      = TextOpts{YAlign: YTop, XAlign: XLeft}
	TopCenter    = TextOpts{YAlign: YTop, XAlign: XCenter}
	TopRight     = TextOpts{YAlign: YTop, XAlign: XRight}
)

// TextOpts provides options for positioning (aligning) the text based on
// its minimum bounding box.
type TextOpts struct {
	// XAlign represents the horizontal alignment of the text.
	// 0=x origin at left (the default), 1=x origin at right, 0.5=center.
	// XLeft, XCenter, and XRight are defined for convenience and
	// readability of the code.
	XAlign float64
	// YAlign represents the vertical alignment of the text.
	// 0=y origin at bottom (the default), 1=y origin at top, 0.5=center.
	// YBottom, YCenter, and YTop are defined for convenience and
	// readbility of the code.
	YAlign float64
	// Rotate rotates the entire message about its anchor point
	// by this number of radians.
	Rotate float64
}

// Pt represents a 2D Point.
type Pt = vec2.T

// MBB represents a minimum bounding box.
type MBB = vec2.Rect

// Render represents a collection of polygons and includes
// the minimum bounding box of their union.
type Render struct {
	// MBB represents the minimum bounding box (MBB) of the render.
	MBB MBB
	// Polygons are the rendered polygons.
	Polygons []*Polygon
	// Info contains the MBB and base position of each glyph.
	// The length of info is identical to the number of runes in
	// the original text message.
	Info []*GlyphInfo
	// Background is the (optional) background color that the
	// "clear" polygons will use for rendering (default=white).
	Background color.Color
	// Foreground is the (optional) foreground color that the
	// "dark" polygons will use for rendering (default=black).
	Foreground color.Color
}

// GlyphInfo contains the MBB and base position of a glyph.
type GlyphInfo struct {
	// Glyph is the rendered rune from the original text message.
	Glyph rune
	// X, Y represent the base position of the glyph.
	X, Y float64
	// Width represents the width of the glyph.
	Width float64
	// MBB represents the minimum bounding box (MBB) of the glyph.
	MBB MBB
	// N represents the number of polygons dedicated to rendering this
	// glyph.
	N int
}

// Polygon represents a dark or clear polygon.
type Polygon struct {
	// RuneIndex is the index of the rune that this polygon belongs
	// to in the original rendered text message.
	RuneIndex int
	// Dark represents if this polygon is rendered dark (true) or clear (false).
	Dark bool
	// Pts is the collection of points making up the polygon.
	Pts []Pt
	// MBB represents the MBB of the polygon.
	MBB MBB
}

// Area calculates the area of the polygon.
func (p *Polygon) Area() float64 {
	return p.MBB.Area()
}

func getFont(fontName string) (*Font, error) {
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
	return font, nil
}

// TextMBB calculates the minimum bounding box of a text message without
// the overhead of rendering it.
func TextMBB(xPos, yPos, xScale, yScale float64, message, fontName string) (*MBB, error) {
	font, err := getFont(fontName)
	if err != nil {
		return nil, err
	}

	if font.HorizAdvX == 0 {
		font.HorizAdvX = font.UnitsPerEm
	}
	if font.MissingHorizAdvX == 0 {
		font.MissingHorizAdvX = font.HorizAdvX
	}

	var result *MBB
	x, y := 0.0, 0.0
	for _, c := range message {
		if c == rune('\n') {
			x, y = 0.0, y-(font.Ascent-font.Descent)
			continue
		}
		if c == rune('\t') {
			x += 2.0 * font.HorizAdvX
			continue
		}
		g, ok := font.Glyphs[c]
		if !ok {
			if c != ' ' {
				log.Printf("Warning: missing glyph %+q: skipping", c)
			}
			x += font.MissingHorizAdvX
			continue
		}
		width := (g.MBB.Max[0] - g.MBB.Min[0])
		height := (g.MBB.Max[1] - g.MBB.Min[1])
		minx := x + g.MBB.Min[0]
		miny := y + g.MBB.Min[1]
		mbb := MBB{
			Min: Pt{minx, miny},
			Max: Pt{minx + width, miny + height},
		}
		dx := g.HorizAdvX
		if dx == 0 {
			dx = font.HorizAdvX
		}
		x += dx
		// log.Printf("Glyph %+q: mbb=%v, x=%v", g.Unicode, mbb, x)
		if result == nil {
			result = &mbb
		} else {
			result.Join(&mbb)
		}
	}

	// log.Printf("TextMBB: xScale,yScale=(%v,%v)", xScale, yScale)
	fsf := 1.0 / font.UnitsPerEm
	xScale *= fsf
	yScale *= fsf
	// log.Printf("TextMBB: UnitsPerEm=%v, fsf=%v, scale=(%v,%v), xPos,yPos=(%v,%v)", font.UnitsPerEm, fsf, xScale, yScale, xPos, yPos)

	minx := xPos + xScale*result.Min[0]
	miny := yPos + yScale*result.Min[1]
	width := xScale * (result.Max[0] - result.Min[0])
	height := yScale * (result.Max[1] - result.Min[1])
	mbb := &MBB{
		Min: Pt{minx, miny},
		Max: Pt{minx + width, miny + height},
	}

	// log.Printf("TextMBB: message=%q, result=%v, mbb=%v", message, result, mbb)
	return mbb, nil
}

// Text returns a Render representing the rendered text.
// All dimensions are in "em"s, the width of the character "M" in the
// desired font.
//
// xScale and yScale are provided to convert the font to any scale desired.
func Text(xPos, yPos, xScale, yScale float64, message, fontName string, opts *TextOpts) (*Render, error) {
	font, err := getFont(fontName)
	if err != nil {
		return nil, err
	}

	mbb, err := TextMBB(xPos, yPos, xScale, yScale, message, fontName)
	if err != nil {
		return nil, err
	}

	if font.HorizAdvX == 0 {
		font.HorizAdvX = font.UnitsPerEm
	}
	if font.MissingHorizAdvX == 0 {
		font.MissingHorizAdvX = font.HorizAdvX
	}

	x, y := xPos, yPos
	fsf := 1.0 / font.UnitsPerEm
	xScale *= fsf
	yScale *= fsf
	var xAlign float64
	var yAlign float64
	if opts != nil {
		xAlign = opts.XAlign
		yAlign = opts.YAlign
	}

	width := (mbb.Max[0] - mbb.Min[0])
	height := (mbb.Max[1] - mbb.Min[1])
	xError := mbb.Min[0] - xPos
	yError := mbb.Min[1] - yPos
	xPos = xPos - xAlign*width - xError
	yPos = yPos - yAlign*height - yError
	x, y = xPos, yPos
	// log.Printf("Text: TextMBB=%v, Pos=(%v,%v), error=(%v,%v), x,y=(%v,%v)", mbb, xPos, yPos, xError, yError, x, y)

	var xformPt func(pt Pt) Pt
	if opts == nil || opts.Rotate == 0.0 {
		xformPt = func(pt Pt) Pt {
			return pt
		}
	} else {
		cos := math.Cos(opts.Rotate)
		sin := math.Sin(opts.Rotate)
		xformPt = func(pt Pt) Pt {
			dx := pt[0] - xPos
			dy := pt[1] - yPos
			return Pt{xPos + dx*cos - dy*sin, yPos + dy*cos + dx*sin}
		}
	}

	result := &Render{}
	for runeIndex, c := range message {
		newPt := xformPt(Pt{x, y})
		gi := &GlyphInfo{
			Glyph: c,
			X:     newPt[0],
			Y:     newPt[1],
			MBB:   MBB{Min: newPt, Max: newPt},
		}
		result.Info = append(result.Info, gi)

		if c == rune('\n') {
			x, y = xPos, y-yScale*(font.Ascent-font.Descent)
			continue
		}
		if c == rune('\t') {
			x += 2.0 * xScale * font.HorizAdvX
			continue
		}
		g, ok := font.Glyphs[c]
		if !ok {
			if c != ' ' {
				log.Printf("Warning: missing glyph %+q: skipping", c)
			}
			x += xScale * font.MissingHorizAdvX
			continue
		}
		dx, r := g.Render(x, y, xScale, yScale)
		gi.N = len(r.Polygons)
		for pi, poly := range r.Polygons {
			poly.RuneIndex = runeIndex
			for i, pt := range poly.Pts {
				newPt := xformPt(pt)
				v := MBB{Min: newPt, Max: newPt}
				if i == 0 {
					poly.MBB = v
				} else {
					poly.MBB.Join(&v)
				}
				poly.Pts[i] = newPt
			}
			if pi == 0 {
				r.MBB = poly.MBB
			} else {
				r.MBB.Join(&poly.MBB)
			}
			if len(result.Polygons) == 0 {
				result.MBB = r.MBB
			} else {
				result.MBB.Join(&r.MBB)
			}
			result.Polygons = append(result.Polygons, poly)
		}
		gi.MBB = r.MBB

		if dx == 0 {
			dx = font.HorizAdvX
		}
		width := dx * xScale
		gi.Width = width
		x += width
	}
	// log.Printf("FINAL MBB=%v", result.MBB)
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
			v := MBB{Min: pt, Max: pt}
			if i == 0 {
				r.MBB = v
			} else {
				r.MBB.Join(&v)
			}
		}
		if len(result.Polygons) == 0 {
			result.MBB = r.MBB
		} else {
			result.MBB.Join(&r.MBB)
		}

		result.Polygons = append(result.Polygons, &Polygon{
			Dark: currentPolarity == "d",
			Pts:  pts,
			MBB:  r.MBB,
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
			pts = []Pt{{x, y}}
		case 'm':
			if len(pts) > 0 {
				dumpPoly()
			}
			x, y = x+xScale*ps.P[0], y+yScale*ps.P[1]
			pts = []Pt{{x, y}}
		case 'L':
			for i := 0; i < len(ps.P); i += 2 {
				x, y = oX+xScale*ps.P[i], oY+yScale*ps.P[i+1]
				pts = append(pts, Pt{x, y})
			}
		case 'l':
			for i := 0; i < len(ps.P); i += 2 {
				x, y = x+xScale*ps.P[i], y+yScale*ps.P[i+1]
				pts = append(pts, Pt{x, y})
			}
		case 'H':
			for i := 0; i < len(ps.P); i++ {
				x = oX + xScale*ps.P[i]
				pts = append(pts, Pt{x, y})
			}
		case 'h':
			for i := 0; i < len(ps.P); i++ {
				x += xScale * ps.P[i]
				pts = append(pts, Pt{x, y})
			}
		case 'V':
			for i := 0; i < len(ps.P); i++ {
				y = oY + yScale*ps.P[i]
				pts = append(pts, Pt{x, y})
			}
		case 'v':
			for i := 0; i < len(ps.P); i++ {
				y += yScale * ps.P[i]
				pts = append(pts, Pt{x, y})
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
					pts = append(pts, p)
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
					pts = append(pts, p)
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
					pts = append(pts, p)
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
					pts = append(pts, p)
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
					pts = append(pts, p)
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
