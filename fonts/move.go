package fonts

import "fmt"

func (r *Render) MoveGlyph(index int, dx, dy float64) error {
	if index >= len(r.Info) {
		return fmt.Errorf("bad index %v; only %v glyphs found", index, len(r.Info))
	}

	r.Info[index].X += dx
	r.Info[index].Y += dy
	r.Info[index].MBB.Min[0] += dx
	r.Info[index].MBB.Max[0] += dx
	r.Info[index].MBB.Min[1] += dy
	r.Info[index].MBB.Max[1] += dy

	polys, err := r.GetPolygonsForGlyph(index)
	if err != nil {
		return err
	}

	for _, p := range polys {
		p.Move(dx, dy)
	}

	return nil
}

func (p *Polygon) Move(dx, dy float64) {
	p.MBB.Min[0] += dx
	p.MBB.Max[0] += dx
	p.MBB.Min[1] += dy
	p.MBB.Max[1] += dy
	for i := range p.Pts {
		p.Pts[i][0] += dx
		p.Pts[i][1] += dy
	}
}

func (r *Render) GetPolygonsForGlyph(index int) ([]*Polygon, error) {
	if index >= len(r.Info) {
		return nil, fmt.Errorf("bad index %v; only %v glyphs found", index, len(r.Info))
	}

	var result []*Polygon
	var pi int
	for i, gi := range r.Info {
		if i == index {
			for len(result) < gi.N {
				result = append(result, r.Polygons[pi])
				pi++
			}
			break
		}
		pi += gi.N
	}

	return result, nil
}
