package latoregular

import (
	"fmt"
	"testing"

	"github.com/gmlewis/go-fonts/fonts"
)

func TestText(t *testing.T) {
	tests := []struct {
		xPos, yPos     float64
		xScale, yScale float64
		message        string
		want           *fonts.Render
	}{
		{
			xScale:  1000.0,
			yScale:  1000.0,
			message: "012",
			want:    &fonts.Render{},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("test %v", tt.message), func(t *testing.T) {
			_, err := fonts.Text(tt.xPos, tt.yPos, tt.xScale, tt.yScale, tt.message, "latoregular")
			if err != nil {
				t.Fatal(err)
			}
			// TODO: Come up with a good way to test this with all the float64s.
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Text: got\n%v\nwant\n%v", goon.Sdump(got), goon.Sdump(tt.want))
			// }
		})
	}
}
