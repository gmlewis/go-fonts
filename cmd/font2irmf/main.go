// font2irmf reads one or more standard SVG webfont file(s) and writes IRMF file(s)
// used to render them in IRMF model shaders.
//
// For more information about IRMF, please see:
//   - https://github.com/gmlewis/irmf
//   - https://github.com/gmlewis/irmf-editor
//   - https://github.com/gmlewis/irmf-slicer
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/davecgh/go-spew/spew"
)

const (
	prefix = "fonts"
)

var (
	message = flag.String("msg", "$", "Message to spell. If empty, whole font is output.")
	digitRE = regexp.MustCompile(`^\d`)
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		log.Printf("Processing file %q ...", arg)

		fontData := &FontData{}
		if buf, err := ioutil.ReadFile(arg); err != nil {
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

		fontDir := filepath.Join(prefix, fontData.Font.ID)
		if err := os.MkdirAll(fontDir, 0755); err != nil {
			log.Fatal(err)
		}

		writeFont(fontData, fontDir, *message)
	}

	fmt.Println("Done.")
}

func writeFont(fontData *FontData, fontDir, msg string) {
	glyphLess := func(a, b int) bool {
		sa, sb := "", ""
		if fontData.Font.Glyphs[a].Unicode != nil {
			sa = *fontData.Font.Glyphs[a].Unicode
		}
		if fontData.Font.Glyphs[b].Unicode != nil {
			sb = *fontData.Font.Glyphs[b].Unicode
		}
		return strings.Compare(sa, sb) < 0
	}

	sort.Slice(fontData.Font.Glyphs, glyphLess)

	// Fix UTF8 rune errors and de-duplicate identical code points.
	dedup := map[rune]*Glyph{}
	var dst rune = 0xfbf0
	for _, g := range fontData.Font.Glyphs {
		if g.Unicode == nil {
			continue
		}
		r := utf8toRune(g.Unicode)
		if r == 0 {
			log.Fatalf("Unicode %+q is mapping to r=0 !!!", *g.Unicode)
			continue
		}
		if msg != "" && !strings.ContainsRune(msg, r) {
			continue
		}
		if _, ok := dedup[r]; ok {
			if dst == 0xfeff { // BOM - disallowed in Go source.
				dst++
			}
			for {
				if _, ok := dedup[dst]; !ok {
					break
				}
				dst++
			}
			log.Printf("WARNING: unicode %+q found multiple times in font. Moving code point to %+q", r, dst)
			rs := fmt.Sprintf("%c", dst)
			g.Unicode = &rs
			dedup[dst] = g
			dst++
			continue
		}
		rs := fmt.Sprintf("%c", r)
		g.Unicode = &rs
		dedup[r] = g
	}

	// re-sort with deduped glyph code points.
	sort.Slice(fontData.Font.Glyphs, glyphLess)

	filename := filepath.Join(fontDir, "font.irmf")
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Create: %v", err)
	}

	// Write helper functions.
	fmt.Fprintf(f, `/*{
  "author": "",
  "copyright": "",
  "date": "",
  "irmf": "1.0",
  "materials": ["PLA"],
  "max": [5,5,5],
  "min": [-5,-5,-5],
  "notes": "",
  "options": {},
  "title": "",
  "units": "mm",
  "version": ""
}*/

float blinnLoop(vec2 A, vec2 B, vec2 C) {
  vec2 v0 = C - A;
  vec2 v1 = B - A;
  vec2 v2 = vec2(0.75,0.5) - A;
  // Compute dot products
  float dot00 = dot(v0, v0);
  float dot01 = dot(v0, v1);
  float dot02 = dot(v0, v2);
  float dot11 = dot(v1, v1);
  float dot12 = dot(v1, v2);
  // Compute barycentric coordinates
  float invDenom = 1.0 / (dot00 * dot11 - dot01 * dot01);
  float u = (dot11 * dot02 - dot01 * dot12) * invDenom;
  float v = (dot00 * dot12 - dot01 * dot02) * invDenom;
  // use the blinn and loop method
  float w = (1.0 - u - v);
  return w;
}

float interpLine(vec2 A, vec2 B, float y) {
  float p = (y - A.y) / (B.y - A.y);
  return p*(B.x-A.x) + A.x;
}

float interpQuadratic(vec2 p0, vec2 p1, vec2 p2, float y) {
  float a = p2.y + p0.y - 2.0*p1.y;
  float b = 2.0 * (p1.y - p0.y);
  float c = p0.y - y;
  if (b*b < 4.0*a*c) { return 0.0; } // bad (imaginary) quadratic
  float det = sqrt(b*b - 4.0*a*c);
  float t = (-b + det) / (2.0 * a);
  float t2 = (-b - det) / (2.0 * a);
  if (t2 >= 0.0 && t2 <= 1.0) {
    t = t2;
  }
  float x = (1.0-t)*(1.0-t)*p0.x + 2.0*(1.0-t)*t*p1.x + t*t*p2.x;
  return x;
}
`)

	for _, g := range fontData.Font.Glyphs {
		r := utf8toRune(g.Unicode)
		if g.Unicode == nil || (msg != "" && !strings.ContainsRune(msg, r)) {
			continue
		}
		g.ParsePath()
		g.GenGerberLP(fontData.Font.FontFace)
		if g.MBB.Area() == 0.0 {
			continue
		}
		log.Printf("Glyph %+q: mbb=%v", *g.Unicode, g.MBB)
		g.rec.process(f, g)
		spew.Dump(g)
	}

	if msg != "" {
		emSize := fontData.Font.FontFace.Ascent

		var lines []string
		var offset float64
		for _, r := range msg {
			g := dedup[r]
			glyphName := *g.Unicode
			if gn, ok := safeGlyphName[glyphName]; ok {
				glyphName = gn
			}
			lines = append(lines, fmt.Sprintf("  result += glyph_%v(height, xyz-vec3(%.2f,0,0));", glyphName, offset))
			offset += g.HorizAdvX
		}

		fmt.Fprintf(f, `
float textMessage(float emSize, float height, in vec3 xyz) {
  xyz *= vec3(%.2f/emSize,%.2f/emSize,1.0);
  float result = 0.0;
  %v
  return result;
}

void mainModel4(out vec4 materials, in vec3 xyz) {
  xyz += vec3(0.0, 5.0, 1.5);
  materials[0] = textMessage(10.0, 3.0, xyz);
}
`, emSize, emSize, strings.Join(lines, "\n"))
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func utf8toRune(s *string) rune {
	if s == nil || *s == "" {
		return 0
	}

	switch *s {
	case "\n":
		return '\n'
	case `\`:
		return '\\'
	case `'`:
		return '\''
	}

	if utf8.RuneCountInString(*s) == 1 {
		r, _ := utf8.DecodeRuneInString(*s)
		return r
	}
	if r, ok := specialCase[*s]; ok {
		return r
	}

	if len(*s) > 1 {
		log.Printf("WARNING: Unhandled unicode seqence: %+q", *s)
	}
	for _, r := range *s { // Return the first rune
		return r
	}
	return 0
}

// This specialCase map converts non-unicode strings (e.g. "ffi" - which
// is a 3-rune string - vs 'ﬃ' which is a 1-rune code point)
// to (basically-random) unicode runes so that they can still be
// rendered with some arbitrary rune code point. Note that the official code
// points were not used in some cases (e.g. "ffi" is not mapped to '\ufb03'
// because one or more of the open source fonts had both "ffi" and '\ufb03'
// code points already in the font!
//
// PRs that fix the code points will be rejected if they prevent one of
// the existing open source fonts from functioning properly.
var specialCase = map[string]rune{
	"1!":                       '\ufb10',
	"1#":                       '\ufb11',
	"1$":                       '\ufb12',
	"1%":                       '\ufb13',
	"1&":                       '\ufb14',
	"1(":                       '\ufb15',
	"1)":                       '\ufb16',
	"1*":                       '\ufb17',
	"1@":                       '\ufb18',
	"1^":                       '\ufb19',
	"Ex":                       '\ufb38',
	"Fi":                       '\ufb3b',
	"Gj":                       '\ufb1d',
	"IJ":                       '\ufb1e',
	"L\u00b7":                  '\ufbb8',
	"Qu":                       '\ufb8a',
	"Th":                       '\ufb58',
	"Ti":                       '\ufb59',
	"Tj":                       '\ufbc7',
	"Yj":                       '\ufb1f',
	"\u00e9x":                  '\ufb5c',
	"\u00edx":                  '\ufb5d',
	"\u00edz":                  '\ufb5e',
	"\u00f3s":                  '\ufb5f',
	"\u00f3sx":                 '\ufb60',
	"\u00f3sz":                 '\ufb61',
	"\u00f3x":                  '\ufb62',
	"\u0105,,":                 '\ufb7b',
	"\u0107,,":                 '\ufb7c',
	"\u0119,,":                 '\ufb7d',
	"\u0142,,":                 '\ufb7e',
	"\u0144,,":                 '\ufb7f',
	"\u017a,,":                 '\ufb80',
	"\u017c,,":                 '\ufb81',
	"\u0457\u0457":             '\ufbcc',
	"\u05f2\u05b7":             '\ufbcd',
	"\u064c\u0651":             '\ufbce',
	"\u064e\u0651":             '\ufbcf',
	"\u064f\u0651":             '\ufbd0',
	"\u0670\u0651":             '\ufbd1',
	"\u0e0d\u0e38":             '\ufbd2',
	"\u0e0d\u0e39":             '\ufbd3',
	"\u0e0d\u0e3a":             '\ufbd4',
	"\u0e10\u0e38":             '\ufbd5',
	"\u0e10\u0e39":             '\ufbd6',
	"\u0e10\u0e3a":             '\ufbd7',
	"\u0e24\u0e32":             '\ufbd8',
	"\u0e26\u0e32":             '\ufbd9',
	"\u1724\u1733":             '\ufbda',
	"\u1725\u1732":             '\ufbdb',
	"\u1725\u1733":             '\ufbdc',
	"\u1726\u1733":             '\ufbdd',
	"\u1727\u1733":             '\ufbde',
	"\u1729\u1732":             '\ufbdf',
	"\u1729\u1733":             '\ufbe0',
	"\u172b\u1733":             '\ufbe1',
	"\u172e\u1733":             '\ufbe2',
	"\u1731\u1732":             '\ufbe3',
	"\u1731\u1733":             '\ufbe4',
	"\u1a15\u1a17\u200d\u1a10": '\ufbe5',
	"\ue001\ue014":             '\u2469',
	"\ue001\ue015":             '\u246a',
	"\ue001\ue016":             '\u246b',
	"\ue001\ue017":             '\u246c',
	"\ue001\ue018":             '\u246d',
	"\ue001\ue019":             '\u246e',
	"\ue001\ue01a":             '\u246f',
	"\ue001\ue01b":             '\u2470',
	"\ue001\ue01c":             '\u2471',
	"\ue001\ue01d":             '\u2472',
	"\ue002\ue014":             '\u2473',
	"\ufedf\ufe82":             '\ufbe6',
	"\ufedf\ufe84":             '\ufbe7',
	"\ufedf\ufe88":             '\ufbe8',
	"\ufedf\ufe8e":             '\ufbe9',
	"\ufedf\ufee0\ufeea":       '\ufbea',
	"\ufee0\ufe84":             '\ufbeb',
	"\ufee0\ufe88":             '\ufbec',
	"\ufee0\ufe8e":             '\ufbed',
	"a,,":                      '\ufb20',
	"ar":                       '\ufba0',
	"as":                       '\ufba1',
	"ax":                       '\ufb30',
	"az":                       '\ufb31',
	"br":                       '\ufba2',
	"bs":                       '\ufb32',
	"bsx":                      '\ufb33',
	"bsz":                      '\ufb34',
	"c,,":                      '\ufb21',
	"c\u02c7\u0313":            '\ufba5',
	"c\u030c\u0313":            '\ufba6',
	"c\u030c\u0315":            '\ufba7',
	"cr":                       '\ufba3',
	"cs":                       '\ufba4',
	"ct":                       '\ufb87',
	"cx":                       '\ufb35',
	"cz":                       '\ufb36',
	"d,,":                      '\ufb22',
	"d\u0313":                  '\ufbaa',
	"d\u0315":                  '\ufbab',
	"dr":                       '\ufba8',
	"ds":                       '\ufba9',
	"e,,":                      '\ufb23',
	"er":                       '\ufbac',
	"es":                       '\ufbad',
	"ex":                       '\ufb37',
	"ez":                       '\ufb39',
	"f\u00ed":                  '\ufb88',
	"fb":                       '\ufb1c',
	"ff":                       '\ufb90', // Not \ufb00 - see above.
	"ffb":                      '\ufb07',
	"ffh":                      '\ufb08',
	"ffi":                      '\ufb93', // Not \ufb03 - see above.
	"ffitrk":                   '\ufbd3',
	"ffj":                      '\ufb09',
	"ffk":                      '\ufb0a',
	"ffl":                      '\ufb94', // Not \ufb04 - see above.
	"fft":                      '\ufb3a',
	"fh":                       '\ufb0b',
	"fi":                       '\ufb91', // Not \ufb01 - see above.
	"fitrk":                    '\ufbd1',
	"fix":                      '\ufb3c',
	"fiz":                      '\ufb3d',
	"fj":                       '\ufb0c',
	"fk":                       '\ufb0d',
	"fl":                       '\ufb92', // Not \ufb02 - see above.
	"flx":                      '\ufb3e',
	"flz":                      '\ufb3f',
	"ft":                       '\ufb05',
	"fu":                       '\ufb40',
	"fx":                       '\ufb41',
	"g,,":                      '\ufb68',
	"gi":                       '\ufbae',
	"gj":                       '\ufb69',
	"gp":                       '\ufb42',
	"gx":                       '\ufb43',
	"gz":                       '\ufb44',
	"h,,":                      '\ufb6a',
	"hr":                       '\ufbaf',
	"hs":                       '\ufbb0',
	"i,,":                      '\ufb6b',
	"ij":                       '\ufb45',
	"ir":                       '\ufbb1',
	"is":                       '\ufbb2',
	"ix":                       '\ufb46',
	"iz":                       '\ufb47',
	"j,,":                      '\ufb6c',
	"jj":                       '\ufb6d',
	"jx":                       '\ufb48',
	"jz":                       '\ufb49',
	"k,,":                      '\ufb6e',
	"kr":                       '\ufbb3',
	"ks":                       '\ufbb4',
	"l,,":                      '\ufb6f',
	"l\u00b7":                  '\ufbb7',
	"l\u0313":                  '\ufbb9',
	"l\u0315":                  '\ufbba',
	"lr":                       '\ufbb5',
	"ls":                       '\ufbb6',
	"lx":                       '\ufb4a',
	"lz":                       '\ufb4b',
	"m,,":                      '\ufb70',
	"mr":                       '\ufbbb',
	"ms":                       '\ufbbc',
	"mx":                       '\ufb4c',
	"mz":                       '\ufb4d',
	"n,,":                      '\ufb71',
	"nr":                       '\ufbbd',
	"ns":                       '\ufbbe',
	"nz":                       '\ufb4e',
	"o\u00e6":                  '\ufb52',
	"or":                       '\ufbbf',
	"os":                       '\ufb4f',
	"osx":                      '\ufb50',
	"osz":                      '\ufb51',
	"ox":                       '\ufb53',
	"pr":                       '\ufbc0',
	"ps":                       '\ufbc1',
	"qf":                       '\ufb1a',
	"qj":                       '\ufb1b',
	"r\u017c":                  '\ufb74',
	"rf":                       '\ufbc2',
	"rr":                       '\ufbc3',
	"rs":                       '\ufbc4',
	"rt":                       '\ufbc5',
	"ru":                       '\ufb73',
	"ru,,":                     '\ufb72',
	"rw":                       '\ufb75',
	"rx":                       '\ufb54',
	"ry":                       '\ufb77',
	"ry,,":                     '\ufb76',
	"rz":                       '\ufb55',
	"rz,,":                     '\ufb78',
	"ss":                       '\ufbc6',
	"st":                       '\ufb8b',
	"sx":                       '\ufb56',
	"sz":                       '\ufb57',
	"t,,":                      '\ufb79',
	"t\u0313":                  '\ufbca',
	"t\u0315":                  '\ufbcb',
	"ti":                       '\ufb8c',
	"tj":                       '\ufb8d',
	"tr":                       '\ufbc8',
	"ts":                       '\ufbc9',
	"tt":                       '\ufb0f',
	"tx":                       '\ufb5a',
	"tz":                       '\ufb5b',
	"u,,":                      '\ufb7a',
	"ur":                       '\ufbee',
	"us":                       '\ufbef',
	"uv":                       '\ufb63',
	"ux":                       '\ufb64',
	"vr":                       '\ufbf0',
	"vs":                       '\ufbf1',
	"wr":                       '\ufbf2',
	"ws":                       '\ufbf3',
	"www":                      '\ufbf4',
	"x,,":                      '\ufb82',
	"y,,":                      '\ufb83',
	"yf":                       '\ufb84',
	"yj":                       '\ufb85',
	"yp":                       '\ufb65',
	"yx":                       '\ufb66',
	"yz":                       '\ufb67',
	"z,,":                      '\ufb86',
}
