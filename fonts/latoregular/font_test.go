package latoregular

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/gmlewis/go-fonts/fonts"
)

const (
	fontName = "latoregular"
	eps      = 1e-3
)

func TestTextMBB(t *testing.T) {
	tests := []struct {
		name       string
		message    string
		xPos, yPos float64
		wantXmin   float64
		wantYmin   float64
		wantXmax   float64
		wantYmax   float64
	}{
		{
			name:     "Origin",
			message:  "012",
			wantXmin: 0.02978515625,
			wantYmin: -0.00732421875,
			wantXmax: 1.68896484375,
			wantYmax: 0.724609375,
		},
		{
			name:     "Offset",
			message:  "012",
			xPos:     10,
			yPos:     20,
			wantXmin: 10 + 0.02978515625,
			wantYmin: 20 + -0.00732421875,
			wantXmax: 10 + 1.68896484375,
			wantYmax: 20 + 0.724609375,
		},
		{
			name:     "Origin",
			message:  "M",
			wantXmin: 0.0869140625,
			wantYmin: 0,
			wantXmax: 0.83251953125,
			wantYmax: 0.71630859375,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v %v", tt.name, tt.message), func(t *testing.T) {
			mbb, err := fonts.TextMBB(tt.xPos, tt.yPos, 1, 1, tt.message, fontName)
			if err != nil {
				t.Fatal(err)
			}
			if math.Abs(mbb.Min[0]-tt.wantXmin) > eps {
				t.Errorf("Xmin = %v, want %v", mbb.Min[0], tt.wantXmin)
			}
			if math.Abs(mbb.Min[1]-tt.wantYmin) > eps {
				t.Errorf("Ymin = %v, want %v", mbb.Min[1], tt.wantYmin)
			}
			if math.Abs(mbb.Max[0]-tt.wantXmax) > eps {
				t.Errorf("Xmax = %v, want %v", mbb.Max[0], tt.wantXmax)
			}
			if math.Abs(mbb.Max[1]-tt.wantYmax) > eps {
				t.Errorf("Ymax = %v, want %v", mbb.Max[1], tt.wantYmax)
			}
		})
	}
}

func sp(s string) *string {
	return &s
}

