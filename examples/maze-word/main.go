// maze-word renders a word in white-on-black.
package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"

	_ "github.com/gmlewis/go-fonts-m/fonts/modak"
	. "github.com/gmlewis/go-fonts/fonts"
)

var (
	fontName = flag.String("font", "modak", "Font to use")
	word     = flag.String("w", "", "Word to render")
	width    = flag.Int("width", 850, "Image width")
	height   = flag.Int("height", 1100, "Image height")
	outFmt   = flag.String("out", "%v-%v-%v.png", "Output image filename format string")
	shift    = flag.Float64("shift", 0.1, "Percentage of glyph width to shift")
)

func main() {
	flag.Parse()

	if *word == "" {
		*word = "Sample"
	}

	m1, err := Text(0, 0, 1, 1, *word, *fontName, nil)
	check(err)
	m1.Foreground = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	m1.Background = color.RGBA{R: 0, G: 0, B: 0, A: 255}

	shiftAmounts := map[int]float64{}
	var totalShift float64
	for i, gi := range m1.Info {
		if i == 0 {
			continue
		}
		log.Printf("gi #%v '%c': mbb=%v", i+1, gi.Glyph, gi.MBB)
		sw := -0.1 * (gi.MBB.Max[0] - gi.MBB.Min[0])
		totalShift += sw
		shiftAmounts[i] = totalShift
		log.Printf("shiftAmounts[%v] = %v", i, totalShift)
	}
	for i, p := range m1.Polygons {
		log.Printf("poly #%v: rune[%v], mbb=%v", i+1, p.RuneIndex, p.MBB)
		sw := shiftAmounts[p.RuneIndex]
		p.MBB.Min[0] += sw
		p.MBB.Min[1] += sw
		for j := range p.Pts {
			p.Pts[j][0] += sw
		}
	}

	all := []*Render{
		m1,
	}

	out := fmt.Sprintf(*outFmt, *width, *word, *fontName)
	if err := SavePNG(out, *width, *height, all...); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done.")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
