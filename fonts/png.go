package fonts

import (
	"log"

	"github.com/fogleman/gg"
)

func (r *Render) SavePNG(filename string, width, height int) error {
	scale := float64(width) / (r.MBB.Max[0] - r.MBB.Min[0])
	if yScale := float64(height) / (r.MBB.Max[1] - r.MBB.Min[1]); yScale < scale {
		scale = yScale
		width = int(0.5 + scale*(r.MBB.Max[0]-r.MBB.Min[0]))
	} else {
		height = int(0.5 + scale*(r.MBB.Max[1]-r.MBB.Min[1]))
	}
	log.Printf("MBB: %v, scale=%.2f, size=(%v,%v)", r.MBB, scale, width, height)

	dc := gg.NewContext(width, height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	for _, poly := range r.Polygons {
		if poly.Dark {
			dc.SetRGB(0, 0, 0)
		} else {
			dc.SetRGB(1, 1, 1)
		}
		for i, pt := range poly.Pts {
			if i == 0 {
				dc.MoveTo(scale*(pt[0]-r.MBB.Min[0]), float64(height)-scale*(pt[1]-r.MBB.Min[1]))
			} else {
				dc.LineTo(scale*(pt[0]-r.MBB.Min[0]), float64(height)-scale*(pt[1]-r.MBB.Min[1]))
			}
		}
		dc.Fill()
	}
	return dc.SavePNG(filename)
}
