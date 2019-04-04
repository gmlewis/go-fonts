package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestWiring(t *testing.T) {
	tests := []struct {
		total          int
		onAxis         bool
		startTR        int
		wantStartLabel string
		wantEndLabel   string
		want           map[string]int
		wantConn       map[string]string
	}{
		{
			total:          20,
			onAxis:         true,
			startTR:        0,
			wantStartLabel: "TR",
			wantEndLabel:   "BL",
			want: map[string]int{
				"TR": 0, "TL": 10, "BR": 10, "BL": 0,
				"2R": 1, "2L": 11, "3R": 9, "3L": 19,
				"4R": 19, "4L": 9, "5R": 11, "5L": 1,
				"6R": 2, "6L": 12, "7R": 8, "7L": 18,
				"8R": 18, "8L": 8, "9R": 12, "9L": 2,
				"10R": 3, "10L": 13, "11R": 7, "11L": 17,
				"12R": 17, "12L": 7, "13R": 13, "13L": 3,
				"14R": 4, "14L": 14, "15R": 6, "15L": 16,
				"16R": 16, "16L": 6, "17R": 14, "17L": 4,
				"18R": 5, "18L": 15, "19R": 5, "19L": 15,
			},
			wantConn: map[string]string{
				"TR": "BL", "BL": "TR", // 0
				"2R": "5L", "5L": "2R", // 1
				"6R": "9L", "9L": "6R", // 2
				"10R": "13L", "13L": "10R", // 3
				"14R": "17L", "17L": "14R", // 4
				"18R": "19R", "19R": "18R", // 5
				"15R": "16L", "16L": "15R", // 6
				"11R": "12L", "12L": "11R", // 7
				"7R": "8L", "8L": "7R", // 8
				"3R": "4L", "4L": "3R", // 9
				"TL": "BR", "BR": "TL", // 10
				"2L": "5R", "5R": "2L", // 11
				"6L": "9R", "9R": "6L", // 12
				"10L": "13R", "13R": "10L", // 13
				"14L": "17R", "17R": "14L", // 14
				"18L": "19L", "19L": "18L", // 15
				"15L": "16R", "16R": "15L", // 16
				"11L": "12R", "12R": "11L", // 17
				"7L": "8R", "8R": "7L", // 18
				"3L": "4R", "4R": "3L", // 19
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v-layer", tt.total), func(t *testing.T) {
			opposite, mirror := genMaps(tt.total, tt.onAxis)
			startLabel, endLabel, got, connection := wiring(tt.startTR, tt.total, opposite, mirror)
			if startLabel != tt.wantStartLabel {
				t.Errorf("startLabel = %v, want %v", startLabel, tt.wantStartLabel)
			}
			if endLabel != tt.wantEndLabel {
				t.Errorf("endLabel = %v, want %v", endLabel, tt.wantEndLabel)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mirror(%v) = %#v, want %#v", tt.total, got, tt.want)
			}
			if !reflect.DeepEqual(connection, tt.wantConn) {
				t.Errorf("mirror(%v) connection = %#v, want %#v", tt.total, connection, tt.wantConn)
			}
		})
	}
}

func TestGenMaps(t *testing.T) {
	tests := []struct {
		total    int
		onAxis   bool
		opposite map[int]int
		mirror   map[int]int
	}{
		{
			total:  20,
			onAxis: true,
			opposite: map[int]int{
				0: 10, 10: 0,
				1: 11, 11: 1,
				2: 12, 12: 2,
				3: 13, 13: 3,
				4: 14, 14: 4,
				5: 15, 15: 5,
				6: 16, 16: 6,
				7: 17, 17: 7,
				8: 18, 18: 8,
				9: 19, 19: 9,
			},
			mirror: map[int]int{
				0: 10, 10: 0,
				1: 9, 9: 1,
				2: 8, 8: 2,
				3: 7, 7: 3,
				4: 6, 6: 4,
				5:  5,
				19: 11, 11: 19,
				18: 12, 12: 18,
				17: 13, 13: 17,
				16: 14, 14: 16,
				15: 15,
			},
		},
		{
			total:  20,
			onAxis: false,
			opposite: map[int]int{
				0: 10, 10: 0,
				1: 11, 11: 1,
				2: 12, 12: 2,
				3: 13, 13: 3,
				4: 14, 14: 4,
				5: 15, 15: 5,
				6: 16, 16: 6,
				7: 17, 17: 7,
				8: 18, 18: 8,
				9: 19, 19: 9,
			},
			mirror: map[int]int{
				0: 9, 9: 0,
				1: 8, 8: 1,
				2: 7, 7: 2,
				3: 6, 6: 3,
				4: 5, 5: 4,
				19: 10, 10: 19,
				18: 11, 11: 18,
				17: 12, 12: 17,
				16: 13, 13: 16,
				15: 14, 14: 15,
			},
		},
		{
			total:  6,
			onAxis: true,
			opposite: map[int]int{
				0: 3, 3: 0,
				1: 4, 4: 1,
				2: 5, 5: 2,
			},
			mirror: map[int]int{
				0: 3, 3: 0,
				1: 2, 2: 1,
				4: 5, 5: 4,
			},
		},
		{
			total:  6,
			onAxis: false,
			opposite: map[int]int{
				0: 3, 3: 0,
				1: 4, 4: 1,
				2: 5, 5: 2,
			},
			mirror: map[int]int{
				0: 2, 2: 0,
				1: 1, 4: 4,
				3: 5, 5: 3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v-layer(%v)", tt.total, tt.onAxis), func(t *testing.T) {
			o, m := genMaps(tt.total, tt.onAxis)
			if !reflect.DeepEqual(o, tt.opposite) {
				t.Errorf("genMaps(%v): opposite = %#v, want %#v", tt.total, o, tt.opposite)
			}
			if !reflect.DeepEqual(m, tt.mirror) {
				t.Errorf("genMaps(%v): mirror = %#v, want %#v", tt.total, m, tt.mirror)
			}
		})
	}
}
