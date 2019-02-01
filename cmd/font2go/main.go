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
	"regexp"
	"sort"
	"strings"
	"text/template"
	"unicode/utf8"
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

		sort.Slice(fontData.Font.Glyphs, func(a, b int) bool {
			sa, sb := "", ""
			if fontData.Font.Glyphs[a].Unicode != nil {
				sa = *fontData.Font.Glyphs[a].Unicode
			}
			if fontData.Font.Glyphs[b].Unicode != nil {
				sb = *fontData.Font.Glyphs[b].Unicode
			}
			return strings.Compare(sa, sb) < 0
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
			log.Fatalf("error formating generated Go code: %v : %v", filename, err)
		}

		if err := ioutil.WriteFile(filename, fmtBuf, 0644); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Done.")
}

func utf8Escape(s *string) string {
	if s == nil || *s == "" {
		return `''`
	}
	switch *s {
	case `\`:
		return `'\\'`
	case `'`:
		return `'\''`
	}

	if utf8.RuneCountInString(*s) == 1 {
		return fmt.Sprintf("'%v'", *s)
	}

	if v, ok := specialCase[*s]; ok {
		return v
	}

	if len(*s) > 1 {
		log.Printf("WARNING: Unhandled unicode seqence: %+q", *s)
	}
	for _, r := range *s { // Return the first rune
		return fmt.Sprintf("'%c'", r)
	}
	return ""
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
	ID: "{{ .ID }}",
	HorizAdvX:  {{ .HorizAdvX }},
	UnitsPerEm: {{ .FontFace.UnitsPerEm }},
	Ascent:     {{ .FontFace.Ascent }},
	Descent:    {{ .FontFace.Descent }},
	MissingHorizAdvX: {{ .MissingGlyph.HorizAdvX }},
	Glyphs: map[rune]*fonts.Glyph{ {{ range .Glyphs }}{{ if .Unicode }}{{ if .PathSteps }}
		{{ .Unicode | utf8 }}: {
			HorizAdvX: {{ .HorizAdvX }},
			Unicode: {{ .Unicode | utf8 }},
			GerberLP: {{ .GerberLP | orEmpty }},
			PathSteps: []*fonts.PathStep{ {{ range .PathSteps }}
				{ C: '{{ .C }}'{{ if .P }}, P: {{ .P | floats }}{{ end }} },{{ end }}
			},
		},{{ end }}{{ end }}{{ end }}
	},
}
`

// This specialCase map converts non-unicode strings (e.g. "ffi" - which
// is a 3-rune string - vs 'ﬃ' which is a 1-rune code point)
// to (basically-random) unicode characters so that they can still be
// rendered with some arbitrary rune code point. Note that the official code
// points were not used in some cases (e.g. "ffi" is not mapped to '\ufb03'
// because one or more of the open source fonts had both "ffi" and '\ufb03'
// code points already in the font!
//
// PRs that fix the code points will be rejected if they prevent one of
// the existing open source fonts from functioning properly.
var specialCase = map[string]string{
	"1!":           `'\ufb10'`,
	"1#":           `'\ufb11'`,
	"1$":           `'\ufb12'`,
	"1%":           `'\ufb13'`,
	"1&":           `'\ufb14'`,
	"1(":           `'\ufb15'`,
	"1)":           `'\ufb16'`,
	"1*":           `'\ufb17'`,
	"1@":           `'\ufb18'`,
	"1^":           `'\ufb19'`,
	"a,,":          `'\ufb20'`,
	"ax":           `'\ufb30'`,
	"az":           `'\ufb31'`,
	"bs":           `'\ufb32'`,
	"bsx":          `'\ufb33'`,
	"bsz":          `'\ufb34'`,
	"ct":           `'\ufb87'`,
	"c,,":          `'\ufb21'`,
	"cx":           `'\ufb35'`,
	"cz":           `'\ufb36'`,
	"d,,":          `'\ufb22'`,
	"e,,":          `'\ufb23'`,
	"ex":           `'\ufb37'`,
	"Ex":           `'\ufb38'`,
	"ez":           `'\ufb39'`,
	"fb":           `'\ufb1c'`,
	"ffb":          `'\ufb07'`,
	"ffh":          `'\ufb08'`,
	"ffi":          `'\ufb93'`, // Not \ufb03 - see above.
	"ffj":          `'\ufb09'`,
	"ffk":          `'\ufb0a'`,
	"ffl":          `'\ufb94'`, // Not \ufb04 - see above.
	"fft":          `'\ufb3a'`,
	"ff":           `'\ufb90'`, // Not \ufb00 - see above.
	"fh":           `'\ufb0b'`,
	"fi":           `'\ufb91'`, // Not \ufb01 - see above.
	"Fi":           `'\ufb3b'`,
	"fix":          `'\ufb3c'`,
	"fiz":          `'\ufb3d'`,
	"fj":           `'\ufb0c'`,
	"fk":           `'\ufb0d'`,
	"fl":           `'\ufb92'`, // Not \ufb02 - see above.
	"flx":          `'\ufb3e'`,
	"flz":          `'\ufb3f'`,
	"ft":           `'\ufb05'`,
	"f\u00ed":      `'\ufb88'`,
	"fu":           `'\ufb40'`,
	"fx":           `'\ufb41'`,
	"Gj":           `'\ufb1d'`,
	"gj":           `'\ufb69'`,
	"gp":           `'\ufb42'`,
	"g,,":          `'\ufb68'`,
	"gx":           `'\ufb43'`,
	"gz":           `'\ufb44'`,
	"h,,":          `'\ufb6a'`,
	"IJ":           `'\ufb1e'`,
	"ij":           `'\ufb45'`,
	"i,,":          `'\ufb6b'`,
	"ix":           `'\ufb46'`,
	"iz":           `'\ufb47'`,
	"jj":           `'\ufb6d'`,
	"j,,":          `'\ufb6c'`,
	"jx":           `'\ufb48'`,
	"jz":           `'\ufb49'`,
	"k,,":          `'\ufb6e'`,
	"l,,":          `'\ufb6f'`,
	"lx":           `'\ufb4a'`,
	"lz":           `'\ufb4b'`,
	"m,,":          `'\ufb70'`,
	"mx":           `'\ufb4c'`,
	"mz":           `'\ufb4d'`,
	"n,,":          `'\ufb71'`,
	"nz":           `'\ufb4e'`,
	"os":           `'\ufb4f'`,
	"osx":          `'\ufb50'`,
	"osz":          `'\ufb51'`,
	"o\u00e6":      `'\ufb52'`,
	"ox":           `'\ufb53'`,
	"qf":           `'\ufb1a'`,
	"qj":           `'\ufb1b'`,
	"Qu":           `'\ufb8a'`,
	"r\u017c":      `'\ufb74'`,
	"ru,,":         `'\ufb72'`,
	"ru":           `'\ufb73'`,
	"rw":           `'\ufb75'`,
	"rx":           `'\ufb54'`,
	"ry,,":         `'\ufb76'`,
	"ry":           `'\ufb77'`,
	"rz":           `'\ufb55'`,
	"rz,,":         `'\ufb78'`,
	"st":           `'\ufb8b'`,
	"sx":           `'\ufb56'`,
	"sz":           `'\ufb57'`,
	"Th":           `'\ufb58'`,
	"ti":           `'\ufb8c'`,
	"Ti":           `'\ufb59'`,
	"tj":           `'\ufb8d'`,
	"tt":           `'\ufb0f'`,
	"t,,":          `'\ufb79'`,
	"tx":           `'\ufb5a'`,
	"tz":           `'\ufb5b'`,
	"\u00e9x":      `'\ufb5c'`,
	"\u00edx":      `'\ufb5d'`,
	"\u00edz":      `'\ufb5e'`,
	"\u00f3s":      `'\ufb5f'`,
	"\u00f3sx":     `'\ufb60'`,
	"\u00f3sz":     `'\ufb61'`,
	"\u00f3x":      `'\ufb62'`,
	"\u0105,,":     `'\ufb7b'`,
	"\u0107,,":     `'\ufb7c'`,
	"\u0119,,":     `'\ufb7d'`,
	"\u0142,,":     `'\ufb7e'`,
	"\u0144,,":     `'\ufb7f'`,
	"\u017a,,":     `'\ufb80'`,
	"\u017c,,":     `'\ufb81'`,
	"\ue001\ue014": `'\u2469'`,
	"\ue001\ue015": `'\u246a'`,
	"\ue001\ue016": `'\u246b'`,
	"\ue001\ue017": `'\u246c'`,
	"\ue001\ue018": `'\u246d'`,
	"\ue001\ue019": `'\u246e'`,
	"\ue001\ue01a": `'\u246f'`,
	"\ue001\ue01b": `'\u2470'`,
	"\ue001\ue01c": `'\u2471'`,
	"\ue001\ue01d": `'\u2472'`,
	"\ue002\ue014": `'\u2473'`,
	"u,,":          `'\ufb7a'`,
	"uv":           `'\ufb63'`,
	"ux":           `'\ufb64'`,
	"x,,":          `'\ufb82'`,
	"yf":           `'\ufb84'`,
	"Yj":           `'\ufb1f'`,
	"yj":           `'\ufb85'`,
	"yp":           `'\ufb65'`,
	"y,,":          `'\ufb83'`,
	"yx":           `'\ufb66'`,
	"yz":           `'\ufb67'`,
	"z,,":          `'\ufb86'`,
}
