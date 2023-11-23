// -*- compile-command: "go run main.go ../../fonts/freesans/FreeSans.svg"; -*-

// font2lua reads one or more standard SVG webfont file(s) and writes Lua file(s)
// used to render them to polygons in Blackjack and go-bjk.
// See: https://github.com/setzer22/blackjack
// and: https://github.com/gmlewis/go-bjk
//
// Usage:
//
//	font2lua
package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/gmlewis/go-fonts/webfont"
)

const (
	prefix = "fonts"
)

var (
	debug = flag.Bool("debug", false, "Turn on debugging info")

	outTemp = template.Must(template.New("out").Funcs(funcMap).Parse(luaTemplate))
	funcMap = template.FuncMap{
		"floats":     floats,
		"orEmpty":    orEmpty,
		"viewFilter": viewFilter,
	}

	digitRE = regexp.MustCompile(`^\d`)
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		logf("Processing file %q ...", arg)

		fontData := &webfont.FontData{}
		if buf, err := os.ReadFile(arg); err != nil {
			log.Fatal(err)
		} else {
			if err := xml.Unmarshal(buf, fontData); err != nil {
				log.Fatal(err)
			}
		}

		fontData.Font.ID = strings.ToLower(fontData.Font.ID)
		fontData.Font.ID = strings.Replace(fontData.Font.ID, "-", "_", -1)
		if digitRE.MatchString(fontData.Font.ID) {
			fontData.Font.ID = "f" + fontData.Font.ID
		}

		writeFont(fontData)
	}

	fmt.Println("Done.")
}

// processor implements the webfont.Processor interface.
type processor struct {
	current *glyphT
	glyphs  map[string]*glyphT
}

var _ webfont.Processor = &processor{}

type glyphT struct {
	horizAdvX float64
	unicode   string
	gerberLP  string
	d         string
	dOrig     string
	xmin      float64
	ymin      float64
	xmax      float64
	ymax      float64

	faces []*faceT
}

type vec2 [2]float64

type faceT struct {
	absCmds []string
	after   []vec2

	cut0Idx int // face[0] vertIdx (in 'after') to cut to connect to this face
	cutIdx  int // vertIdx (in 'after') to cut to connect to face[0]
}

func (p *processor) NewGlyph(g *webfont.Glyph) {
	var d string
	if g.D != nil {
		d = *g.D
	}
	dOrig := d
	if g.DOrig != nil {
		dOrig = *g.DOrig
	}

	glyph := &glyphT{
		horizAdvX: g.HorizAdvX,
		unicode:   *g.Unicode,
		d:         d,
		dOrig:     dOrig,
		xmin:      g.MBB.Min[0],
		ymin:      g.MBB.Min[1],
		xmax:      g.MBB.Max[0],
		ymax:      g.MBB.Max[1],
	}

	p.glyphs[*g.Unicode] = glyph
	p.current = glyph
}

func (p *processor) ProcessGlyph(r rune, g *webfont.Glyph) {
	glyph := p.current
	if g.GerberLP != nil {
		glyph.gerberLP = *g.GerberLP
	}
	if glyph.gerberLP != "" {
		logf("glyph.gerberLP=%q", glyph.gerberLP)
		glyph.regenerateFace()
	}
}

func (p *processor) addCmd(glyph *glyphT, oldX, oldY float64, cmd string, x, y float64, absCmd string) {
	if cmd == "M" { // start a new face
		glyph.faces = append(glyph.faces, &faceT{
			absCmds: []string{absCmd},
			after:   []vec2{{x + glyph.xmin, y + glyph.ymin}},
		})
		return
	}

	face := len(glyph.faces) - 1
	glyph.faces[face].absCmds = append(glyph.faces[face].absCmds, absCmd)
	glyph.faces[face].after = append(glyph.faces[face].after,
		vec2{x + glyph.xmin, y + glyph.ymin})
}

func (p *processor) MoveTo(g *webfont.Glyph, oldX, oldY float64, cmd string, x, y float64) {
	glyph := p.current
	logf("p.MoveTo(g,%v,%v,%q,%v,%v)", oldX+glyph.xmin, oldY+glyph.ymin, cmd, x+glyph.xmin, y+glyph.ymin)
	absCmd := fmt.Sprintf("M%v %v", x+glyph.xmin, y+glyph.ymin)
	p.addCmd(glyph, oldX, oldY, cmd, x, y, absCmd)
}

