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

	dx = 0.02
	ts = 0.25
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
	text3L, err := Text(coil1.MBB.Min[0]+dx, coil1.MBB.Min[1], ts, ts, "3L", textFont, opts)
	check(err)
	textBL, err := Text(coil2.MBB.Min[0]+dx, coil2.MBB.Min[1], ts, ts, "BL", textFont, opts)
	check(err)
	text5L, err := Text(coil3.MBB.Min[0]+dx, coil3.MBB.Min[1], ts, ts, "5L", textFont, opts)
	check(err)
	text3R, err := Text(coil4.MBB.Min[0]+dx, coil4.MBB.Min[1], ts, ts, "3R", textFont, opts)
	check(err)
	textBR, err := Text(coil5.MBB.Min[0]+dx, coil5.MBB.Min[1], ts, ts, "BR", textFont, opts)
	check(err)
	text5R, err := Text(coil6.MBB.Min[0]+dx, coil6.MBB.Min[1], ts, ts, "5R", textFont, opts)
	check(err)
	text2R, err := Text(coil6.MBB.Max[0]+dx, coil6.MBB.Min[1], ts, ts, "2R", textFont, opts)
	check(err)

	coils := Merge(
		coil1, coil2, coil3, coil4, coil5, coil6,
		text3L, textBL, text5L, text3R, textBR, text5R, text2R,
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
