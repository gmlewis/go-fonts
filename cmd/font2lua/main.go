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
	"os"
	"path/filepath"
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
	outTemp = template.Must(template.New("out").Funcs(funcMap).Parse(luaTemplate))
	funcMap = template.FuncMap{
		"floats":     floats,
		"orEmpty":    orEmpty,
		"viewFilter": viewFilter,
	}
	readmeTemp = template.Must(template.New("readme").Parse(readmeTemplate))

	digitRE = regexp.MustCompile(`^\d`)
)

func main() {
	flag.Parse()

	// for _, arg := range flag.Args() {
	arg := "../../fonts/freesans/FreeSans.svg"
	log.Printf("Processing file %q ...", arg)

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

	readme := genReadme(fontData)
	license := genLicense(filepath.Dir(arg))

	writeFont(fontData, readme, license)
	// }

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
	dParts  []string
	verts   []vec2
	center  vec2
	cut0Idx int // face[0] vertIdx to cut to connect to this face
	cutIdx  int // vertIdx to cut to connect to face[0]
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
		glyph.regenerateFace()
	}
}

func (p *processor) findDParts(cmd string) string {
	glyph := p.current
	dParts := glyph.d
	if len(glyph.faces) > 0 {
		face := glyph.faces[len(glyph.faces)-1]
		dParts = face.dParts[len(face.dParts)-1]
	}
	var strIdx int
	if len(glyph.faces) > 0 {
		strIdx = 1 + strings.Index(dParts[1:], cmd)
		face := glyph.faces[len(glyph.faces)-1]
		dIdx := len(face.dParts) - 1
		face.dParts[dIdx] = face.dParts[dIdx][:strIdx]
		if face.dParts[dIdx] == "" {
			log.Fatalf("programming error: findDParts(%q): face.dParts[%v]='' glyph=%#v", cmd, dIdx, *glyph)
		}
	}
	result := dParts[strIdx:]
	if result == "" {
		log.Fatalf("programming error: findDParts(%q): glyph=%#v", cmd, *glyph)
	}
	return result
}

func (p *processor) addCmd(glyph *glyphT, cmd string, x, y float64) {
	dParts := p.findDParts(cmd)
	if cmd == "M" { // start a new face
		glyph.faces = append(glyph.faces, &faceT{
			dParts: []string{dParts},
			verts:  []vec2{{x + glyph.xmin, y + glyph.ymin}},
		})
		return
	}

	face := len(glyph.faces) - 1
	glyph.faces[face].dParts = append(glyph.faces[face].dParts, dParts)
	glyph.faces[face].verts = append(glyph.faces[face].verts,
		vec2{x + glyph.xmin, y + glyph.ymin})
}

func (p *processor) MoveTo(g *webfont.Glyph, cmd string, x, y float64) {
	glyph := p.current
	log.Printf("p.MoveTo(g,%q, %v,%v)", cmd, x+glyph.xmin, y+glyph.ymin)
	p.addCmd(glyph, cmd, x, y)
}

func (p *processor) LineTo(g *webfont.Glyph, cmd string, x, y float64) {
	glyph := p.current
	log.Printf("p.LineTo(g,%q, %v,%v)", cmd, x+glyph.xmin, y+glyph.ymin)
	p.addCmd(glyph, cmd, x, y)
}

func (p *processor) CubicTo(g *webfont.Glyph, cmd string, x1, y1, x2, y2, ex, ey float64) {
	glyph := p.current
	log.Printf("p.CubicTo(g,%q,%v,%v,%v,%v,%v,%v)", cmd, x1+glyph.xmin, y1+glyph.ymin, x2+glyph.xmin, y2+glyph.ymin, ex+glyph.xmin, ey+glyph.ymin)
	p.addCmd(glyph, cmd, ex, ey)
}

func (p *processor) QuadraticTo(g *webfont.Glyph, cmd string, x1, y1, x2, y2 float64) {
	glyph := p.current
	log.Printf("p.QuadraticTo(g,%q,%v,%v,%v,%v)", cmd, x1+glyph.xmin, y1+glyph.ymin, x2+glyph.xmin, y2+glyph.ymin)
	p.addCmd(glyph, cmd, x2, y2)
}

