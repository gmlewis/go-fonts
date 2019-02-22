package fonts

import (
	"github.com/yofu/dxf"
)

func (r *Render) SaveDXF(filename string, width, height int) error {
	d := dxf.NewDrawing()
	d.Header().LtScale = 100.0
	d.AddLayer("Text", dxf.DefaultColor, dxf.DefaultLineType, true)
	d.ChangeLayer("Text")

	for _, poly := range r.Polygons {
		lastX, lastY := poly.Pts[0][0]-r.MBB.Min[0], poly.Pts[0][1]-r.MBB.Min[1]
		for _, pt := range poly.Pts[1:] {
			x, y := pt[0]-r.MBB.Min[0], pt[1]-r.MBB.Min[1]
			d.Line(lastX, lastY, 0.0, x, y, 0.0)
			lastX, lastY = x, y
		}
		d.Line(lastX, lastY, 0.0, poly.Pts[0][0]-r.MBB.Min[0], poly.Pts[0][1]-r.MBB.Min[1], 0.0)
	}

	return d.SaveAs(filename)
}