func TestText(t *testing.T) {
	const (
		message    = "012"
		wantWidth  = 1.6591796875
		wantHeight = 0.73193359375
		angle      = math.Pi / 3.0
	)

	tests := []struct {
		name       string
		message    *string
		xPos, yPos float64
		rotation   float64
		opts       fonts.TextOpts
		wantXmin   float64
		wantYmin   float64
		wantXmax   float64
		wantYmax   float64
	}{
		{
			name:     "XLeft,YBottom",
			message:  sp("M"),
			wantXmin: 0,
			wantYmin: 0,
			wantXmax: 0.74560546875,
			wantYmax: 0.71630859375,
		},
		{
			name:     "XLeft,YBottom",
			wantXmin: 0,
			wantYmin: 0,
			wantXmax: wantWidth,
			wantYmax: wantHeight,
		},
		{
			name:     "XLeft,YBottom w/ offset",
			xPos:     10,
			yPos:     20,
			wantXmin: 10,
			wantYmin: 20,
			wantXmax: (10 + wantWidth),
			wantYmax: (20 + wantHeight),
		},
		{
			name:     "XCenter,YBottom",
			opts:     fonts.BottomCenter,
			wantXmin: -0.5 * wantWidth,
			wantYmin: 0,
			wantXmax: 0.5 * wantWidth,
			wantYmax: wantHeight,
		},
		{
			name:     "XCenter,YBottom w/ offset",
			xPos:     10,
			yPos:     20,
			opts:     fonts.BottomCenter,
			wantXmin: (10 - 0.5*wantWidth),
			wantYmin: 20,
			wantXmax: (10 + 0.5*wantWidth),
			wantYmax: (20 + wantHeight),
		},
		{
			name:     "XRight,YBottom",
			opts:     fonts.BottomRight,
			wantXmin: -wantWidth,
			wantYmin: 0,
			wantXmax: 0,
			wantYmax: wantHeight,
		},
		{
			name:     "XRight,YBottom w/ offset",
			xPos:     10,
			yPos:     20,
			opts:     fonts.BottomRight,
			wantXmin: (10 - wantWidth),
			wantYmin: 20,
			wantXmax: 10,
			wantYmax: (20 + wantHeight),
		},
		{
			name:     "XLeft,YCenter",
			opts:     fonts.CenterLeft,
			wantXmin: 0,
			wantYmin: -0.5 * wantHeight,
			wantXmax: wantWidth,
			wantYmax: 0.5 * wantHeight,
		},
		{
			name:     "XLeft,YCenter w/ offset",
			xPos:     10,
			yPos:     20,
			opts:     fonts.CenterLeft,
			wantXmin: 10,
			wantYmin: (20 - 0.5*wantHeight),
			wantXmax: (10 + wantWidth),
			wantYmax: (20 + 0.5*wantHeight),
		},
		{
			name:     "XCenter,YCenter",
			opts:     fonts.Center,
			wantXmin: -0.5 * wantWidth,
			wantYmin: -0.5 * wantHeight,
			wantXmax: 0.5 * wantWidth,
			wantYmax: 0.5 * wantHeight,
		},
		{
			name:     "XCenter,YCenter w/ offset",
			xPos:     10,
			yPos:     20,
			opts:     fonts.Center,
			wantXmin: (10 - 0.5*wantWidth),
			wantYmin: (20 - 0.5*wantHeight),
			wantXmax: (10 + 0.5*wantWidth),
			wantYmax: (20 + 0.5*wantHeight),
		},
		{
			name:     "XRight,YCenter",
			opts:     fonts.CenterRight,
			wantXmin: -wantWidth,
			wantYmin: -0.5 * wantHeight,
			wantXmax: 0,
			wantYmax: 0.5 * wantHeight,
		},
		{
			name:     "XRight,YCenter w/ offset",
			xPos:     10,
			yPos:     20,
			opts:     fonts.CenterRight,
			wantXmin: (10 - wantWidth),
			wantYmin: (20 - 0.5*wantHeight),
			wantXmax: 10,
			wantYmax: (20 + 0.5*wantHeight),
		},
		{
			name:     "XLeft,YTop",
			opts:     fonts.TopLeft,
			wantXmin: 0,
			wantYmin: -wantHeight,
			wantXmax: wantWidth,
			wantYmax: 0,
		},
		{
			name:     "XLeft,YTop w/ offset",
			xPos:     10,
			yPos:     20,
			opts:     fonts.TopLeft,
			wantXmin: 10,
			wantYmin: (20 - wantHeight),
			wantXmax: (10 + wantWidth),
			wantYmax: 20,
		},
		{
			name:     "XCenter,YTop",
			opts:     fonts.TopCenter,
			wantXmin: -0.5 * wantWidth,
			wantYmin: -wantHeight,
			wantXmax: 0.5 * wantWidth,
			wantYmax: 0,
		},
		{
			name:     "XCenter,YTop w/ offset",
			xPos:     10,
			yPos:     20,
			opts:     fonts.TopCenter,
			wantXmin: (10 - 0.5*wantWidth),
			wantYmin: (20 - wantHeight),
			wantXmax: (10 + 0.5*wantWidth),
			wantYmax: 20,
		},
		{
			name:     "XRight,YTop",
			opts:     fonts.TopRight,
			wantXmin: -wantWidth,
			wantYmin: -wantHeight,
			wantXmax: 0,
			wantYmax: 0,
		},
		{
			name:     "XRight,YTop w/ offset",
			xPos:     10,
			yPos:     20,
			opts:     fonts.TopRight,
			wantXmin: (10 - wantWidth),
			wantYmin: (20 - wantHeight),
			wantXmax: 10,
			wantYmax: 20,
		},
		// Tests with rotation:
		{
			name:     "XLeft,YBottom w/ rotation",
			rotation: angle,
			wantXmin: -0.5454372744403848,
			wantYmin: 0.1364595374423504,
			wantXmax: 0.8146972656249999,
			wantYmax: 1.7488829588350612,
		},
		{
			name:     "XLeft,YBottom w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			wantXmin: 9.454562725559613,
			wantYmin: 20.13645953744235,
			wantXmax: 10.814697265625,
			wantYmax: 21.74888295883506,
		},
		{
			name:     "XCenter,YBottom w/ rotation",
			rotation: angle,
			opts:     fonts.BottomCenter,
			wantXmin: -1.3750271181903848,
			wantYmin: 0.13645953744235043,
			wantXmax: -0.014892578125000111,
			wantYmax: 1.7488829588350612,
		},
		{
			name:     "XCenter,YBottom w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			opts:     fonts.BottomCenter,
			wantXmin: 8.624972881809615,
			wantYmin: 20.13645953744235,
			wantXmax: 9.985107421875,
			wantYmax: 21.74888295883506,
		},
		{
			name:     "XRight,YBottom w/ rotation",
			rotation: angle,
			opts:     fonts.BottomRight,
			wantXmin: -2.204616961940385,
			wantYmin: 0.13645953744235034,
			wantXmax: -0.8444824218750001,
			wantYmax: 1.7488829588350612,
		},
		{
			name:     "XRight,YBottom w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			opts:     fonts.BottomRight,
			wantXmin: 7.795383038059613,
			wantYmin: 20.13645953744235,
			wantXmax: 9.155517578125,
			wantYmax: 21.74888295883506,
		},
		{
			name:     "XLeft,YCenter w/ rotation",
			rotation: angle,
			opts:     fonts.CenterLeft,
			wantXmin: -0.5454372744403848,
			wantYmin: -0.2295072594326496,
			wantXmax: 0.8146972656249999,
			wantYmax: 1.3829161619600612,
		},
		{
			name:     "XLeft,YCenter w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			opts:     fonts.CenterLeft,
			wantXmin: 9.454562725559613,
			wantYmin: 19.77049274056735,
			wantXmax: 10.814697265625,
			wantYmax: 21.38291616196006,
		},
		{
			name:     "XCenter,YCenter w/ rotation",
			rotation: angle,
			opts:     fonts.Center,
			wantXmin: -1.3750271181903848,
			wantYmin: -0.22950725943264957,
			wantXmax: -0.014892578125000111,
			wantYmax: 1.3829161619600612,
		},
		{
			name:     "XCenter,YCenter w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			opts:     fonts.Center,
			wantXmin: 8.624972881809615,
			wantYmin: 19.77049274056735,
			wantXmax: 9.985107421875,
			wantYmax: 21.38291616196006,
		},
		{
			name:     "XRight,YCenter w/ rotation",
			rotation: angle,
			opts:     fonts.CenterRight,
			wantXmin: -2.204616961940385,
			wantYmin: -0.22950725943264966,
			wantXmax: -0.8444824218750001,
			wantYmax: 1.3829161619600612,
		},
		{
			name:     "XRight,YCenter w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			opts:     fonts.CenterRight,
			wantXmin: 7.795383038059613,
			wantYmin: 19.77049274056735,
			wantXmax: 9.155517578125,
			wantYmax: 21.38291616196006,
		},
		{
			name:     "XLeft,YTop w/ rotation",
			rotation: angle,
			opts:     fonts.TopLeft,
			wantXmin: -0.5454372744403848,
			wantYmin: -0.5954740563076496,
			wantXmax: 0.8146972656249999,
			wantYmax: 1.0169493650850612,
		},
		{
			name:     "XLeft,YTop w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			opts:     fonts.TopLeft,
			wantXmin: 9.454562725559613,
			wantYmin: 19.40452594369235,
			wantXmax: 10.814697265625,
			wantYmax: 21.01694936508506,
		},
		{
			name:     "XCenter,YTop w/ rotation",
			rotation: angle,
			opts:     fonts.TopCenter,
			wantXmin: -1.3750271181903848,
			wantYmin: -0.5954740563076496,
			wantXmax: -0.014892578125000111,
			wantYmax: 1.0169493650850612,
		},
		{
			name:     "XCenter,YTop w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			opts:     fonts.TopCenter,
			wantXmin: 8.624972881809615,
			wantYmin: 19.40452594369235,
			wantXmax: 9.985107421875,
			wantYmax: 21.01694936508506,
		},
		{
			name:     "XRight,YTop w/ rotation",
			rotation: angle,
			opts:     fonts.TopRight,
			wantXmin: -2.204616961940385,
			wantYmin: -0.5954740563076497,
			wantXmax: -0.8444824218750001,
			wantYmax: 1.0169493650850612,
		},
		{
			name:     "XRight,YTop w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			opts:     fonts.TopRight,
			wantXmin: 7.795383038059613,
			wantYmin: 19.40452594369235,
			wantXmax: 9.155517578125,
			wantYmax: 21.01694936508506,
		},
	}

	for _, tt := range tests {
		m := message
		if tt.message != nil {
			m = *tt.message
		}
		t.Run(fmt.Sprintf("test %v %v", tt.name, m), func(t *testing.T) {
			opts := tt.opts // make a copy
			opts.Rotate = tt.rotation
			got, err := fonts.Text(tt.xPos, tt.yPos, 1, 1, m, fontName, &opts)
			if err != nil {
				t.Fatal(err)
			}
			mbb := got.MBB
			if math.Abs(mbb.Min[0]-tt.wantXmin) > eps {
				t.Errorf("Xmin = %v, want %v", mbb.Min[0], tt.wantXmin)
			}
			if math.Abs(mbb.Min[1]-tt.wantYmin) > eps {
				t.Errorf("Ymin = %v, want %v", mbb.Min[1], tt.wantYmin)
			}
			if math.Abs(mbb.Max[0]-tt.wantXmax) > eps {
				t.Errorf("Xmax = %v, want %v", mbb.Max[0], tt.wantXmax)
			}
			if math.Abs(mbb.Max[1]-tt.wantYmax) > eps {
				t.Errorf("Ymax = %v, want %v", mbb.Max[1], tt.wantYmax)
			}
		})
	}
}

