package main

import (
	"reflect"
	"testing"

	"github.com/gmlewis/go3d/float64/vec2"
)

func TestRecorder_xIntercepts(t *testing.T) {
	tests := []struct {
		name    string
		topY    float64
		botY    float64
		segType segmentType
		pts     []vec2.T
		want    []vec2.T
	}{
		{
			name:    "latoregular 'i' dot left 1484-1494",
			segType: quadratic,
			topY:    1494,
			botY:    1484,
			pts: []vec2.T{
				{79.5, 1484}, {103, 1494}, {129, 1494},
			},
			want: []vec2.T{
				{129, 1494}, {79.5, 1484},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &segment{
				segType: tt.segType,
				pts:     tt.pts,
			}
			got := s.xIntercepts(tt.topY, tt.botY)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("xIntercepts = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestRecorder_interpFunc(t *testing.T) {
	tests := []struct {
		name    string
		segType segmentType
		pts     []vec2.T
		want    string
	}{
		{
			name:    "latoregular 'i' dot left 1484-1494",
			segType: quadratic,
			pts: []vec2.T{
				{79.5, 1484}, {103, 1494}, {129, 1494},
			},
			want: "interpQuadratic(vec2(79.50,1484.00),vec2(103.00,1494.00),vec2(129.00,1494.00),xyz.y)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &segment{
				segType: tt.segType,
				pts:     tt.pts,
			}
			got := s.interpFunc()
			if got != tt.want {
				t.Errorf("interpFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
