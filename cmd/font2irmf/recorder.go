package main

import (
	"fmt"
	"io"
	"log"
	"math"
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
	// spew.Fdump(f, g)
	fmt.Fprintf(f, `
float glyph_%v(float height, vec3 xyz) {
  if (any(lessThan(xyz, vec3(%.2f,%.2f,0.0))) || any(greaterThan(xyz, vec3(%.2f,%.2f,height)))) { return 0.0; }
  xyz -= vec3(%.2f,%.2f,0.0);
  float result = glyph_%v_1(xyz);
`, *g.Unicode, g.MBB.Min[0], g.MBB.Min[1], g.MBB.Max[0], g.MBB.Max[1], g.MBB.Min[0], g.MBB.Min[1], *g.Unicode)
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
	var mbb vec2.Rect
	for i, seg := range r.segments[segNum] {
		if i == 0 {
			mbb = vec2.Rect{vec2.T{seg.minX, seg.minY}, vec2.T{seg.maxX, seg.maxY}}
		} else {
			mbb.Join(&vec2.Rect{vec2.T{seg.minX, seg.minY}, vec2.T{seg.maxX, seg.maxY}})
		}
		yvals = append(yvals, seg.maxY, seg.minY)
	}
	sort.Slice(yvals, func(a, b int) bool {
		return yvals[b] < yvals[a]
	})

	fmt.Fprintf(f, `
float glyph_%v_%v(vec3 xyz) {
  if (any(lessThan(xyz.xy, vec2(%.2f,%.2f))) || any(greaterThan(xyz.xy, vec2(%.2f,%.2f)))) { return 0.0; }

`, glyphName, segNum+1, mbb.Min[0], mbb.Min[1], mbb.Max[0], mbb.Max[1])

	r.lastY = yvals[0]
	var sliceNum int
	for _, y := range yvals {
		if y >= r.lastY {
			continue
		}
		r.processSlice(f, r.lastY, y, sliceNum, segNum)
		r.lastY = y
		sliceNum++
	}
	fmt.Fprintln(f, `  return 1.0;
}`)
}

func (r *recorder) processSlice(f io.Writer, topY, botY float64, sliceNum, segNum int) {
	segs := r.getRange(topY, botY, segNum)
	op := "<"
	if sliceNum == 0 {
		op = "<="
	}
	switch len(segs) {
	case 1: // Only 1 segment falls within range?!?
		log.Fatalf("only 1 segment falls between y=%v and y=%v: %#v", botY, topY, segs[0])
	case 2:
		r.processTwoSegs(f, op, topY, botY, segs)
	case 4:
		r.processFourSegs(f, op, topY, botY, segs)
	default:
		log.Fatalf("%v segments not yet supported: botY=%v, topY=%v", len(segs), botY, topY)
	}
}

func (r *recorder) processTwoSegs(f io.Writer, op string, topY, botY float64, segs []*segment) {
	xs := [][]vec2.T{}
	// spew.Fdump(f, segs)
	xs = append(xs, segs[0].xIntercepts(topY, botY))
	xs = append(xs, segs[1].xIntercepts(topY, botY))
	left, right := segs[0], segs[1]
	if xs[0][0][0] <= xs[1][0][0] && xs[0][1][0] <= xs[1][1][0] {
		// left is left of right.
	} else if xs[0][0][0] >= xs[1][0][0] && xs[0][1][0] >= xs[1][1][0] {
		left, right = right, left // swap
	} else {
		log.Fatalf("two segments cross mid-y-slice: %v", spew.Sdump(segs))
	}

	fmt.Fprintf(f, "  if (xyz.y >= %0.2f && xyz.y %v %0.2f && (xyz.x < %v || xyz.x > %v)) { return 0.0; }\n",
		botY, op, topY, left.interpFunc(), right.interpFunc())
}

