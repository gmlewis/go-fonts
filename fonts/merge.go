package fonts

func Merge(renders ...*Render) *Render {
	merged := &Render{}
	for i, r := range renders {
		if i == 0 {
			merged.MBB = r.MBB
		} else {
			merged.MBB.Join(&r.MBB)
		}
		merged.Polygons = append(merged.Polygons, r.Polygons...)
		merged.Info = append(merged.Info, r.Info...)
	}
	return merged
}
