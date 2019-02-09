package latoregular

import (
	"fmt"
	"math"
	"testing"

	"github.com/gmlewis/go-fonts/fonts"
)

const (
	fontName = "latoregular"
	eps      = 1e-6
)

func TestTextMBB(t *testing.T) {
	const (
		message = "012"
	)

	tests := []struct {
		name       string
		xPos, yPos float64
		wantXmin   float64
		wantYmin   float64
		wantXmax   float64
		wantYmax   float64
	}{
		{
			name:     "Origin",
			wantXmin: 0.02978515625,
			wantYmin: -0.00732421875,
			wantXmax: 1.68896484375,
			wantYmax: 0.724609375,
		},
		{
			name:     "Offset",
			xPos:     10,
			yPos:     20,
			wantXmin: 10 + 0.02978515625,
			wantYmin: 20 + -0.00732421875,
			wantXmax: 10 + 1.68896484375,
			wantYmax: 20 + 0.724609375,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("test %v", message), func(t *testing.T) {
			mbb, err := fonts.TextMBB(tt.xPos, tt.yPos, 1, 1, message, fontName)
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

func TestText(t *testing.T) {
	const (
		message    = "012"
		wantWidth  = 1.6591796875
		wantHeight = 0.73193359375
		angle      = math.Pi / 3.0
	)

	tests := []struct {
		name       string
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
		// Tests with rotation...
		{
			name:     "XLeft,YBottom w/ rotation",
			rotation: angle,
			wantXmin: -0.014448926977988859,
			wantYmin: 0.03388734704564836,
			wantXmax: 0.8891355407855701,
			wantYmax: 1.5989415541422152,
		},
		{
			name:     "XLeft,YBottom w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			wantXmin: 9.98555107302201,
			wantYmin: 20.03388734704565,
			wantXmax: 10.88913554078557,
			wantYmax: 21.598941554142215,
		},
		{
			name:     "XCenter,YBottom w/ rotation",
			rotation: angle,
			opts:     fonts.BottomCenter,
			wantXmin: -0.8440387707279888,
			wantYmin: 0.03388734704564836,
			wantXmax: 0.05954569703557011,
			wantYmax: 1.5989415541422152,
		},
		{
			name:     "XCenter,YBottom w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			opts:     fonts.BottomCenter,
			wantXmin: 9.15596122927201,
			wantYmin: 20.03388734704565,
			wantXmax: 10.05954569703557,
			wantYmax: 21.598941554142215,
		},
		{
			name:     "XRight,YBottom w/ rotation",
			rotation: angle,
			opts:     fonts.BottomRight,
			wantXmin: -1.6736286144779888,
			wantYmin: 0.03388734704564836,
			wantXmax: -0.7700441467144299,
			wantYmax: 1.5989415541422152,
		},
		{
			name:     "XRight,YBottom w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			opts:     fonts.BottomRight,
			wantXmin: 8.32637138552201,
			wantYmin: 20.03388734704565,
			wantXmax: 9.22995585328557,
			wantYmax: 21.598941554142215,
		},
		{
			name:     "XLeft,YCenter w/ rotation",
			rotation: angle,
			opts:     fonts.CenterLeft,
			wantXmin: -0.014448926977988859,
			wantYmin: -0.33207944982935167,
			wantXmax: 0.8891355407855701,
			wantYmax: 1.2329747572672152,
		},
		{
			name:     "XLeft,YCenter w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			opts:     fonts.CenterLeft,
			wantXmin: 9.98555107302201,
			wantYmin: 19.66792055017065,
			wantXmax: 10.88913554078557,
			wantYmax: 21.232974757267215,
		},
		{
			name:     "XCenter,YCenter w/ rotation",
			rotation: angle,
			opts:     fonts.Center,
			wantXmin: -0.8440387707279888,
			wantYmin: -0.33207944982935167,
			wantXmax: 0.05954569703557011,
			wantYmax: 1.2329747572672152,
		},
		{
			name:     "XCenter,YCenter w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			opts:     fonts.Center,
			wantXmin: 9.15596122927201,
			wantYmin: 19.66792055017065,
			wantXmax: 10.05954569703557,
			wantYmax: 21.232974757267215,
		},
		{
			name:     "XRight,YCenter w/ rotation",
			rotation: angle,
			opts:     fonts.CenterRight,
			wantXmin: -1.6736286144779888,
			wantYmin: -0.33207944982935167,
			wantXmax: -0.7700441467144299,
			wantYmax: 1.2329747572672152,
		},
		{
			name:     "XRight,YCenter w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			opts:     fonts.CenterRight,
			wantXmin: 8.32637138552201,
			wantYmin: 19.66792055017065,
			wantXmax: 9.22995585328557,
			wantYmax: 21.232974757267215,
		},
		{
			name:     "XLeft,YTop w/ rotation",
			rotation: angle,
			opts:     fonts.TopLeft,
			wantXmin: -0.014448926977988859,
			wantYmin: -0.6980462467043517,
			wantXmax: 0.8891355407855701,
			wantYmax: 0.8670079603922152,
		},
		{
			name:     "XLeft,YTop w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			opts:     fonts.TopLeft,
			wantXmin: 9.98555107302201,
			wantYmin: 19.30195375329565,
			wantXmax: 10.88913554078557,
			wantYmax: 20.867007960392215,
		},
		{
			name:     "XCenter,YTop w/ rotation",
			rotation: angle,
			opts:     fonts.TopCenter,
			wantXmin: -0.8440387707279888,
			wantYmin: -0.6980462467043517,
			wantXmax: 0.05954569703557011,
			wantYmax: 0.8670079603922152,
		},
		{
			name:     "XCenter,YTop w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			opts:     fonts.TopCenter,
			wantXmin: 9.15596122927201,
			wantYmin: 19.30195375329565,
			wantXmax: 10.05954569703557,
			wantYmax: 20.867007960392215,
		},
		{
			name:     "XRight,YTop w/ rotation",
			rotation: angle,
			opts:     fonts.TopRight,
			wantXmin: -1.6736286144779888,
			wantYmin: -0.6980462467043517,
			wantXmax: -0.7700441467144299,
			wantYmax: 0.8670079603922152,
		},
		{
			name:     "XRight,YTop w/ offset w/ rotation",
			rotation: angle,
			xPos:     10,
			yPos:     20,
			opts:     fonts.TopRight,
			wantXmin: 8.32637138552201,
			wantYmin: 19.30195375329565,
			wantXmax: 9.22995585328557,
			wantYmax: 20.867007960392215,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("test Text %v", message), func(t *testing.T) {
			opts := tt.opts // make a copy
			opts.Rotate = tt.rotation
			got, err := fonts.Text(tt.xPos, tt.yPos, 1, 1, message, fontName, &opts)
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
