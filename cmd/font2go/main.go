// font2go reads one or more standard SVG webfont file(s) and writes Go file(s)
// used to render them to polygons.
package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

const (
	prefix = "fonts"
)

var (
	outTemp = template.Must(template.New("out").Funcs(funcMap).Parse(goTemplate))
	funcMap = template.FuncMap{
		"floats":  floats,
		"orEmpty": orEmpty,
		"utf8":    utf8Escape,
	}
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

		sort.Slice(fontData.Font.Glyphs, func(a, b int) bool {
			return strings.Compare(*fontData.Font.Glyphs[a].Unicode, *fontData.Font.Glyphs[b].Unicode) < 0
		})

		for _, g := range fontData.Font.Glyphs {
			g.ParsePath()
			g.GenGerberLP(fontData.Font.FontFace)
		}

		var buf bytes.Buffer
		if err := outTemp.Execute(&buf, fontData.Font); err != nil {
			log.Fatal(err)
		}

		fontDir := filepath.Join(prefix, fontData.Font.ID)
		if err := os.MkdirAll(fontDir, 0755); err != nil {
			log.Fatal(err)
		}
		filename := filepath.Join(fontDir, "font.go")
		fmtBuf, err := format.Source(buf.Bytes())
		if err != nil {
			ioutil.WriteFile(filename, buf.Bytes(), 0644) // Dump the unformatted output.
			log.Fatal(err)
		}

		if err := ioutil.WriteFile(filename, fmtBuf, 0644); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Done.")
}

func utf8Escape(s *string) string {
	if s == nil || *s == "" {
		return `""`
	}
	return fmt.Sprintf("%+q", *s)
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

var goTemplate = `// Auto-generated - DO NOT EDIT!

package {{ .ID }}

import (
	"github.com/gmlewis/go-fonts/fonts"
)

func init() {
  fonts.Fonts["{{ .ID }}"] = {{ .ID }}Font
}

var {{ .ID }}Font = &fonts.Font{
	// ID: "{{ .ID }}",
	HorizAdvX:  {{ .HorizAdvX }},
	UnitsPerEm: {{ .FontFace.UnitsPerEm }},
	Ascent:     {{ .FontFace.Ascent }},
	Descent:    {{ .FontFace.Descent }},
	MissingHorizAdvX: {{ .MissingGlyph.HorizAdvX }},
	Glyphs: map[string]*fonts.Glyph{ {{ range .Glyphs }}{{ if .Unicode }}
		{{ .Unicode | utf8 }}: {
			HorizAdvX: {{ .HorizAdvX }},
			Unicode: {{ .Unicode | utf8 }},
			GerberLP: {{ .GerberLP | orEmpty }},
			PathSteps: []*fonts.PathStep{ {{ range .PathSteps }}
				{ C: '{{ .C }}'{{ if .P }}, P: {{ .P | floats }}{{ end }} },{{ end }}
			},
		},{{ end }}{{ end }}
	},
}
`