func TestFillBox(t *testing.T) {
	const wideMsg = "this is a long, horizontal line"
	var tallMsg = strings.Join(strings.Split("this is a tall, vertical line", ""), "\n")

	tests := []struct {
		name    string
		mbb     fonts.MBB
		xScale  float64
		message string
		opts    *fonts.TextOpts
		want    fonts.MBB
	}{
		{
			name:    "Square, no opts - XLeft,YBottom",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: "M",
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 96.07072691552064}},
		},
		{
			name:    "Square, mirror, no opts - XLeft,YBottom",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: "M",
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 96.07072691552064}},
		},

		{
			name:    "Wide, no opts - XLeft,YBottom",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: wideMsg,
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 7.738168660828091}},
		},
		{
			name:    "Wide, mirror, no opts - XLeft,YBottom",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: wideMsg,
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 7.738168660828091}},
		},

		{
			name:    "Wide, BottomLeft",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: wideMsg,
			opts:    &fonts.BottomLeft,
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 7.738168660828091}},
		},
		{
			name:    "Wide, mirror, BottomLeft",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: wideMsg,
			opts:    &fonts.BottomLeft,
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 7.738168660828091}},
		},

		{
			name:    "Wide, BottomCenter",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: wideMsg,
			opts:    &fonts.BottomCenter,
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 7.738168660828091}},
		},
		{
			name:    "Wide, mirror, BottomCenter",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: wideMsg,
			opts:    &fonts.BottomCenter,
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 7.738168660828091}},
		},

		{
			name:    "Wide, BottomRight",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: wideMsg,
			opts:    &fonts.BottomRight,
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 7.738168660828091}},
		},
		{
			name:    "Wide, mirror, BottomRight",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: wideMsg,
			opts:    &fonts.BottomRight,
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 7.738168660828091}},
		},

		{
			name:    "Wide, CenterLeft",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: wideMsg,
			opts:    &fonts.CenterLeft,
			want:    fonts.MBB{Min: fonts.Pt{0, 46.1309156696}, Max: fonts.Pt{100, 53.8690843304}},
		},
		{
			name:    "Wide, mirror, CenterLeft",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: wideMsg,
			opts:    &fonts.CenterLeft,
			want:    fonts.MBB{Min: fonts.Pt{0, 46.1309156696}, Max: fonts.Pt{100, 53.8690843304}},
		},

		{
			name:    "Wide, Center",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: wideMsg,
			opts:    &fonts.Center,
			want:    fonts.MBB{Min: fonts.Pt{0, 46.1309156696}, Max: fonts.Pt{100, 53.8690843304}},
		},
		{
			name:    "Wide, mirror, Center",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: wideMsg,
			opts:    &fonts.Center,
			want:    fonts.MBB{Min: fonts.Pt{0, 46.1309156696}, Max: fonts.Pt{100, 53.8690843304}},
		},

		{
			name:    "Wide, CenterRight",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: wideMsg,
			opts:    &fonts.CenterRight,
			want:    fonts.MBB{Min: fonts.Pt{0, 46.1309156696}, Max: fonts.Pt{100, 53.8690843304}},
		},
		{
			name:    "Wide, mirror, CenterRight",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: wideMsg,
			opts:    &fonts.CenterRight,
			want:    fonts.MBB{Min: fonts.Pt{0, 46.1309156696}, Max: fonts.Pt{100, 53.8690843304}},
		},

		{
			name:    "Wide, TopLeft",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: wideMsg,
			opts:    &fonts.TopLeft,
			want:    fonts.MBB{Min: fonts.Pt{0, 92.2618313391719}, Max: fonts.Pt{100, 100}},
		},
		{
			name:    "Wide, mirror, TopLeft",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: wideMsg,
			opts:    &fonts.TopLeft,
			want:    fonts.MBB{Min: fonts.Pt{0, 92.2618313391719}, Max: fonts.Pt{100, 100}},
		},

		{
			name:    "Wide, TopCenter",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: wideMsg,
			opts:    &fonts.TopCenter,
			want:    fonts.MBB{Min: fonts.Pt{0, 92.2618313391719}, Max: fonts.Pt{100, 100}},
		},
		{
			name:    "Wide, mirror, TopCenter",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: wideMsg,
			opts:    &fonts.TopCenter,
			want:    fonts.MBB{Min: fonts.Pt{0, 92.2618313391719}, Max: fonts.Pt{100, 100}},
		},

		{
			name:    "Wide, TopRight",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: wideMsg,
			opts:    &fonts.TopRight,
			want:    fonts.MBB{Min: fonts.Pt{0, 92.2618313391719}, Max: fonts.Pt{100, 100}},
		},
		{
			name:    "Wide, mirror, TopRight",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: wideMsg,
			opts:    &fonts.TopRight,
			want:    fonts.MBB{Min: fonts.Pt{0, 92.2618313391719}, Max: fonts.Pt{100, 100}},
		},

		// Tall
		{
			name:    "Tall, no opts - XLeft,YBottom",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: tallMsg,
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{1.721378464891371, 100}},
		},
		{
			name:    "Tall, mirror, no opts - XLeft,YBottom",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: tallMsg,
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{1.721378464891371, 100}},
		},

		{
			name:    "Tall, BottomLeft",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: tallMsg,
			opts:    &fonts.BottomLeft,
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{1.721378464891371, 100}},
		},
		{
			name:    "Tall, mirror, BottomLeft",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: tallMsg,
			opts:    &fonts.BottomLeft,
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{1.721378464891371, 100}},
		},

		{
			name:    "Tall, BottomCenter",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: tallMsg,
			opts:    &fonts.BottomCenter,
			want:    fonts.MBB{Min: fonts.Pt{49.1393107676, 0}, Max: fonts.Pt{50.8606892324, 100}},
		},
		{
			name:    "Tall, mirror, BottomCenter",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: tallMsg,
			opts:    &fonts.BottomCenter,
			want:    fonts.MBB{Min: fonts.Pt{49.1393107676, 0}, Max: fonts.Pt{50.8606892324, 100}},
		},

		{
			name:    "Tall, BottomRight",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: tallMsg,
			opts:    &fonts.BottomRight,
			want:    fonts.MBB{Min: fonts.Pt{98.2786215351, 0}, Max: fonts.Pt{100, 100}},
		},
		{
			name:    "Tall, mirror, BottomRight",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: tallMsg,
			opts:    &fonts.BottomRight,
			want:    fonts.MBB{Min: fonts.Pt{98.2786215351, 0}, Max: fonts.Pt{100, 100}},
		},

		{
			name:    "Tall, CenterLeft",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: tallMsg,
			opts:    &fonts.CenterLeft,
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{1.721378464891371, 100}},
		},
		{
			name:    "Tall, mirror, CenterLeft",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: tallMsg,
			opts:    &fonts.CenterLeft,
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{1.721378464891371, 100}},
		},

		{
			name:    "Tall, Center",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: tallMsg,
			opts:    &fonts.Center,
			want:    fonts.MBB{Min: fonts.Pt{49.1393107676, 0}, Max: fonts.Pt{50.8606892324, 100}},
		},
		{
			name:    "Tall, mirror, Center",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: tallMsg,
			opts:    &fonts.Center,
			want:    fonts.MBB{Min: fonts.Pt{49.1393107676, 0}, Max: fonts.Pt{50.8606892324, 100}},
		},

		{
			name:    "Tall, CenterRight",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: tallMsg,
			opts:    &fonts.CenterRight,
			want:    fonts.MBB{Min: fonts.Pt{98.2786215351, 0}, Max: fonts.Pt{100, 100}},
		},
		{
			name:    "Tall, mirror, CenterRight",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: tallMsg,
			opts:    &fonts.CenterRight,
			want:    fonts.MBB{Min: fonts.Pt{98.2786215351, 0}, Max: fonts.Pt{100, 100}},
		},

		{
			name:    "Tall, TopLeft",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: tallMsg,
			opts:    &fonts.TopLeft,
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{1.721378464891371, 100}},
		},
		{
			name:    "Tall, mirror, TopLeft",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: tallMsg,
			opts:    &fonts.TopLeft,
			want:    fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{1.721378464891371, 100}},
		},

		{
			name:    "Tall, TopCenter",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: tallMsg,
			opts:    &fonts.TopCenter,
			want:    fonts.MBB{Min: fonts.Pt{49.1393107676, 0}, Max: fonts.Pt{50.8606892324, 100}},
		},
		{
			name:    "Tall, mirror, TopCenter",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: tallMsg,
			opts:    &fonts.TopCenter,
			want:    fonts.MBB{Min: fonts.Pt{49.1393107676, 0}, Max: fonts.Pt{50.8606892324, 100}},
		},

		{
			name:    "Tall, TopRight",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  1.0,
			message: tallMsg,
			opts:    &fonts.TopRight,
			want:    fonts.MBB{Min: fonts.Pt{98.2786215351, 0}, Max: fonts.Pt{100, 100}},
		},
		{
			name:    "Tall, mirror, TopRight",
			mbb:     fonts.MBB{Min: fonts.Pt{0, 0}, Max: fonts.Pt{100, 100}},
			xScale:  -1.0,
			message: tallMsg,
			opts:    &fonts.TopRight,
			want:    fonts.MBB{Min: fonts.Pt{98.2786215351, 0}, Max: fonts.Pt{100, 100}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y, pts, err := fonts.FillBox(tt.mbb, tt.xScale, 1.0, tt.message, fontName, tt.opts)
			if err != nil {
				t.Fatalf("FillBox: %v", err)
			}
			// log.Printf("x=%v, y=%v, pts=%v", x, y, pts)
			got, err := fonts.Text(x, y, tt.xScale*pts, pts, tt.message, fontName, tt.opts)
			if err != nil {
				t.Fatalf("Text: %v", err)
			}
			// log.Printf("w,h=(%v,%v), got=%#v", got.MBB.Max[0]-got.MBB.Min[0], got.MBB.Max[1]-got.MBB.Min[1], got.MBB)

			mbb := got.MBB
			if math.Abs(mbb.Min[0]-tt.want.Min[0]) > eps {
				t.Errorf("Xmin = %v, want %v", mbb.Min[0], tt.want.Min[0])
			}
			if math.Abs(mbb.Min[1]-tt.want.Min[1]) > eps {
				t.Errorf("Ymin = %v, want %v", mbb.Min[1], tt.want.Min[1])
			}
			if math.Abs(mbb.Max[0]-tt.want.Max[0]) > eps {
				t.Errorf("Xmax = %v, want %v", mbb.Max[0], tt.want.Max[0])
			}
			if math.Abs(mbb.Max[1]-tt.want.Max[1]) > eps {
				t.Errorf("Ymax = %v, want %v", mbb.Max[1], tt.want.Max[1])
			}
		})
	}
}
