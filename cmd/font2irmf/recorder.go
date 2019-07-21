package main

import (
	"fmt"
	"io"
	"log"
	"sort"

	"github.com/davecgh/go-spew/spew"
	"github.com/gmlewis/go3d/float64/vec2"
)

type recorder struct {
	lastX, lastY float64

	segments [][]*segment
}

func (r *recorder) process(f io.Writer, g *Glyph) {
	for i := range r.segments {
		r.processGerberLP(f, *g.Unicode, *g.GerberLP, i)
	}
	fmt.Fprintf(f, `
float glyph_%v(float height, vec3 xyz) {
  if (any(lessThan(xyz, vec3(%.2f,%.2f,0.0))) || any(greaterThan(xyz, vec3(%.2f,%.2f,height)))) { return 0.0; }
  float result = glyph_%v_1(xyz);
`, *g.Unicode, g.MBB.Min[0], g.MBB.Min[1], g.MBB.Max[0], g.MBB.Max[1], *g.Unicode)
	for i := range r.segments[1:] {
		op := "+"
		if (*g.GerberLP)[i+1:i+2] == "c" {
			op = "-"
		}
		fmt.Fprintf(f, "  result %v= glyph_%v_%v(xyz);\n", op, *g.Unicode, i+2)
	}
	fmt.Fprintln(f, `  return result;
}`)
}

func (r *recorder) processGerberLP(f io.Writer, glyphName string, gerberLP string, segNum int) {
	// Sort all Y values, descending.
	var yvals []float64
	for _, seg := range r.segments[segNum] {
		yvals = append(yvals, seg.maxY, seg.minY)
	}
	sort.Slice(yvals, func(a, b int) bool {
		return yvals[b] < yvals[a]
	})

	fmt.Fprintf(f, `
float glyph_%v_%v(vec3 xyz) {
`, glyphName, segNum+1)

	r.lastY = yvals[0]
	for _, y := range yvals {
		if y >= r.lastY {
			continue
		}
		r.processSlice(f, r.lastY, y, segNum)
		r.lastY = y
	}
	fmt.Fprintln(f, `  return 1.0;
}`)
}

func (r *recorder) processSlice(f io.Writer, topY, botY float64, segNum int) {
	segs := r.getRange(topY, botY, segNum)
	switch len(segs) {
	case 1: // Only 1 segment falls within range?!?
		log.Fatalf("only 1 segment falls between y=%v and y=%v: %#v", botY, topY, segs[0])
	default:
		fmt.Fprintf(f, "  if (xyz.y <= %0.2f && xyz.y >= %0.2f) { return 0.0; } // %v segs\n", topY, botY, len(segs))
		spew.Fdump(f, segs)
	}
}

func (r *recorder) getRange(topY, botY float64, segNum int) []*segment {
	var result []*segment
	for _, seg := range r.segments[segNum] {
		if seg.minY <= botY && seg.maxY >= topY {
			result = append(result, seg)
		}
	}
	return result
}

type segmentType uint8

const (
	line segmentType = iota
	cubic
	quadratic
)

func segTypeName(segType segmentType) string {
	switch segType {
	case line:
		return "line"
	case cubic:
		return "cubic"
	case quadratic:
		return "quadratic"
	}
	return "unknown"
}

type segment struct {
	segType    segmentType
	pts        []vec2.T
	minX, maxX float64
	minY, maxY float64
}

func newSeg(segType segmentType, pts []vec2.T) *segment {
	s := &segment{
		segType: segType,
		pts:     pts,
	}
	for i, pt := range pts {
		if i == 0 || pt[0] < s.minX {
			s.minX = pt[0]
		}
		if i == 0 || pt[0] > s.maxX {
			s.maxX = pt[0]
		}
		if i == 0 || pt[1] < s.minY {
			s.minY = pt[1]
		}
		if i == 0 || pt[1] > s.maxY {
			s.maxY = pt[1]
		}
	}
	return s
}

func (r *recorder) moveTo(x, y float64) {
	r.lastX, r.lastY = x, y
	r.segments = append(r.segments, []*segment{})
}

func (r *recorder) lineTo(x, y float64) {
	s := newSeg(line, []vec2.T{{r.lastX, r.lastY}, {x, y}})
	if s.minY != s.maxY {
		segNum := len(r.segments) - 1
		r.segments[segNum] = append(r.segments[segNum], s)
	}
	r.lastX, r.lastY = x, y
}

func (r *recorder) cubicTo(x1, y1, x2, y2, ex, ey float64) {
	s := newSeg(cubic, []vec2.T{{r.lastX, r.lastY}, {x1, y1}, {x2, y2}, {ex, ey}})
	if s.minY != s.maxY {
		segNum := len(r.segments) - 1
		r.segments[segNum] = append(r.segments[segNum], s)
	}
	r.lastX, r.lastY = ex, ey
}

func (r *recorder) quadraticTo(x1, y1, x2, y2 float64) {
	s := newSeg(quadratic, []vec2.T{{r.lastX, r.lastY}, {x1, y1}, {x2, y2}})
	if s.minY != s.maxY {
		segNum := len(r.segments) - 1
		r.segments[segNum] = append(r.segments[segNum], s)
	}
	r.lastX, r.lastY = x2, y2
}
