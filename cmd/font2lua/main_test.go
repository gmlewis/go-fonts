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

	wantD := "M53 369q6 170 210 170q97 0 147 -36.5t50 -106.5v-308q0 -41 45 -41q9 0 18 2v-63q-35 -9 -57 -9q-40 0 -60.5 18.5t-25.5 58.5q-84 -77 -178 -77q-79 0 -125.5 42t-46.5 113q0 23 4.5 42t10 33.5t19.5 27t23.5 20.5t32 15.5t34.5 11t42 8.5t44 7" +
		"L301.5 239" + // jump over to face1
		"t-87.5 -14t-67 -28.5t-30 -62.5t27 -62t76 -22" + // rest of face1 - strip "z"!
		"q69 0 113 36.5t44 78.5v94q-25 -12 -75.5 -20" + // start of face1
		"L240 297" + // jump back to face0
		"t50 7q49 6 68 19t19 39\nv22q0 38 -30.5 58t-86.5 20q-58 0 -88 -22.5t-35 -70.5h-84z" // rest of face0
	t.Logf("wantD:\n%v", wantD)
	t.Logf("gotD:\n%v", glyph.d)
	if diff := cmp.Diff(wantD, glyph.d); diff != "" {
		t.Errorf("d mismatch (-want +got):\n%v", diff)
	}

	wantFaces := []*faceT{
		{
			dParts: []string{
				"M53 369", "q6 170 210 170", "q97 0 147 -36.5", "t50 -106.5", "v-308", "q0 -41 45 -41", "q9 0 18 2", "v-63", "q-35 -9 -57 -9", "q-40 0 -60.5 18.5", "t-25.5 58.5", "q-84 -77 -178 -77", "q-79 0 -125.5 42", "t-46.5 113", "q0 23 4.5 42", "t10 33.5", "t19.5 27", "t23.5 20.5", "t32 15.5", "t34.5 11", "t42 8.5",
				"t44 7", // cut0Idx=21
				"t50 7", "q49 6 68 19", "t19 39\n", "v22", "q0 38 -30.5 58", "t-86.5 20", "q-58 0 -88 -22.5", "t-35 -70.5",
				"h-84z", // idx=30
			},
			verts: []vec2{
				{53, 369}, {263, 539}, {410, 502.5}, {460, 396}, {460, 88}, {505, 47}, {523, 49}, {523, -14}, {466, -23}, {405.5, -4.5}, {380, 54}, {202, -23}, {76.5, 19}, {30, 132}, {34.5, 174}, {44.5, 207.5}, {64, 234.5}, {87.5, 255}, {119.5, 270.5}, {154, 281.5}, {196, 290},
				{240, 297}, // cut0Idx=21
				{290, 304}, {358, 323}, {377, 362}, {377, 384}, {346.5, 442}, {260, 462}, {172, 439.5}, {137, 369},
				{53, 369}, // idx=30
			},
		},
		{
			dParts: []string{
				"M220 50", "q69 0 113 36.5", "t44 78.5", "v94", "q-25 -12 -75.5 -20",
				"t-87.5 -14", // cutIdx=5
				"t-67 -28.5", "t-30 -62.5", "t27 -62",
				"t76 -22z", // idx=9
			},
			verts: []vec2{
				{220, 50}, {333, 86.5}, {377, 165}, {377, 259},
				{301.5, 239}, // cutIdx=5
				{214, 225}, {147, 196.5}, {117, 134}, {144, 72},
				{220, 50}, // idx=9
			},
			center:  vec2{245.05, 147.7},
			cut0Idx: 21,
			cutIdx:  5,
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

/*
func TestRegenerateFace(t *testing.T) {
	tests := []struct {
		name  string
		glyph *glyphT
		want  string
	}{
		{
			name: "a",
			glyph: &glyphT{
				d: "",
				faces: []*faceT{
					{
						dParts: []string{
							"M53 369q6 170 210 170q97 0 147 -36.5t50 -106.5",
							"",
							"",
							"v-308q0 -41 45 -41q9 0 18 2v-63q-35 -9 -57 -9q-40 0 -60.5 18.5t-25.5 58.5q-84 -77 -178 -77q-79 0 -125.5 42t-46.5 113q0 23 4.5 42t10 33.5t19.5 27t23.5 20.5t32 15.5t34.5 11t42 8.5t44 7t50 7q49 6 68 19t19 39\nv22q0 38 -30.5 58t-86.5 20q-58 0 -88 -22.5t-35 -70.5",
							"h-84z",
						},
						verts: []vec2{{53, 369}, {460, 88}, {523, -14}, {377, 384}, {53, 369}},
					},
					{
						dParts: []string{
							"M220 50q69 0 113 36.5t44 78.5",
							"v94q-25 -12 -75.5 -20t-87.5 -14t-67 -28.5t-30 -62.5t27 -62t76 -22z",
						},
						verts:   []vec2{{220, 50}, {377, 259}},
						center:  vec2{298.5, 154.5},
						cut0Idx: 1,
						cutIdx:  1,
					},
				},
			},
			want: "",
		},
	}

}
*/
