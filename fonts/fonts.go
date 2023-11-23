// Package fonts provides a collection of open source fonts converted to Go.
//
// To use the fonts in your own program, import the main package and
// import the font(s) that you want.
//
// For example:
//
// import (
//
//	"github.com/gmlewis/go-fonts/fonts"
//	_ "github.com/gmlewis/go-fonts/fonts/znikomitno24"
//
// )
//
//	func main() {
//	  polys := fonts.Text(x, y, xs, ys, "znikomitno24")
//	  //...
//	}
//
// Default units are in "em"s which typically represent the width of the
// character "M" in the font. Note that positive X is to the right
// and positive Y is up.
//
// Each polygon is either "dark" or "clear". Dark polygons should be
// rendered before clear ones and should be returned in a natural drawing
// order.
package fonts

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"io"
	"log"
	"unicode/utf8"

	"github.com/gmlewis/go-fonts/pb/glyphs"
	"google.golang.org/protobuf/proto"
)

// Font represents a webfont.
//
// Each font uses its own native units (typically not "ems") that are later
// scaled to "ems" when rendered.
type Font struct {
	// ID is the name used to identify the font.
	ID string
	// HorizAdvX is the default font units to advance per glyph.
	HorizAdvX float64
	// UnitsPerEm is the number of font units per "em".
	UnitsPerEm float64
	// Ascent is the height of the font above the baseline.
	Ascent float64
	// Descent is the negative vertical distance below the baseline
	// of the font.
	Descent float64
	// MissingHorizAdvX is the amount of font units to advance
	// in the case of missing glyphs.
	MissingHorizAdvX float64
	// Glyphs is a map of the available glyphs, mapped by rune.
	Glyphs map[rune]*Glyph
}

// Glyph represents an individual character of the webfont data.
type Glyph struct {
	// HorizAdvX is the number of font units to advance for this glyph.
	HorizAdvX float64
	// Unicode is the rune representing this glyph.
	Unicode rune
	// GerberLP is a string of "d" (for "dark") and "c" (for "clear")
	// representing the nature of each subsequent font curve subpath
	// contained within PathSteps. Its length matches the number of
	// subpaths in PathSteps (Starting from 'M' or 'm' path commands.)
	GerberLP string
	// PathSteps represents the SVG commands that define the glyph.
	PathSteps []*PathStep
	// MBB represents the minimum bounding box of the glyph in native units.
	MBB MBB
}

// PathStep represents a single subpath command.
//
// There are 20 possible commands, broken up into 6 types,
// with each command having an "absolute" (upper case) and
// a "relative" (lower case) version.
//
// Note that not all subpath types are currently supported.
// Only the ones needed for the provided fonts have been
// implemented.
//
// See https://developer.mozilla.org/en-US/docs/Web/SVG/Attribute/d
// for more details.
//
// MoveTo: M, m
// LineTo: L, l, H, h, V, v
// Cubic Bézier Curve: C, c, S, s
// Quadratic Bézier Curve: Q, q, T, t
// Elliptical Arc Curve: A, a
// ClosePath: Z, z
type PathStep struct {
	C byte      // C is the command.
	P []float64 // P are the parameters of the command.
}

// Fonts is a map of all the available fonts.
//
// The map is initialized at runtime by `init` functions
// in order to reduce the overall initial compile time
// of the package.
var Fonts = map[string]*Font{}

// InitFromFontData is a workaround for a Go compiler error when
// compiling large map literals. The data is marshaled to a string
// from a protobuf then base64 encoded into a font package source file.
//
// This function decodes the base64 data, then unmarshals the protobuf
// and populates the glyphs into the font. Fortunately, this happens
// extremely quickly in the package's `init` function.
func InitFromFontData(font *Font, fontData string) {
	data, err := base64.StdEncoding.DecodeString(fontData)
	if err != nil {
		log.Fatalf("unable to base64 decode %v fontData: %v", font.ID, err)
	}
	b := bytes.NewBuffer(data)
	r, err := zlib.NewReader(b)
	if err != nil {
		log.Fatalf("zlib.NewReader: %v", err)
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		log.Fatalf("io.Copy: %v", err)
	}
	if err := r.Close(); err != nil {
		log.Fatalf("zlib.Close: %v", err)
	}

	gs := &glyphs.Glyphs{}
	if err := proto.Unmarshal(buf.Bytes(), gs); err != nil {
		log.Fatalf("unable to unmarshal %v proto data: %v", font.ID, err)
	}
	for _, glyph := range gs.Glyphs {
		r, _ := utf8.DecodeRuneInString(glyph.Unicode)
		var pathSteps []*PathStep
		for _, ps := range glyph.PathSteps {
			pathSteps = append(pathSteps, &PathStep{
				C: byte(ps.C),
				P: ps.P,
			})
		}
		mbb := glyph.GetMbb()
		font.Glyphs[r] = &Glyph{
			HorizAdvX: glyph.HorizAdvX,
			Unicode:   r,
			GerberLP:  glyph.GerberLP,
			PathSteps: pathSteps,
			MBB: MBB{
				Min: Pt{mbb.GetXmin(), mbb.GetYmin()},
				Max: Pt{mbb.GetXmax(), mbb.GetYmax()},
			},
		}
	}
}
