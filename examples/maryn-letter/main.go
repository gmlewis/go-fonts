// maryn-letter renders a letter to Maryn.
package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"

	. "github.com/gmlewis/go-fonts/fonts"
	_ "github.com/gmlewis/go-fonts/fonts/latoregular"
	_ "github.com/gmlewis/go-fonts/fonts/printersornamentsone"
	_ "github.com/gmlewis/go-fonts/fonts/scriptinapro"
	_ "github.com/gmlewis/go-fonts/fonts/sofia_regular"
	_ "github.com/gmlewis/go-fonts/fonts/spirax_regular"
	_ "github.com/gmlewis/go-fonts/fonts/tangerine_bold"
	_ "github.com/gmlewis/go-fonts/fonts/topsecret_bold"
	_ "github.com/gmlewis/go-fonts/fonts/typemymusic_notation"
)

var (
	width  = flag.Int("width", 8500, "Image width")
	height = flag.Int("height", 11000, "Image height")
	out    = flag.String("out", "maryn-letter.png", "Output image filename")
)

const (
	carrFont = "carrelectronicdingbats"
	textFont = "latoregular"

	dy      = 1
	upArrow = "⇧"

	msg1 = `Hi Maryn,

    Recently I have been working on a fun project that allows me to make
printed circuit boards to try out experiments. In the process, I needed
to be able to draw letters on the boards so I wrote some code that let me
use different fonts (which are collections of letters that use the same
style). This paragraph uses a font called "Sofia Regular" which I thought
you might like.`
	msg2 = `
    This letter was written to you using that code. I'm changing the font
so that you can see different styles of writing. This font is called
"Spirax Regular". I didn't actually need different font colors for my
printed circuit boards, but I thought you might like to see the letters in
different colors, so I modified the code to let me draw with colors.`
	msg3 = `
    This font is called "Top Secret Bold". It only has
upper-case letters in it, unfortunately.`
	msg4 = `
    On the second page, I'll print out a copy of the code that I used to
"render" (which just means "to draw") this letter. The programming
language that I used to write the code is called Go. It is one of my
favorite programming languages. I enjoy learning different programming
languages just for fun. Some people probably think that is weird! Right
now I'm reading about a language called "Nim". This is a fun font
called "Tangerine Bold".`
	msg5 = `
    There are even fonts that render musical notes, like this one which
is called "Type My Music Notation":`
	msg6 = `
    12345678`
	msg7 = `
    Some fonts just have pretty pictures. Here is one of those that
is called "Printers Ornaments One":`
	msg8 = `
abcdefghijklmnopqrst`
	msg9 = `
    If you ever want to learn more about programming, I'm always happy
to talk about it, or anything else you want to talk about.`
	msg10 = `Love,`
)

func main() {
	flag.Parse()

	m1, err := Text(0, 0, 1, 1, msg1, "sofia_regular", nil)
	check(err)
	m1.Foreground = color.RGBA{R: 255, G: 105, B: 180, A: 255}
	m2, err := Text(0, m1.MBB.Min[1]-dy, 1, 1, msg2, "spirax_regular", &TopLeft)
	check(err)
	m2.Foreground = color.RGBA{R: 200, G: 0, B: 200, A: 255}
	m3, err := Text(0, m2.MBB.Min[1]-dy, 1, 1, msg3, "topsecret_bold", &TopLeft)
	check(err)
	m3.Foreground = color.RGBA{R: 200, G: 0, B: 0, A: 255}
	m4, err := Text(0, m3.MBB.Min[1]-dy, 1.7, 1.7, msg4, "tangerine_bold", &TopLeft)
	check(err)
	m4.Foreground = color.RGBA{R: 0, G: 0, B: 200, A: 255}
	m5, err := Text(0, m4.MBB.Min[1]-dy, 1, 1, msg5, "latoregular", &TopLeft)
	check(err)
	m5.Foreground = color.RGBA{R: 200, G: 110, B: 0, A: 255}
	m6, err := Text(0, m5.MBB.Min[1]-dy, 10, 10, msg6, "typemymusic_notation", &TopLeft)
	check(err)
	m7, err := Text(0, m6.MBB.Min[1]-dy, 1, 1, msg7, "latoregular", &TopLeft)
	check(err)
	m7.Foreground = color.RGBA{R: 0, G: 100, B: 0, A: 255}
	m8, err := Text(0, m7.MBB.Min[1]-dy, 1.25, 1.25, msg8, "printersornamentsone", &TopLeft)
	check(err)
	m8.Foreground = color.RGBA{R: 0, G: 100, B: 0, A: 255}
	m9, err := Text(0, m8.MBB.Min[1]-dy, 1, 1, msg9, "latoregular", &TopLeft)
	check(err)
	m10, err := Text(m6.MBB.Max[0], m9.MBB.Min[1]-dy, 2.5, 2.5, msg10, "scriptinapro", &TopRight)
	check(err)

	all := []*Render{
		m1, m2, m3, m4, m5, m6, m7, m8, m9, m10,
	}

	if err := SavePNG(*out, *width, *height, all...); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done.")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