func (p *processor) LineTo(g *webfont.Glyph, oldX, oldY float64, cmd string, x, y float64) {
	glyph := p.current
	logf("p.LineTo(g,%v,%v,%q,%v,%v)", oldX+glyph.xmin, oldY+glyph.ymin, cmd, x+glyph.xmin, y+glyph.ymin)
	absCmd := fmt.Sprintf("L%v %v", x+glyph.xmin, y+glyph.ymin)
	p.addCmd(glyph, oldX, oldY, cmd, x, y, absCmd)
}

func (p *processor) CubicTo(g *webfont.Glyph, oldX, oldY float64, cmd string, x1, y1, x2, y2, ex, ey float64) {
	glyph := p.current
	logf("p.CubicTo(g,%v,%v,%q,%v,%v,%v,%v,%v,%v)", oldX+glyph.xmin, oldY+glyph.ymin, cmd, x1+glyph.xmin, y1+glyph.ymin, x2+glyph.xmin, y2+glyph.ymin, ex+glyph.xmin, ey+glyph.ymin)
	absCmd := fmt.Sprintf("C%v %v %v %v %v %v", x1+glyph.xmin, y1+glyph.ymin, x2+glyph.xmin, y2+glyph.ymin, ex+glyph.xmin, ey+glyph.ymin)
	p.addCmd(glyph, oldX, oldY, cmd, ex, ey, absCmd)
}

func (p *processor) QuadraticTo(g *webfont.Glyph, oldX, oldY float64, cmd string, x1, y1, x2, y2 float64) {
	glyph := p.current
	logf("p.QuadraticTo(g,%v,%v,%q,%v,%v,%v,%v)", oldX+glyph.xmin, oldY+glyph.ymin, cmd, x1+glyph.xmin, y1+glyph.ymin, x2+glyph.xmin, y2+glyph.ymin)
	absCmd := fmt.Sprintf("Q%v %v %v %v", x1+glyph.xmin, y1+glyph.ymin, x2+glyph.xmin, y2+glyph.ymin)
	p.addCmd(glyph, oldX, oldY, cmd, x2, y2, absCmd)
}

func (g *glyphT) findClosestVerts(face *faceT) {
	bestDistSq := math.MaxFloat64
	for f0vertIdx, f0vertAfter := range g.faces[0].after {
		for fivertIdx, fivertAfter := range face.after {
			dx := fivertAfter[0] - f0vertAfter[0]
			dy := fivertAfter[1] - f0vertAfter[1]
			distSq := dx*dx + dy*dy
			if distSq < bestDistSq {
				logf("c=%q: fivertAfter[%v]=%v, f0vertAfter[%v]=%v, distSq=%v", g.unicode, fivertIdx, fivertAfter, f0vertIdx, f0vertAfter, distSq)
				bestDistSq = distSq
				face.cut0Idx = f0vertIdx
				face.cutIdx = fivertIdx
			}
		}
	}
}

func (g *glyphT) findCutPoints() {
	logf("c=%q: d=%v", g.unicode, g.d)
	logf("c=%q: got %v faces", g.unicode, len(g.faces))

	for faceIdx, face := range g.faces {
		logf("c=%q: face[%v]: absCmds(%v): %#v", g.unicode, faceIdx, len(face.absCmds), face.absCmds)
		logf("c=%q: face[%v]: after(%v): %+v", g.unicode, faceIdx, len(face.after), face.after)
		if len(face.after) != len(face.absCmds) {
			log.Fatalf("programming error - absCmds not split correctly with 'after' verts")
		}
		if faceIdx == 0 {
			continue
		}

		g.findClosestVerts(face)

		logf("c=%q: face[%v] cut0Idx=%v %v", g.unicode, faceIdx, face.cut0Idx, g.faces[0].after[face.cut0Idx])
		logf("c=%q: face[%v] cutIdx=%v %v", g.unicode, faceIdx, face.cutIdx, face.after[face.cutIdx])
	}
}

