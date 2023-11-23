package main

import (
	"testing"

	"github.com/gmlewis/go-fonts/webfont"
	"github.com/google/go-cmp/cmp"
)

func String(s string) *string { return &s }

func TestProcessor_Letter_a(t *testing.T) {
	g := &webfont.Glyph{
		Unicode:  String("a"),
		GerberLP: String("dc"),
		D: String(`M53 369q6 170 210 170q97 0 147 -36.5t50 -106.5v-308q0 -41 45 -41q9 0 18 2v-63q-35 -9 -57 -9q-40 0 -60.5 18.5t-25.5 58.5q-84 -77 -178 -77q-79 0 -125.5 42t-46.5 113q0 23 4.5 42t10 33.5t19.5 27t23.5 20.5t32 15.5t34.5 11t42 8.5t44 7t50 7q49 6 68 19t19 39
v22q0 38 -30.5 58t-86.5 20q-58 0 -88 -22.5t-35 -70.5h-84zM220 50q69 0 113 36.5t44 78.5v94q-25 -12 -75.5 -20t-87.5 -14t-67 -28.5t-30 -62.5t27 -62t76 -22z`),
	}

	p := &processor{glyphs: map[string]*glyphT{}}
	p.NewGlyph(g)
	p.MoveTo(g, "M", 53, 369)
	p.QuadraticTo(g, "q", 59, 539, 263, 539)
	p.QuadraticTo(g, "q", 360, 539, 410, 502.5)
	p.QuadraticTo(g, "t", 460, 466, 460, 396)
	p.LineTo(g, "v", 460, 88)
	p.QuadraticTo(g, "q", 460, 47, 505, 47)
	p.QuadraticTo(g, "q", 514, 47, 523, 49)
	p.LineTo(g, "v", 523, -14)
	p.QuadraticTo(g, "q", 488, -23, 466, -23)
	p.QuadraticTo(g, "q", 426, -23, 405.5, -4.5)
	p.QuadraticTo(g, "t", 385, 14, 380, 54)
	p.QuadraticTo(g, "q", 296, -23, 202, -23)
	p.QuadraticTo(g, "q", 123, -23, 76.5, 19)
	p.QuadraticTo(g, "t", 30, 61, 30, 132)
	p.QuadraticTo(g, "q", 30, 155, 34.5, 174)
	p.QuadraticTo(g, "t", 39, 193, 44.5, 207.5)
	p.QuadraticTo(g, "t", 50, 222, 64, 234.5)
	p.QuadraticTo(g, "t", 78, 247, 87.5, 255)
	p.QuadraticTo(g, "t", 97, 263, 119.5, 270.5)
	p.QuadraticTo(g, "t", 142, 278, 154, 281.5)
	p.QuadraticTo(g, "t", 166, 285, 196, 290)
	p.QuadraticTo(g, "t", 226, 295, 240, 297)
	p.QuadraticTo(g, "t", 254, 299, 290, 304)
	p.QuadraticTo(g, "q", 339, 310, 358, 323)
	p.QuadraticTo(g, "t", 377, 336, 377, 362)
	p.LineTo(g, "v", 377, 384)
	p.QuadraticTo(g, "q", 377, 422, 346.5, 442)
	p.QuadraticTo(g, "t", 316, 462, 260, 462)
	p.QuadraticTo(g, "q", 202, 462, 172, 439.5)
	p.QuadraticTo(g, "t", 142, 417, 137, 369)
	p.LineTo(g, "h", 53, 369)
	p.MoveTo(g, "M", 220, 50)
	p.QuadraticTo(g, "q", 289, 50, 333, 86.5)
	p.QuadraticTo(g, "t", 377, 123, 377, 165)
	p.LineTo(g, "v", 377, 259)
	p.QuadraticTo(g, "q", 352, 247, 301.5, 239)
	p.QuadraticTo(g, "t", 251, 231, 214, 225)
	p.QuadraticTo(g, "t", 177, 219, 147, 196.5)
	p.QuadraticTo(g, "t", 117, 174, 117, 134)
	p.QuadraticTo(g, "t", 117, 94, 144, 72)
	p.QuadraticTo(g, "t", 171, 50, 220, 50)
	p.ProcessGlyph('a', g)

	glyph, ok := p.glyphs["a"]
	if !ok {
		t.Fatalf("p.glyphs missing 'a': %+v", p.glyphs)
	}
	if p.current != glyph {
		t.Fatalf("p.current != glyph")
	}

	wantD := "M53 369" +
		"Q59 539 263 539" +
		"Q360 539 410 502.5" +
		"Q460 466 460 396" +
		"L460 88" +
		"Q460 47 505 47" +
		"Q514 47 523 49" +
		"L523 -14" +
		"Q488 -23 466 -23" +
		"Q426 -23 405.5 -4.5" +
		"Q385 14 380 54" + // cut0Idx=10
		"L333 86.5" + // jump over to face[1] - after[cutIdx=1]
		"Q377 123 377 165L377 259Q352 247 301.5 239Q251 231 214 225Q177 219 147 196.5Q117 174 117 134Q117 94 144 72Q171 50 220 50" + // remainder of face[1] after[2:] (after[cuIdx+1:])
		// completely SKIP face[1] after[0] - "M220 50"
		"Q289 50 333 86.5" + // cutIdx=1
		"L380 54" + // jump back to face[0] - after[cut0Idx=10]
		"Q296 -23 202 -23Q123 -23 76.5 19Q30 61 30 132Q30 155 34.5 174Q39 193 44.5 207.5Q50 222 64 234.5Q78 247 87.5 255Q97 263 119.5 270.5Q142 278 154 281.5Q166 285 196 290Q226 295 240 297Q254 299 290 304Q339 310 358 323Q377 336 377 362L377 384Q377 422 346.5 442Q316 462 260 462Q202 462 172 439.5Q142 417 137 369L53 369" + // remainder of face[0] after[11:] (after[cut0Idx+1:])
		"Z" // terminate the new face.

	// t.Logf("wantD:\n%v", wantD)
	// t.Logf("gotD:\n%v", glyph.d)
	if diff := cmp.Diff(wantD, glyph.d); diff != "" {
		t.Errorf("d mismatch (-want +got):\n%v", diff)
	}

	wantFaces := []*faceT{
		{
			absCmds: []string{
				"M53 369",
				"Q59 539 263 539",
				"Q360 539 410 502.5",
				"Q460 466 460 396",
				"L460 88",
				"Q460 47 505 47",
				"Q514 47 523 49",
				"L523 -14",
				"Q488 -23 466 -23",
				"Q426 -23 405.5 -4.5",
				"Q385 14 380 54", // cut0Idx=10
				"Q296 -23 202 -23",
				"Q123 -23 76.5 19",
				"Q30 61 30 132",
				"Q30 155 34.5 174",
				"Q39 193 44.5 207.5",
				"Q50 222 64 234.5",
				"Q78 247 87.5 255",
				"Q97 263 119.5 270.5",
				"Q142 278 154 281.5",
				"Q166 285 196 290",
				"Q226 295 240 297",
				"Q254 299 290 304",
				"Q339 310 358 323",
				"Q377 336 377 362",
				"L377 384",
				"Q377 422 346.5 442",
				"Q316 462 260 462",
				"Q202 462 172 439.5",
				"Q142 417 137 369",
				"L53 369",
			},
			// before: []vec2{{0, 0}, {460, 396}, {523, 49}, {377, 362}, {137, 369}},
			after: []vec2{
				{53, 369}, {263, 539}, {410, 502.5}, {460, 396}, {460, 88}, {505, 47},
				{523, 49}, {523, -14}, {466, -23}, {405.5, -4.5},
				{380, 54}, // cut0Idx=10
				{202, -23}, {76.5, 19}, {30, 132}, {34.5, 174}, {44.5, 207.5}, {64, 234.5}, {87.5, 255},
				{119.5, 270.5}, {154, 281.5}, {196, 290}, {240, 297}, {290, 304}, {358, 323},
				{377, 362}, {377, 384}, {346.5, 442}, {260, 462}, {172, 439.5}, {137, 369}, {53, 369},
			},
		},
		{
			absCmds: []string{
				"M220 50",
				"Q289 50 333 86.5", // cutIdx=1
				"Q377 123 377 165",
				"L377 259",
				"Q352 247 301.5 239",
				"Q251 231 214 225",
				"Q177 219 147 196.5",
				"Q117 174 117 134",
				"Q117 94 144 72",
				"Q171 50 220 50",
			},
			after: []vec2{
				{220, 50}, {333, 86.5}, {377, 165}, {377, 259}, {301.5, 239},
				{214, 225}, {147, 196.5}, {117, 134}, {144, 72}, {220, 50},
			},
			cut0Idx: 10,
			cutIdx:  1,
		},
	}

	if got, want := len(glyph.faces), len(wantFaces); got != want {
		t.Fatalf("processor = %v faces, want %v", got, want)
	}

	for faceIdx, face := range glyph.faces {
		if diff := cmp.Diff(face, wantFaces[faceIdx], cmp.AllowUnexported(faceT{})); diff != "" {
			t.Errorf("face[%v] mismatch (-want +got):\n%v", faceIdx, diff)
		}
	}
}