func (r *recorder) processFourSegs(f io.Writer, op string, topY, botY float64, segs []*segment) {
	sort.Slice(segs, func(a, b int) bool { // TODO: Precalculate the xIntercepts if slow.
		ax := segs[a].xIntercept(0.5 * (topY + botY))
		bx := segs[b].xIntercept(0.5 * (topY + botY))
		return ax < bx
	})

	fmt.Fprintf(f, "  if (xyz.y >= %0.2f && xyz.y %v %0.2f && (xyz.x < %v || (xyz.x > %v && xyz.x < %v) || xyz.x > %v)) { return 0.0; }\n",
		botY, op, topY, segs[0].interpFunc(), segs[1].interpFunc(), segs[2].interpFunc(), segs[3].interpFunc())
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

func (s *segment) xIntercept(y float64) float64 {
	switch s.segType {
	case line:
		return interpLine(s.pts[0][0], s.pts[1][0], y, s.pts[0][1], s.pts[1][1])
	// case cubic:
	case quadratic:
		return interpQuadratic(s.pts, y)
	default:
		log.Fatalf("Unknown segment type %v", s.segType)
	}
	return 0
}

func (s *segment) xIntercepts(topY, botY float64) []vec2.T {
	var result []vec2.T
	result = append(result, vec2.T{s.xIntercept(topY), topY})
	result = append(result, vec2.T{s.xIntercept(botY), botY})
	return result
}

func (s *segment) interpFunc() string {
	switch s.segType {
	case line:
		top, bot := s.pts[0], s.pts[1]
		if top[1] < bot[1] {
			top, bot = bot, top
		}
		return fmt.Sprintf("interpLine(vec2(%.2f,%.2f),vec2(%.2f,%.2f),xyz.y)", bot[0], bot[1], top[0], top[1])
	// case cubic:
	case quadratic:
		p0, p1, p2 := s.pts[0], s.pts[1], s.pts[2]
		if p2[1] < p0[1] {
			p0, p2 = p2, p0
		}
		return fmt.Sprintf("interpQuadratic(vec2(%.2f,%.2f),vec2(%.2f,%.2f),vec2(%.2f,%.2f),xyz.y)", p0[0], p0[1], p1[0], p1[1], p2[0], p2[1])
	default:
		log.Fatalf("Unknown segment type %v", s.segType)
	}
	return ""
}

func interpLine(to1, to2, val, from1, from2 float64) float64 {
	p := (val - from1) / (from2 - from1)
	return p*(to2-to1) + to1
}

func interpQuadratic(pts []vec2.T, y float64) float64 {
	a := pts[2][1] + pts[0][1] - 2*pts[1][1]
	b := 2 * (pts[1][1] - pts[0][1])
	c := pts[0][1] - y
	if b*b < 4*a*c {
		log.Printf("pts[0]=%#v, pts[1]=%#v, pts[2]=%#v", pts[0], pts[1], pts[2])
		log.Fatalf("bad quadratic equation: b^2=%v, 4ac=%v", b*b, 4*a*c)
	}
	det := math.Sqrt(b*b - 4*a*c)
	t := (-b + det) / (2 * a)
	if t2 := (-b - det) / (2 * a); t2 >= 0 && t2 <= 1 {
		t = t2
	}
	x := (1-t)*(1-t)*pts[0][0] + 2*(1-t)*t*pts[1][0] + t*t*pts[2][0]
	return x
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
	log.Printf("moveTo(%v,%v)", r.lastX, r.lastY)
	r.segments = append(r.segments, []*segment{})
}

func (r *recorder) lineTo(x, y float64) {
	log.Printf("from(%v,%v) - lineTo(%v,%v)", r.lastX, r.lastY, x, y)
	s := newSeg(line, []vec2.T{{r.lastX, r.lastY}, {x, y}})
	if s.minY != s.maxY {
		segNum := len(r.segments) - 1
		r.segments[segNum] = append(r.segments[segNum], s)
	}
	r.lastX, r.lastY = x, y
}

func (r *recorder) cubicTo(x1, y1, x2, y2, ex, ey float64) {
	log.Printf("from(%v,%v) - cubicTo((%v,%v),(%v,%v),(%v,%v))", r.lastX, r.lastY, x1, y1, x2, y2, ex, ey)
	s := newSeg(cubic, []vec2.T{{r.lastX, r.lastY}, {x1, y1}, {x2, y2}, {ex, ey}})
	if s.minY != s.maxY {
		segNum := len(r.segments) - 1
		r.segments[segNum] = append(r.segments[segNum], s)
	}
	r.lastX, r.lastY = ex, ey
}

func (r *recorder) quadraticTo(x1, y1, x2, y2 float64) {
	log.Printf("from(%v,%v) - quadraticTo((%v,%v),(%v,%v))", r.lastX, r.lastY, x1, y1, x2, y2)
	s := newSeg(quadratic, []vec2.T{{r.lastX, r.lastY}, {x1, y1}, {x2, y2}})
	if s.minY != s.maxY {
		segNum := len(r.segments) - 1
		r.segments[segNum] = append(r.segments[segNum], s)
	}
	r.lastX, r.lastY = x2, y2
}