func (g *glyphT) regenerateFace() {
	g.findCutPoints()

	for i, face := range g.faces {
		logf("face[%v] new absCmd:\n%v", i, strings.Join(face.absCmds, ""))
	}

	var d strings.Builder
	face0 := g.faces[0]
	for f0vertIdx := range face0.after {
		d.WriteString(face0.absCmds[f0vertIdx])
		for _, face := range g.faces[1:] {
			if face.cut0Idx != f0vertIdx {
				continue
			}
			v := face.after[face.cutIdx]
			d.WriteString(fmt.Sprintf("L%v %v", v[0], v[1])) // jump over to face
			for i := range face.after {
				newIdx := (i + face.cutIdx + 1) % len(face.after)
				if newIdx == 0 {
					continue // skip after[0] - initial "M"
				}
				d.WriteString(face.absCmds[newIdx])
			}
			v = face0.after[face.cut0Idx]
			d.WriteString(fmt.Sprintf("L%v %v", v[0], v[1])) // jump back to face0
		}
	}
	g.d = d.String()
}

func writeFont(fontData *webfont.FontData) {
	p := &processor{glyphs: map[string]*glyphT{}}
	if err := webfont.ParseNeededGlyphs(fontData, "a", p); err != nil { // DEBUGGING ONLY
		log.Fatalf("webfont: %v", err)
	}

	var lines []string
	for unicode, glyph := range p.glyphs {
		lines = append(lines, fmt.Sprintf("%v={", unicode))
		lines = append(lines, fmt.Sprintf("    horiz_adv_x=%v,", glyph.horizAdvX))
		lines = append(lines, fmt.Sprintf(`    gerber_lp="%v",`, glyph.gerberLP))
		lines = append(lines, fmt.Sprintf(`    d="%v",`, glyph.d))
		// lines = append(lines, fmt.Sprintf(`    d_orig="%v",`, glyph.dOrig))
		lines = append(lines, fmt.Sprintf("    xmin=%v,", glyph.xmin))
		lines = append(lines, fmt.Sprintf("    ymin=%v,", glyph.ymin))
		lines = append(lines, fmt.Sprintf("    xmax=%v,", glyph.xmax))
		lines = append(lines, fmt.Sprintf("    ymax=%v,", glyph.ymax))
		lines = append(lines, "},")
	}

	fontData.Font.Data = strings.Join(lines, "\n    ")

	var buf bytes.Buffer
	if err := outTemp.Execute(&buf, fontData.Font); err != nil {
		log.Fatal(err)
	}

	filename := fmt.Sprintf("font_%v.lua", fontData.Font.ID)
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
}

func logf(fmt string, args ...any) {
	if *debug {
		log.Printf(fmt, args...)
	}
}

func viewFilter(s *string) string {
	if s == nil || !utf8.ValidString(*s) {
		return ""
	}

	r := webfont.UTF8toRune(s)
	if r == 0xfeff {
		return "" // BOM disallowed in Go source.
	}

	switch *s {
	case "\n", "\r", "\t":
		return ""
	default:
		return *s
	}
}

func orEmpty(s *string) string {
	if s == nil || *s == "" {
		return `""`
	}
	return fmt.Sprintf("%q", *s)
}

func floats(f []float64) string {
	return fmt.Sprintf("%#v", f)
}

var readmeTemplate = `# {{ .ID }}

![{{ .ID }}]({{ .ID }}.png)

To use this font in your code, simply import it:

` + "```" + `go
import (
  . "github.com/gmlewis/go-fonts/fonts"
  _ "github.com/gmlewis/go-fonts/fonts/{{ .ID }}"
)

func main() {
  // ...
  render, err := fonts.Text(xPos, yPos, xScale, yScale, message, "{{ .ID }}", Center)
  if err != nil {
    return err
  }
  log.Printf("MBB: %v", render.MBB)
  for _, poly := range render.Polygons {
    // ...
  }
  // ...
}
` + "```" + `
`

var luaTemplate = `-- Auto-generated by font2lua - DO NOT EDIT!
local F = require("font_library")
local P = require("params")
local V = require("vector_math")

local glyphs = {
    {{ .Data }}
}

FontLibrary:addFonts(
    {
        {{ .ID }} = {
            id = "{{ .ID }}",
            all_glyphs = [[{{ range .Glyphs }}{{ .Unicode | viewFilter }}{{ end }}]],
            horiz_adv_x = {{ .HorizAdvX }},
            units_per_em = {{ .FontFace.UnitsPerEm }},
            ascent = {{ .FontFace.Ascent }},
            descent = {{ .FontFace.Descent }},
            missing_horiz_adv_x = {{ .MissingGlyph.HorizAdvX }},
            glyphs = glyphs,
        }
    }
)
`