func (g *glyphT) regenerateFace() {
	log.Printf("c=%q: d=%v", g.unicode, g.d)
	log.Printf("c=%q: got %v faces", g.unicode, len(g.faces))

	for faceIdx, face := range g.faces {
		log.Printf("c=%q: face[%v]: verts(%v): %#v", g.unicode, faceIdx, len(face.verts), face.verts)
		log.Printf("c=%q: face[%v]: dParts(%v): %#v", g.unicode, faceIdx, len(face.dParts), face.dParts)
		if len(face.verts) != len(face.dParts) {
			log.Fatalf("programming error - dParts not split correctly")
		}
		if faceIdx == 0 {
			continue
		}
		for _, vert := range face.verts {
			face.center[0] += vert[0]
			face.center[1] += vert[1]
		}
		numVerts := float64(len(face.verts))
		face.center[0] /= numVerts
		face.center[1] /= numVerts
		log.Printf("c=%q: face[%v] center=%v", g.unicode, faceIdx, face.center)
		var bestDistSq float64
		for vertIdx, vert := range g.faces[0].verts {
			dx := face.center[0] - vert[0]
			dy := face.center[1] - vert[1]
			distSq := dx*dx + dy*dy
			if vertIdx == 0 || distSq < bestDistSq {
				face.cut0Idx = vertIdx
				bestDistSq = distSq
				// log.Printf("c=%q: vert[%v]=%v - bestDistSq=%v", g.unicode, vertIdx, vert, distSq)
			}
		}
		f0x := g.faces[0].verts[face.cut0Idx][0]
		f0y := g.faces[0].verts[face.cut0Idx][1]
		log.Printf("c=%q: face[%v] cut0Idx=%v [%v %v]", g.unicode, faceIdx, face.cut0Idx, f0x, f0y)
		for vertIdx, vert := range face.verts {
			dx := f0x - vert[0]
			dy := f0y - vert[1]
			distSq := dx*dx + dy*dy
			if vertIdx == 0 || distSq < bestDistSq {
				face.cutIdx = vertIdx
				bestDistSq = distSq
				// log.Printf("c=%q: vert[%v]=%v - bestDistSq=%v", g.unicode, vertIdx, vert, distSq)
			}
		}
		log.Printf("c=%q: face[%v] cutIdx=%v %v", g.unicode, faceIdx, face.cutIdx, face.verts[face.cutIdx])
	}

	// Now re-generate the 'd' path so that it is one closed path with no holes.
	var d strings.Builder
	face0 := g.faces[0]
	for f0vertIdx, f0vert := range face0.verts {
		log.Printf("face0[%v]=%v", f0vertIdx, f0vert)
		d.WriteString(face0.dParts[f0vertIdx])
		for _, face := range g.faces[1:] {
			if face.cut0Idx != f0vertIdx {
				continue
			}
			v := face.verts[face.cutIdx]
			d.WriteString(fmt.Sprintf("L%v %v", v[0], v[1]))
			for i, vert := range face.verts {
				log.Printf("facei[%v]=%v", i, vert)
				newIdx := (i + face.cutIdx) % len(face.verts)
				d.WriteString(face.dParts[newIdx])
			}
		}
	}
}

func writeFont(fontData *webfont.FontData, readme, license string) {
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

func genReadme(fontData *webfont.FontData) string {
	var buf bytes.Buffer
	if err := readmeTemp.Execute(&buf, fontData.Font); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

func genLicense(srcDir string) string {
	// Copy any license along with the font.
	txtFiles, err := filepath.Glob(filepath.Join(srcDir, "*.txt"))
	if err != nil || len(txtFiles) == 0 {
		log.Printf("WARNING: unable to find license file in %v : %v", srcDir, err)
		return ""
	}
	var lines []string
	for _, txtFile := range txtFiles {
		buf, err := os.ReadFile(txtFile)
		if err != nil {
			log.Printf("WARNING: unable to read text file %v : %v", txtFile, err)
			continue
		}
		lines = append(lines, string(buf))
	}

	return strings.Join(lines, "\n")
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
