// -*- compile-command: "go run main.go ../../fonts/freesans/FreeSans.svg && cp font_*.lua /Users/glenn/src/github.com/gmlewis/blackjack/blackjack_lua/run"; -*-

// font2lua reads one or more standard SVG webfont file(s) and writes Lua file(s)
// used to render them to polygons in Blackjack and go-bjk.
// See: https://github.com/setzer22/blackjack
// and: https://github.com/gmlewis/go-bjk
//
// Usage:
//
//	font2lua fonts/*/*.svg
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
	"sort"
	"strconv"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/gmlewis/go-fonts/webfont"
)

var (
	debug = flag.String("debug", "", "Turn on debugging info for specific glyph ('all' for all glyphs)")

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
		logf("all", "Processing file %q ...", arg)

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

	// if glyph.gerberLP != "" {
	logf(glyph.unicode, "r=%q, unicode=%q, glyph.gerberLP=%q", r, glyph.unicode, glyph.gerberLP)
	glyph.regenerateFace()
	// }
}

func (p *processor) addCmd(glyph *glyphT, cmd string, x, y float64, absCmd string) {
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

func mySprintf(fmtStr string, coords ...float64) string {
	args := make([]any, 0, len(coords))
	for _, arg := range coords {
		v := fmt.Sprintf("%v", arg)
		if len(v) > 14 {
			v = fmt.Sprintf("%0.2f", arg)
			dotIdx := len(v) - 3
			if strings.HasSuffix(v, ".00") {
				v = v[:dotIdx]
			} else if v[dotIdx+2] == '0' {
				v = v[:dotIdx+2]
			}
		}
		args = append(args, v)
	}
	return fmt.Sprintf(fmtStr, args...)
}

func (p *processor) MoveTo(g *webfont.Glyph, cmd string, x, y float64) {
	glyph := p.current
	logf(glyph.unicode, "p.MoveTo(g,%q,%v,%v)", cmd, x+glyph.xmin, y+glyph.ymin)
	absCmd := mySprintf("M%v %v", x+glyph.xmin, y+glyph.ymin)
	p.addCmd(glyph, cmd, x, y, absCmd)
}

func (p *processor) LineTo(g *webfont.Glyph, cmd string, x, y float64) {
	glyph := p.current
	logf("p.LineTo(g,%q,%v,%v)", cmd, x+glyph.xmin, y+glyph.ymin)
	absCmd := mySprintf("L%v %v", x+glyph.xmin, y+glyph.ymin)
	p.addCmd(glyph, cmd, x, y, absCmd)
}

func (p *processor) CubicTo(g *webfont.Glyph, cmd string, x1, y1, x2, y2, ex, ey float64) {
	glyph := p.current
	logf("p.CubicTo(g,%q,%v,%v,%v,%v,%v,%v)", cmd, x1+glyph.xmin, y1+glyph.ymin, x2+glyph.xmin, y2+glyph.ymin, ex+glyph.xmin, ey+glyph.ymin)
	absCmd := mySprintf("C%v %v %v %v %v %v", x1+glyph.xmin, y1+glyph.ymin, x2+glyph.xmin, y2+glyph.ymin, ex+glyph.xmin, ey+glyph.ymin)
	p.addCmd(glyph, cmd, ex, ey, absCmd)
}

func (p *processor) QuadraticTo(g *webfont.Glyph, cmd string, x1, y1, x2, y2 float64) {
	glyph := p.current
	logf("p.QuadraticTo(g,%q,%v,%v,%v,%v)", cmd, x1+glyph.xmin, y1+glyph.ymin, x2+glyph.xmin, y2+glyph.ymin)
	absCmd := mySprintf("Q%v %v %v %v", x1+glyph.xmin, y1+glyph.ymin, x2+glyph.xmin, y2+glyph.ymin)
	p.addCmd(glyph, cmd, x2, y2, absCmd)
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
	// g.findCutPoints()

	for i, face := range g.faces {
		logf(g.unicode, "face[%v] new absCmd:\n%v", i, strings.Join(face.absCmds, ""))
	}

	var d strings.Builder
	/*
		face0 := g.faces[0]
		for f0vertIdx := range face0.after {
			d.WriteString(face0.absCmds[f0vertIdx])
			for _, face := range g.faces[1:] {
				if face.cut0Idx != f0vertIdx {
					continue
				}
				v := face.after[face.cutIdx]
				d.WriteString(mySprintf("L%v %v", v[0], v[1])) // jump over to face
				for i := range face.after {
					newIdx := (i + face.cutIdx + 1) % len(face.after)
					if newIdx == 0 {
						continue // skip after[0] - initial "M"
					}
					d.WriteString(face.absCmds[newIdx])
				}
				v = face0.after[face.cut0Idx]
				d.WriteString(mySprintf("L%v %v", v[0], v[1])) // jump back to face0
			}
		}
		d.WriteString("Z") // terminate the face.
	*/

	for _, face := range g.faces {
		d.WriteString(strings.Join(face.absCmds, ""))
		d.WriteString("Z") // terminate each face.
	}

	g.d = d.String()
}

func writeFont(fontData *webfont.FontData) {
	p := &processor{glyphs: map[string]*glyphT{}}
	if err := webfont.ParseNeededGlyphs(fontData, "", p); err != nil { // DEBUGGING ONLY
		log.Fatalf("webfont: %v", err)
	}

	keys := make([]string, 0, len(p.glyphs))
	for k := range p.glyphs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var lines []string
	for _, unicode := range keys {
		glyph := p.glyphs[unicode]

		luaChar := unicode
		if v, ok := luaMap[luaChar]; ok {
			luaChar = v
		} else {
			v := fmt.Sprintf("%+q", unicode)
			if strings.Contains(v, `\u`) {
				luaChar = strings.Replace(strings.Trim(v, `"`), `\u`, "u", -1)
			} else if strings.Contains(v, `\U`) {
				luaChar = strings.Replace(strings.Trim(v, `"`), `\U`, "U", -1)
			}
		}

		charOrHex := fmt.Sprintf("%q", unicode)
		if strings.HasPrefix(charOrHex, `"\u`) {
			v, err := strconv.ParseInt(charOrHex[3:len(charOrHex)-1], 16, 64)
			if err != nil {
				log.Fatalf("unable to parse hex '%v': %v", charOrHex, err)
			}
			charOrHex = hex2luaUnicode(v) // fmt.Sprintf(`"\%v"`, v)
		}

		lines = append(lines, fmt.Sprintf("%v={", luaChar))
		lines = append(lines, fmt.Sprintf(`    char=%v,`, charOrHex))
		lines = append(lines, fmt.Sprintf("    horiz_adv_x=%v,", glyph.horizAdvX))
		// lines = append(lines, fmt.Sprintf(`    gerber_lp="%v",`, glyph.gerberLP))
		lines = append(lines, fmt.Sprintf(`    d="%v",`, glyph.d))
		lines = append(lines, fmt.Sprintf("    xmin=%v,", glyph.xmin))
		lines = append(lines, fmt.Sprintf("    ymin=%v,", glyph.ymin))
		lines = append(lines, fmt.Sprintf("    xmax=%v,", glyph.xmax))
		lines = append(lines, fmt.Sprintf("    ymax=%v,", glyph.ymax))
		lines = append(lines, "},")
	}

	fontData.Font.Data = strings.Join(lines, "\n    ")
	if fontData.Font.HorizAdvX == 0 {
		fontData.Font.HorizAdvX = fontData.Font.MissingGlyph.HorizAdvX
	}

	var buf bytes.Buffer
	if err := outTemp.Execute(&buf, fontData.Font); err != nil {
		log.Fatal(err)
	}

	filename := fmt.Sprintf("font_%v.lua", fontData.Font.ID)
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
}

func hex2luaUnicode(v int64) string {
	if v == 0 {
		return `"\000"`
	}
	var out strings.Builder
	out.WriteString(`"`)
	for v > 0 {
		rem := v % 256
		out.WriteString(fmt.Sprintf(`\%03d`, rem))
		v /= 256
	}
	out.WriteString(`"`)
	return out.String()
}

func logf(unicode, fmt string, args ...any) {
	if *debug == "all" || *debug == unicode {
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

var luaTemplate = `-- Auto-generated by font2lua - DO NOT EDIT!
-- For font license and info, see: https://github.com/gmlewis/go-fonts/
local F = require("font_library")

local glyphs = {
    {{ .Data }}
}

F:addFonts(
    {
        {{ .ID }} = {
            id = "{{ .ID }}",
            horiz_adv_x = {{ .HorizAdvX }},
            units_per_em = {{ .FontFace.UnitsPerEm }},
            ascent = {{ .FontFace.Ascent }},
            descent = {{ .FontFace.Descent }},
            glyphs = glyphs,
        }
    }
)
`

//             all_glyphs = [[{{ range .Glyphs }}{{ .Unicode | viewFilter }}{{ end }}]],

var luaMap = map[string]string{
	"\t": "tab",
	"\n": "newline",
	"\r": "ctrlm",
	" ":  "space",
	"!":  "bang",
	`"`:  "dblquote",
	"#":  "hash",
	"$":  "dollar",
	"%":  "percent",
	"&":  "ampersand",
	"'":  "quote",
	"(":  "openparen",
	")":  "closeparen",
	"{":  "openbrace",
	"}":  "closebrace",
	"[":  "openbracket",
	"]":  "closebracket",
	"*":  "asterisk",
	"+":  "plus",
	",":  "comma",
	"-":  "minus",
	".":  "period",
	"/":  "slash",
	"\\": "backslash",
	"0":  "d0",
	"1":  "d1",
	"2":  "d2",
	"3":  "d3",
	"4":  "d4",
	"5":  "d5",
	"6":  "d6",
	"7":  "d7",
	"8":  "d8",
	"9":  "d9",
	":":  "colon",
	";":  "semicolon",
	"<":  "lt",
	">":  "gt",
	"=":  "eq",
	"?":  "questionmark",
	"@":  "at",
	"^":  "caret",
	"_":  "underscore",
	"`":  "backtick",
	"|":  "pipe",
	"~":  "tilde",
}
