// render-fonts renders fonts with the given message in all imported fonts
// to PNG files in order to test out the font rendering.
package main

// To import any desired fonts, import them below:
// _ "github.com/gmlewis/go-fonts/fonts/ubuntumonoregular"
// _ "github.com/gmlewis/go-fonts/fonts/znikomitno24"
// etc.

import (
	"flag"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"

	. "github.com/gmlewis/go-fonts/fonts"
	_ "github.com/gmlewis/go-fonts/fonts/latoregular"
)

var (
	all = flag.Bool("all", false, "All renders all glyphs and overrides -msg")
	msg = flag.String("msg", `0123456789
ABCDEFGHIJKLM
NOPQRSTUVWXYZ
abcdefghijklm
nopqrstuvwxyz
~!@#$%^&*()-_=/?
+[]{}\|;':",.<>`, "Message to write to Gerber file silkscreen")
	center = flag.Bool("center", false, "Center justify all text")
	width  = flag.Int("width", 800, "Image width")
	height = flag.Int("height", 800, "Image height")
	dxf    = flag.String("dxf", "out.dxf", "Output DXF filename")
	out    = flag.String("out", "out.png", "Output image filename")
	rot    = flag.Float64("rot", 0, "Rotate message by this number of degrees")
)

func main() {
	flag.Parse()

	for name, font := range Fonts {
		message := *msg
		if *all {
			var glyphs []rune
			for g := range font.Glyphs {
				glyphs = append(glyphs, g)
			}
			sort.Slice(glyphs, func(a, b int) bool { return glyphs[a] < glyphs[b] })

			lineLength := int(0.5 + math.Sqrt(float64(len(glyphs))))
			var lines []string
			for len(glyphs) > 0 {
				end := lineLength
				if end > len(glyphs) {
					end = len(glyphs)
				}
				lines = append(lines, string(glyphs[0:end]))
				glyphs = glyphs[end:]
			}
			message = strings.Join(lines, "\n")
		}

		var opts *TextOpts
		if *rot != 0.0 {
			opts = &TextOpts{
				Rotate: *rot * math.Pi / 180.0,
			}
		}

		var render *Render
		if *center {
			if opts == nil {
				opts = &TextOpts{}
			}
			opts.XAlign, opts.YAlign = XCenter, YTop
			var lines []*Render
			lastY := 0.0
			for _, line := range strings.Split(message, "\n") {
				r, err := Text(0, lastY, 1.0, 1.0, line, name, opts)
				if err != nil {
					log.Fatal(err)
				}
				lastY = r.MBB.Min[1] + font.Descent/font.UnitsPerEm
				lines = append(lines, r)
			}
			render = Merge(lines...)
		} else {
			var err error
			render, err = Text(0, 0, 1.0, 1.0, message, name, opts)
			if err != nil {
				log.Fatal(err)
			}
		}

		if *out != "" {
			if err := render.SavePNG(*out, *width, *height); err != nil {
				log.Fatal(err)
			}
		}
		if *dxf != "" {
			if err := render.SaveDXF(*dxf, *width, *height); err != nil {
				log.Fatal(err)
			}
		}
	}

	fmt.Println("Done.")
}
