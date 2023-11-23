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
	p.MoveTo(g, 0, 0, "M", 53, 369)
	p.LineTo(g, 460, 396, "v", 460, 88)
	p.LineTo(g, 523, 49, "v", 523, -14)
	p.LineTo(g, 377, 362, "v", 377, 384)
	p.LineTo(g, 137, 369, "h", 53, 369)
	p.MoveTo(g, 53, 369, "M", 220, 50)
	p.LineTo(g, 377, 165, "v", 377, 259)
	p.ProcessGlyph('a', g)

	glyph, ok := p.glyphs["a"]
	if !ok {
		t.Fatalf("p.glyphs missing 'a': %+v", p.glyphs)
	}
	if p.current != glyph {
		t.Fatalf("p.current != glyph")
	}

	// 	wantD := "M53 369q6 170 210 170q97 0 147 -36.5t50 -106.5" + // start of face0
	// 		`v-308` + // MOVED FROM BELOW!
	// 		"L377 165" + // jump over to face[1] 259-94=165
	// 		"v94q-25 -12 -75.5 -20t-87.5 -14t-67 -28.5t-30 -62.5t27 -62t76 -22" + // rest of face[1] - strip "z"!
	// 		// "M220 50"+  // completely REMOVE 'M' of face[1]!
	// 		"q69 0 113 36.5t44 78.5" + // start of face[1] (after leading "M")
	// 		"L460 88" + // jump back to face[0]
	// 		// `v-308` + // MOVED ABOVE!!!
	// 		`q0 -41 45 -41q9 0 18 2v-63q-35 -9 -57 -9q-40 0 -60.5 18.5t-25.5 58.5q-84 -77 -178 -77q-79 0 -125.5 42t-46.5 113q0 23 4.5 42t10 33.5t19.5 27t23.5 20.5t32 15.5t34.5 11t42 8.5t44 7t50 7q49 6 68 19t19 39
	// v22q0 38 -30.5 58t-86.5 20q-58 0 -88 -22.5t-35 -70.5h-84z` // rest of face[0]

	wantD := "M53 369q6 170 210 170q97 0 147 -36.5t50 -106.5" + // start of face[0] dParts[0]
		"v-308q0 -41 45 -41q9 0 18 2" + // face[0] dParts[1]
		"v-63q-35 -9 -57 -9q-40 0 -60.5 18.5t-25.5 58.5q-84 -77 -178 -77q-79 0 -125.5 42t-46.5 113q0 23 4.5 42t10 33.5t19.5 27t23.5 20.5t32 15.5t34.5 11t42 8.5t44 7t50 7q49 6 68 19t19 39\n" + // face[0] dParts[2]
		"v22q0 38 -30.5 58t-86.5 20q-58 0 -88 -22.5t-35 -70.5" + // face[0] dParts[3]
		"h-84" + // face[0] dParts[4]
		"L377 165" + // jump over to face[1] before[1]
		"v94q-25 -12 -75.5 -20t-87.5 -14t-67 -28.5t-30 -62.5t27 -62t76 -22" + // face[1] dParts[1] - strip "z"!
		`v-308` + // MOVED FROM BELOW!
		"L377 165" + // jump over to face[1] 259-94=165
		"v94q-25 -12 -75.5 -20t-87.5 -14t-67 -28.5t-30 -62.5t27 -62t76 -22" + // rest of face[1] - strip "z"!
		// "M220 50"+  // completely REMOVE 'M' of face[1]!
		"q69 0 113 36.5t44 78.5" + // start of face[1] (after leading "M")
		"L460 88" + // jump back to face[0]
		// `v-308` + // MOVED ABOVE!!!
		`q0 -41 45 -41q9 0 18 2v-63q-35 -9 -57 -9q-40 0 -60.5 18.5t-25.5 58.5q-84 -77 -178 -77q-79 0 -125.5 42t-46.5 113q0 23 4.5 42t10 33.5t19.5 27t23.5 20.5t32 15.5t34.5 11t42 8.5t44 7t50 7q49 6 68 19t19 39
		v22q0 38 -30.5 58t-86.5 20q-58 0 -88 -22.5t-35 -70.5h-84z` // rest of face[0]

	t.Logf("wantD:\n%v", wantD)
	t.Logf("gotD:\n%v", glyph.d)
	if diff := cmp.Diff(wantD, glyph.d); diff != "" {
		t.Errorf("d mismatch (-want +got):\n%v", diff)
	}

	wantFaces := []*faceT{
		{
			dParts: []string{
				"M53 369q6 170 210 170q97 0 147 -36.5t50 -106.5",
				"v-308q0 -41 45 -41q9 0 18 2",
				"v-63q-35 -9 -57 -9q-40 0 -60.5 18.5t-25.5 58.5q-84 -77 -178 -77q-79 0 -125.5 42t-46.5 113q0 23 4.5 42t10 33.5t19.5 27t23.5 20.5t32 15.5t34.5 11t42 8.5t44 7t50 7q49 6 68 19t19 39\n",
				"v22q0 38 -30.5 58t-86.5 20q-58 0 -88 -22.5t-35 -70.5",
				"h-84z",
			},
			before: []vec2{{0, 0}, {460, 396}, {523, 49}, {377, 362}, {137, 369}},
			after:  []vec2{{53, 369}, {460, 88}, {523, -14}, {377, 384}, {53, 369}},
		},
		{
			dParts: []string{
				"M220 50q69 0 113 36.5t44 78.5",
				"v94q-25 -12 -75.5 -20t-87.5 -14t-67 -28.5t-30 -62.5t27 -62t76 -22z",
			},
			before:  []vec2{{53, 369}, {377, 165}},
			after:   []vec2{{220, 50}, {377, 259}},
			center:  vec2{215, 267},
			cut0Idx: 4,
			cutIdx:  0,
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
