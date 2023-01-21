package freeserif

import (
	"testing"

	"github.com/gmlewis/go-fonts/fonts"
)

func TestTextMBB(t *testing.T) {
	// See: https://github.com/gmlewis/go-gerber/issues/8
	if _, err := fonts.TextMBB(0, 0, 1, 1, "", "freeserif"); err == nil {
		t.Error("TextMBB = nil, want err")
	}
}
