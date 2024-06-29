// -*- compile-command: "go run main.go 'Sissy' '&Joe'"; -*-

// make-puzzle generates a text puzzle for 3D-printing from a string.
//
// Usage:
//
// make-puzzle 'Sissy' '&Joe'
package main

import (
	"flag"
	"log"
	"strings"

	_ "github.com/gmlewis/go-fonts-b/fonts/baloo"
	. "github.com/gmlewis/go-fonts/fonts"
)

var (
	filename = flag.String("filename", "out.dxf", "DXF file to write output")
	overlap  = flag.Float64("overlap", 0.025, "Overlap for adjacent glyphs")
	scale    = flag.Float64("scale", 288, "Font scale")
)

func main() {
	flag.Parse()

	lines := flag.Args()
	t, err := Text(0, 0, 1, 1, strings.Join(lines, "\n"), "baloo", nil)
	must(err)

	squishTogether(t)

	must(t.SaveDXF(*filename, 1.0))
	must(t.SaveSVG("baloo", strings.Replace(*filename, ".dxf", ".svg", -1), *scale))

	log.Printf("Done.")
}

func squishTogether(t *Render) {
	if len(t.Info) < 1 {
		log.Fatal("no text to render")
	}

	var xi int
	var lastX float64
	for i, gi := range t.Info {
		log.Printf("i=%v, %q: (%.2f,%0.2f) w=%0.2f, mbb=%v", i, gi.Glyph, gi.X, gi.Y, gi.Width, gi.MBB)
		if xi == 0 {
			xi++
			lastX = gi.MBB.Max[0]
			continue
		}
		if gi.Glyph == '\n' {
			xi = 0
			continue
		}

		xi++
		dx := gi.MBB.Min[0] - lastX
		moveX := -dx - *overlap
		t.MoveGlyph(i, moveX, 0)
		lastX = gi.MBB.Max[0]
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
