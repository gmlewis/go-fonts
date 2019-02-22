package fonts

import (
	"github.com/yofu/dxf"
)

func (r *Render) SaveDXF(filename string, scale float64) error {
	d := dxf.NewDrawing()
	d.Header().LtScale = 100.0
	d.AddLayer("Lines", dxf.DefaultColor, dxf.DefaultLineType, true)
	d.ChangeLayer("Lines")
	for _, poly := range r.Polygons {
		lastX, lastY := scale*(poly.Pts[0][0]-r.MBB.Min[0]), scale*(poly.Pts[0][1]-r.MBB.Min[1])
		for _, pt := range poly.Pts[1:] {
			x, y := scale*(pt[0]-r.MBB.Min[0]), scale*(pt[1]-r.MBB.Min[1])
			d.Line(lastX, lastY, 0.0, x, y, 0.0)
			lastX, lastY = x, y
		}
		d.Line(lastX, lastY, 0.0, scale*(poly.Pts[0][0]-r.MBB.Min[0]), scale*(poly.Pts[0][1]-r.MBB.Min[1]), 0.0)
	}

	return d.SaveAs(filename)
}
