// font2irmf reads one or more standard SVG webfont file(s) and writes IRMF file(s)
// used to render them in IRMF model shaders.
//
// For more information about IRMF, please see:
//   - https://github.com/gmlewis/irmf
//   - https://github.com/gmlewis/irmf-editor
//   - https://github.com/gmlewis/irmf-slicer
package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/gmlewis/go-fonts/webfont"
	"github.com/gmlewis/go3d/float64/vec2"
)

var (
	message = flag.String("msg", "IRMF fonts", "Message to spell. If empty, whole font is output.")
	verbose = flag.Bool("v", false, "Verbose debugging output")

	digitRE = regexp.MustCompile(`^\d`)
)

const (
	mmPerEm = 10
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		log.Printf("\n\nProcessing file %q ...", arg)

		fontData := &webfont.FontData{}
		if buf, err := os.ReadFile(arg); err != nil {
			log.Fatal(err)
		} else {
			if err := xml.Unmarshal(buf, fontData); err != nil {
				log.Fatal(err)
			}
		}

		sanitizeFontName(fontData)

		outFilename := fmt.Sprintf("%v.irmf", fontData.Font.ID)
		w, err := os.Create(outFilename)
		if err != nil {
			log.Fatalf("Create: %v", err)
		}

		writeFont(w, fontData, *message)

		if err := w.Close(); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Done.")
}

func sanitizeFontName(fontData *webfont.FontData) {
	fontData.Font.ID = strings.ToLower(fontData.Font.ID)
	fontData.Font.ID = strings.Replace(fontData.Font.ID, "'", "", -1)
	fontData.Font.ID = strings.Replace(fontData.Font.ID, "-", "_", -1)
	fontData.Font.ID = strings.TrimSuffix(fontData.Font.ID, "_")
	if digitRE.MatchString(fontData.Font.ID) {
		fontData.Font.ID = "f" + fontData.Font.ID
	}
}

func writeFont(w io.Writer, fontData *webfont.FontData, msg string) {
	buf := &bytes.Buffer{}
	rec := &recorder{f: buf, dedup: map[rune]*webfont.Glyph{}}
	if err := webfont.ParseNeededGlyphs(fontData, msg, rec); err != nil {
		log.Fatalf("webfont: %v", err)
	}

	emSize := fontData.Font.FontFace.Ascent

	var mbb *webfont.MBB
	if msg != "" {
		var lines []string
		var offset float64
		for _, r := range msg {
			g := rec.dedup[r]
			if g == nil {
				offset += emSize
				continue
			}

			glyphName := *g.Unicode
			if gn, ok := safeGlyphName[glyphName]; ok {
				glyphName = gn
			}
			logf("glyph %q: mbb=%v", glyphName, g.MBB)
			if g.MBB.Min[0] < g.MBB.Max[0] {
				lines = append(lines, fmt.Sprintf("  result += glyph_%v(xyz.xy-vec2(%v,0));", glyphName, offset))
			}

			if mbb == nil {
				mbb = &webfont.MBB{Min: g.MBB.Min, Max: g.MBB.Max}
				logf("Initial mbb=%v", mbb)
			} else {
				shiftedMBB := &webfont.MBB{
					Min: vec2.T{g.MBB.Min[0] + offset, g.MBB.Min[1]},
					Max: vec2.T{g.MBB.Max[0] + offset, g.MBB.Max[1]},
				}
				logf("shiftedMBB=%v", shiftedMBB)
				mbb.Join(shiftedMBB)
				logf("Updated mbb=%v", mbb)
			}

			offset += g.HorizAdvX
		}

		fmt.Fprintf(buf, `
float textMessage(in float mmPerEm, in float height, in vec3 xyz) {
  xyz *= vec3(%v,%v,1) / vec3(mmPerEm,mmPerEm,height);
  xyz += vec3(%v,%v,0);
  if (abs(xyz.z) > 0.5) { return 0.0; }
  float result = 0.0;
%v
  return result;
}

void mainModel4(out vec4 materials, in vec3 xyz) {
  materials[0] = textMessage(float(%v),0.1,xyz);
}
`, emSize, emSize,
			0.5*(mbb.Min[0]+mbb.Max[0]), 0.5*(mbb.Min[1]+mbb.Max[1]),
			strings.Join(lines, "\n"),
			mmPerEm)
	}

	logf("\n\nFinal mbb=%v", mbb)

	// Write header with helper functions.
	width := (mbb.Max[0] - mbb.Min[0]) * mmPerEm / emSize
	height := (mbb.Max[1] - mbb.Min[1]) * mmPerEm / emSize
	fmt.Fprintf(w, header, 0.5*width, 0.5*height, -0.5*width, -0.5*height)

	// Write methods.
	fmt.Fprintf(w, "%s", buf.Bytes())
}

func logf(fmtStr string, args ...interface{}) {
	if *verbose {
		log.Printf(fmtStr, args...)
	}
}

var header = `/*{
  "author": "",
  "copyright": "",
  "date": "",
  "irmf": "1.0",
  "materials": ["PLA"],
  "max": [%v,%v,0.5],
  "min": [%v,%v,-0.5],
  "notes": "",
  "options": {},
  "title": "",
  "units": "mm",
  "version": ""
}*/

float blinnLoop(in vec2 A, in vec2 B, in vec2 C) {
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

float interpLine(in vec2 A, in vec2 B, in float y) {
  float p = (y - A.y) / (B.y - A.y);
  return p*(B.x-A.x) + A.x;
}

float interpQuadratic(in vec2 p0, in vec2 p1, in vec2 p2, in float y) {
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
`
