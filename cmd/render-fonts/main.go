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

	"github.com/fogleman/gg"
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
	width  = flag.Int("width", 800, "Image width")
	height = flag.Int("height", 800, "Image height")
	out    = flag.String("out", "out.png", "Output image filename")
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

		render, err := Text(0, 0, 1.0, 1.0, message, name)
		if err != nil {
			log.Fatal(err)
		}

		scale := float64(*width) / (render.Xmax - render.Xmin)
		if yScale := float64(*height) / (render.Ymax - render.Ymin); yScale < scale {
			scale = yScale
			*width = int(0.5 + scale*(render.Xmax-render.Xmin))
		} else {
			*height = int(0.5 + scale*(render.Ymax-render.Ymin))
		}
		log.Printf("MBB: (%.2f,%.2f)-(%.2f,%.2f), scale=%.2f, size=(%v,%v)", render.Xmin, render.Ymin, render.Xmax, render.Ymax, scale, *width, *height)

		dc := gg.NewContext(*width, *height)
		dc.SetRGB(1, 1, 1)
		dc.Clear()
		for _, poly := range render.Polygons {
			if poly.Dark {
				dc.SetRGB(0, 0, 0)
			} else {
				dc.SetRGB(1, 1, 1)
			}
			for i, pt := range poly.Pts {
				if i == 0 {
					dc.MoveTo(scale*(pt.X-render.Xmin), float64(*height)-scale*(pt.Y-render.Ymin))
				} else {
					dc.LineTo(scale*(pt.X-render.Xmin), float64(*height)-scale*(pt.Y-render.Ymin))
				}
			}
			dc.Fill()
		}
		dc.SavePNG(*out)
	}

	fmt.Println("Done.")
}
