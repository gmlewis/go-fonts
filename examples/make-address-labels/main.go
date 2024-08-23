// -*- compile-command: "go run main.go postcard-labels.txt"; -*-

// make-address-labels reads a text file and generates a PDF that can be used
// to print on a standard 8.5x11 inch page of 30 labels arranged in a 3x10 grid.
// Each label is 66mm wide and 25mm high.
//
// The format of the text file is:
//
// Recipient1 Name
// Address Line 1
// Address Line 2
// Address Line 3
//
// Recipient2 Name
// etc...
//
// Usage:
//
//	make-address-labels addresses.txt
//
// This will write a file called "addresses.pdf".
package main

import (
	_ "embed"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-pdf/fpdf"
)

const (
	pdfXMarginMM  = 15
	pdfYMarginMM  = 25
	pdfDPI        = 300
	mmPerInch     = 25.4
	dpmm          = pdfDPI / mmPerInch // pixels/mm for 300DPI images
	widthInches   = 8.5
	heightInches  = 11
	widthMM       = widthInches * mmPerInch
	heightMM      = heightInches * mmPerInch
	labelWidthMM  = 66
	labelHeightMM = 25
	numLabelsX    = 3
	numLabelsY    = 10
	font1Family   = "Balsamiq Sans"
)

//go:embed BalsamiqSans-Regular.json
var BalsamiqSansRegularJSON []byte

//go:embed BalsamiqSans-Regular.z
var BalsamiqSansRegularZ []byte

func main() {
	flag.Parse()

	for _, filename := range flag.Args() {
		process(filename)
	}

	log.Printf("Done.")
}

func process(filename string) {
	buf, err := os.ReadFile(filename)
	must(err)

	addresses := strings.Split(string(buf), "\n\n")
	log.Printf("Got %v addresses from %v", len(addresses), filename)

	p := newPage()
	_, lineHeight1 := p.GetFontSize()

	for i, label := range addresses {
		nx, ny := i%numLabelsX, i/numLabelsX
		x := pdfXMarginMM + float64(nx)*(widthMM-pdfXMarginMM)/numLabelsX
		y := 2*lineHeight1 + float64(ny)*(heightMM-pdfYMarginMM)/numLabelsY
		lines := strings.Split(label, "\n")
		for j, line := range lines {
			p.SetXY(x, y+(lineHeight1+1.0)*float64(j))
			// log.Printf("(%v,%v)(%v,%v): %q", nx, ny, x, y, line)
			p.CellFormat(1, lineHeight1+10.0, line, "", 2, "AL", false, 0, "")
		}
	}

	pdfFilename := filepath.Join(filepath.Dir(filename), filepath.Base(filename)+".pdf")
	must(p.OutputFileAndClose(pdfFilename))
	log.Printf("Wrote %v", pdfFilename)
}

func newPage() fpdf.Pdf {
	p := fpdf.NewCustom(&fpdf.InitType{
		UnitStr: "mm",
		Size:    fpdf.SizeType{Wd: widthMM, Ht: heightMM},
	})
	p.AddPage()
	p.AddFontFromBytes(font1Family, "", BalsamiqSansRegularJSON, BalsamiqSansRegularZ)
	p.SetAutoPageBreak(false, 0)
	p.SetTextColor(0, 0, 0)
	p.SetFont(font1Family, "", 10)
	return p
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
