package fonts

import (
	"fmt"
	"os"
	"strings"

	svg "github.com/gmlewis/ponoko2d/float"
)

const (
	defaultStyle = "stroke:#000000;stroke-opacity:1;fill:none"
)

func (t *Render) SaveSVG(fontName, filename string, scale float64) error {
	font, ok := Fonts[fontName]
	if !ok {
		return fmt.Errorf("cannot find font %q", fontName)
	}
	fsf := 1.0 / font.UnitsPerEm

	w, err := os.Create(filename)
	if err != nil {
		return err
	}

	width := (t.MBB.Max[0] - t.MBB.Min[0]) * scale
	height := (t.MBB.Max[1] - t.MBB.Min[1]) * scale

	s := svg.New(w)
	s.Start(width, height)
	s.ScaleXY(1, -1)        // flip vertically
	s.Translate(0, -height) // and re-center
	s.ScaleXY(scale, scale)

	for _, gi := range t.Info {
		if gi.Glyph == '\n' {
			continue
		}

		s.Translate(gi.X, gi.Y)
		s.ScaleXY(fsf, fsf)
		path, err := font.SVGPath(gi)
		if err != nil {
			return err
		}
		s.Path(path, defaultStyle)
		s.Gend() // s.ScaleXY
		s.Gend() // s.Translate
	}

	s.Gend() // s.ScaleXY(scale,scale)
	s.Gend() // s.Translate(0,-height)
	s.Gend() // s.ScaleXY(1,-1)
	s.End()  // s.Start(width,height)

	return w.Close()
}

func (f *Font) SVGPath(gi *GlyphInfo) (string, error) {
	if gi.Glyph == '\n' {
		return "", nil
	}

	g, ok := f.Glyphs[gi.Glyph]
	if !ok {
		return "", fmt.Errorf("glyph %q not found", gi.Glyph)
	}

	var parts []string
	for _, ps := range g.PathSteps {
		lineParts := []string{fmt.Sprintf("%c", ps.C)}
		for _, p := range ps.P {
			lineParts = append(lineParts, fmt.Sprintf("%v", p))
		}
		parts = append(parts, strings.Join(lineParts, " "))
	}
	return strings.Join(parts, ""), nil
}
