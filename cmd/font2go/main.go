// font2go reads one or more standard SVG webfont file(s) and writes Go file(s)
// used to render them to polygons.
package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/xml"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/gmlewis/go-fonts/pb/glyphs"
	"github.com/gmlewis/go-fonts/webfont"
	"google.golang.org/protobuf/proto"
)

const (
	prefix = "fonts"
)

var (
	readmeOnly = flag.Bool("readme_only", false, "Only write the README.md file")

	outTemp = template.Must(template.New("out").Funcs(funcMap).Parse(goTemplate))
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

	for _, arg := range flag.Args() {
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
		fontData.Font.ID = strings.TrimSuffix(fontData.Font.ID, "_")
		if digitRE.MatchString(fontData.Font.ID) {
			fontData.Font.ID = "f" + fontData.Font.ID
		}

		fontDir := filepath.Join(prefix, fontData.Font.ID)
		if err := os.MkdirAll(fontDir, 0755); err != nil {
			log.Fatal(err)
		}

		if !*readmeOnly {
			writeFont(fontData, fontDir)
		}

		writeReadme(fontData, fontDir)

		if !*readmeOnly {
			writeLicense(filepath.Dir(arg), fontDir)
		}
	}

	fmt.Println("Done.")
}

// processor implements the webfont.Processor interface.
type processor struct {
	gs *glyphs.Glyphs
}

var _ webfont.Processor = &processor{}

func (p *processor) ProcessGlyph(r rune, g *webfont.Glyph) {
	var pathSteps []*glyphs.PathStep
	for _, ps := range g.PathSteps {
		pathSteps = append(pathSteps, &glyphs.PathStep{C: uint32(ps.C[0]), P: ps.P})
	}
	gerberLP := ""
	if g.GerberLP != nil {
		gerberLP = *g.GerberLP
	}
	p.gs.Glyphs = append(p.gs.Glyphs, &glyphs.Glyph{
		HorizAdvX: g.HorizAdvX,
		Unicode:   *g.Unicode,
		GerberLP:  gerberLP,
		PathSteps: pathSteps,
		Mbb: &glyphs.MBB{
			Xmin: g.MBB.Min[0],
			Ymin: g.MBB.Min[1],
			Xmax: g.MBB.Max[0],
			Ymax: g.MBB.Max[1],
		},
	})
}

func (p *processor) NewGlyph(g *webfont.Glyph)                                            {}
func (p *processor) MoveTo(g *webfont.Glyph, cmd string, x, y float64)                    {}
func (p *processor) LineTo(g *webfont.Glyph, cmd string, x, y float64)                    {}
func (p *processor) CubicTo(g *webfont.Glyph, cmd string, x1, y1, x2, y2, ex, ey float64) {}
func (p *processor) QuadraticTo(g *webfont.Glyph, cmd string, x1, y1, x2, y2 float64)     {}

func writeFont(fontData *webfont.FontData, fontDir string) {
	p := &processor{gs: &glyphs.Glyphs{}}
	if err := webfont.ParseNeededGlyphs(fontData, "", p); err != nil {
		log.Fatalf("webfont: %v", err)
	}

	{
		data, err := proto.Marshal(p.gs)
		if err != nil {
			log.Fatal(err)
		}
		var b bytes.Buffer
		w, err := zlib.NewWriterLevel(&b, zlib.BestCompression)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(data)
		if err := w.Close(); err != nil {
			log.Fatalf("zlib.Close: %v", err)
		}
		fontData.Font.Data = base64.StdEncoding.EncodeToString(b.Bytes())
	}

	var buf bytes.Buffer
	if err := outTemp.Execute(&buf, fontData.Font); err != nil {
		log.Fatal(err)
	}

	filename := filepath.Join(fontDir, "font.go")
	fmtBuf, err := format.Source(buf.Bytes())
	if err != nil {
		os.WriteFile(filename, buf.Bytes(), 0644) // Dump the unformatted output.
		log.Fatalf("error formating generated Go code: %v : %v", filename, err)
	}

	if err := os.WriteFile(filename, fmtBuf, 0644); err != nil {
		log.Fatal(err)
	}
}

func writeReadme(fontData *webfont.FontData, fontDir string) {
	// Create README.md.
	var buf bytes.Buffer
	if err := readmeTemp.Execute(&buf, fontData.Font); err != nil {
		log.Fatal(err)
	}
	readmeName := filepath.Join(fontDir, "README.md")
	if err := os.WriteFile(readmeName, buf.Bytes(), 0644); err != nil {
		log.Printf("WARNING: unable to write %v : %v", readmeName, err)
	}
}

func copyFiles(filenames []string, fontDir string) {
	for _, filename := range filenames {
		buf, err := os.ReadFile(filename)
		if err != nil {
			log.Printf("WARNING: unable to read text file %v : %v", filename, err)
			continue
		}
		baseName := filepath.Base(filename)
		dstName := filepath.Join(fontDir, baseName)
		if err := os.WriteFile(dstName, buf, 0644); err != nil {
			log.Printf("WARNING: unable to write text file %v : %v", dstName, err)
			continue
		}
		log.Printf("Copied file to %v", dstName)
	}
}

func writeLicense(srcDir, fontDir string) {
	// Copy any license along with the font.
	txtFiles, err := filepath.Glob(filepath.Join(srcDir, "*.txt"))
	if err == nil && len(txtFiles) > 0 {
		copyFiles(txtFiles, fontDir)
	} else {
		log.Printf("WARNING: unable to find license file in %v : %v", srcDir, err)
	}
	// Also copy any README-orig.md files.
	txtFiles, err = filepath.Glob(filepath.Join(srcDir, "README-orig.md"))
	if err == nil && len(txtFiles) > 0 {
		copyFiles(txtFiles, fontDir)
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

var goTemplate = `// Auto-generated - DO NOT EDIT!

package {{ .ID }}

import (
	"github.com/gmlewis/go-fonts/fonts"
)

// Available glyphs:
// {{ range .Glyphs }}{{ .Unicode | viewFilter }}{{ end }}

var {{ .ID }}Font = &fonts.Font{
	ID:               "{{ .ID }}",
	HorizAdvX:  {{ .HorizAdvX }},
	UnitsPerEm: {{ .FontFace.UnitsPerEm }},
	Ascent:     {{ .FontFace.Ascent }},
	Descent:    {{ .FontFace.Descent }},
	MissingHorizAdvX: {{ .MissingGlyph.HorizAdvX }},
	Glyphs:           map[rune]*fonts.Glyph{},
}

func init() {
	fonts.Fonts["{{ .ID }}"] = {{ .ID }}Font
	fonts.InitFromFontData({{ .ID }}Font, fontData)
}

var fontData = ` + "`{{ .Data }}`" + ` 
`
