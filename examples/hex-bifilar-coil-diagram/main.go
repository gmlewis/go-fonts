// hex-bifilar-coil-diagram renders a schematic diagram
// representing the wiring of the hex bifilar coil design
// from https://github.com/gmlewis/go-gerber/tree/master/examples/hex-bifilar-coil .
package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	. "github.com/gmlewis/go-fonts/fonts"
	_ "github.com/gmlewis/go-fonts/fonts/carrelectronicdingbats"
	_ "github.com/gmlewis/go-fonts/fonts/freeserif"
	_ "github.com/gmlewis/go-fonts/fonts/latoregular"
)

var (
	width  = flag.Int("width", 800, "Image width")
	height = flag.Int("height", 800, "Image height")
	out    = flag.String("out", "hex-bifilar-coil-diagram.png", "Output image filename")
)

const (
	carrFont = "carrelectronicdingbats"
	textFont = "latoregular"

	dx      = 0.02
	upArrow = "⇧"
)

func main() {
	flag.Parse()

	opts := &TextOpts{Rotate: math.Pi / 2}
	coil1, err := Text(0, 0, 1.0, 1.0, "}", carrFont, opts)
	check(err)
	log.Printf("coil1=%v", coil1.MBB)
	coil2, err := Text(coil1.MBB.Max[0]-coil1.MBB.Min[0]-dx, 0, 1.0, 1.0, "}", carrFont, opts)
	check(err)
	log.Printf("coil2=%v", coil2.MBB)
	coil3, err := Text(coil2.MBB.Max[0]-coil1.MBB.Min[0]-dx, 0, 1.0, 1.0, "}", carrFont, opts)
	check(err)
	log.Printf("coil3=%v", coil3.MBB)
	coil4, err := Text(coil3.MBB.Max[0]-coil1.MBB.Min[0]-dx, 0, 1.0, 1.0, "}", carrFont, opts)
	check(err)
	log.Printf("coil4=%v", coil4.MBB)
	coil5, err := Text(coil4.MBB.Max[0]-coil1.MBB.Min[0]-dx, 0, 1.0, 1.0, "}", carrFont, opts)
	check(err)
	log.Printf("coil5=%v", coil5.MBB)
	coil6, err := Text(coil5.MBB.Max[0]-coil1.MBB.Min[0]-dx, 0, 1.0, 1.0, "}", carrFont, opts)
	check(err)
	log.Printf("coil6=%v", coil6.MBB)

	opts = &TopCenter
	ts := 0.25
	text3L, err := Text(coil1.MBB.Min[0]+dx, coil1.MBB.Min[1], ts, ts, "3L", textFont, opts)
	check(err)
	textBL, err := Text(coil2.MBB.Min[0]+dx, coil2.MBB.Min[1], ts, ts, "BL\n4L", textFont, opts)
	check(err)
	text5L, err := Text(coil3.MBB.Min[0]+dx, coil3.MBB.Min[1], ts, ts, "TL\n5L", textFont, opts)
	check(err)
	text3R, err := Text(coil4.MBB.Min[0]+dx, coil4.MBB.Min[1], ts, ts, "3R\n2L", textFont, opts)
	check(err)
	textBR, err := Text(coil5.MBB.Min[0]+dx, coil5.MBB.Min[1], ts, ts, "4R\nBR", textFont, opts)
	check(err)
	text5R, err := Text(coil6.MBB.Min[0]+dx, coil6.MBB.Min[1], ts, ts, "5R\nTR", textFont, opts)
	check(err)
	text2R, err := Text(coil6.MBB.Max[0]+dx, coil6.MBB.Min[1], ts, ts, "2R", textFont, opts)
	check(err)

	// smaller font for inner connection points
	ts = 0.125
	up1, err := Text(0.5*(coil1.MBB.Min[0]+coil1.MBB.Max[0]), coil1.MBB.Min[1], ts, ts, upArrow, "freeserif", opts)
	check(err)
	text3L4L, err := Text(0.5*(coil1.MBB.Min[0]+coil1.MBB.Max[0]), up1.MBB.Min[1], ts, ts, "3L/4L", textFont, opts)
	check(err)
	up2, err := Text(0.5*(coil2.MBB.Min[0]+coil2.MBB.Max[0]), coil1.MBB.Min[1], ts, ts, upArrow, "freeserif", opts)
	check(err)
	textTLBL, err := Text(0.5*(coil2.MBB.Min[0]+coil2.MBB.Max[0]), up2.MBB.Min[1], ts, ts, "TL/BL", textFont, opts)
	check(err)
	up3, err := Text(0.5*(coil3.MBB.Min[0]+coil3.MBB.Max[0]), coil1.MBB.Min[1], ts, ts, upArrow, "freeserif", opts)
	check(err)
	text2L5L, err := Text(0.5*(coil3.MBB.Min[0]+coil3.MBB.Max[0]), up3.MBB.Min[1], ts, ts, "2L/5L", textFont, opts)
	check(err)
	up4, err := Text(0.5*(coil4.MBB.Min[0]+coil4.MBB.Max[0]), coil1.MBB.Min[1], ts, ts, upArrow, "freeserif", opts)
	check(err)
	text3R4R, err := Text(0.5*(coil4.MBB.Min[0]+coil4.MBB.Max[0]), up4.MBB.Min[1], ts, ts, "3R/4R", textFont, opts)
	check(err)
	up5, err := Text(0.5*(coil5.MBB.Min[0]+coil5.MBB.Max[0]), coil1.MBB.Min[1], ts, ts, upArrow, "freeserif", opts)
	check(err)
	textTRBR, err := Text(0.5*(coil5.MBB.Min[0]+coil5.MBB.Max[0]), up5.MBB.Min[1], ts, ts, "TR/BR", textFont, opts)
	check(err)
	up6, err := Text(0.5*(coil6.MBB.Min[0]+coil6.MBB.Max[0]), coil1.MBB.Min[1], ts, ts, upArrow, "freeserif", opts)
	check(err)
	text2R5R, err := Text(0.5*(coil6.MBB.Min[0]+coil6.MBB.Max[0]), up6.MBB.Min[1], ts, ts, "2R/5R", textFont, opts)
	check(err)

	coils := Merge(
		coil1, coil2, coil3, coil4, coil5, coil6,
		text3L, textBL, text5L, text3R, textBR, text5R, text2R,
		up1, up2, up3, up4, up5, up6,
		text3L4L, textTLBL, text2L5L, text3R4R, textTRBR, text2R5R,
	)

	if err := coils.SavePNG(*out, *width, *height); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done.")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
