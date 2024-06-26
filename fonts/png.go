package fonts

import (
	"errors"
	"image"
	"log"

	"github.com/fogleman/gg"
)

func GetRenderMBB(inWidth, inHeight int, renders ...*Render) (mbb MBB, scale float64, width, height int, err error) {
	if len(renders) == 0 {
		return mbb, 0, 0, 0, errors.New("must have at least one render")
	}

	for i, render := range renders {
		if i == 0 {
			mbb = render.MBB
		} else {
			mbb.Join(&render.MBB)
		}
	}

	scale, width, height = maximize(&mbb, inWidth, inHeight)
	return mbb, scale, width, height, nil
}

func SavePNG(filename string, inWidth, inHeight int, renders ...*Render) error {
	mbb, scale, width, height, err := GetRenderMBB(inWidth, inHeight, renders...)
	if err != nil {
		return err
	}
	log.Printf("MBB: %v, scale=%.2f, size=(%v,%v)", mbb, scale, width, height)

	dc := gg.NewContext(width, height)
	background(dc, renders[0])
	dc.Clear()
	for _, r := range renders {
		for _, poly := range r.Polygons {
			if poly.Dark {
				foreground(dc, r)
			} else {
				background(dc, r)
			}
			for i, pt := range poly.Pts {
				if i == 0 {
					dc.MoveTo(scale*(pt[0]-mbb.Min[0]), float64(height)-scale*(pt[1]-mbb.Min[1]))
				} else {
					dc.LineTo(scale*(pt[0]-mbb.Min[0]), float64(height)-scale*(pt[1]-mbb.Min[1]))
				}
			}
			dc.Fill()
		}
	}
	return dc.SavePNG(filename)
}

func background(dc *gg.Context, render *Render) {
	const max = 0xffff
	if render.Background != nil {
		r, g, b, a := render.Background.RGBA()
		dc.SetRGBA(float64(r)/max, float64(g)/max, float64(b)/max, float64(a)/max)
	} else {
		dc.SetRGB(1, 1, 1)
	}
}

func foreground(dc *gg.Context, render *Render) {
	const max = 0xffff
	if render.Foreground != nil {
		r, g, b, a := render.Foreground.RGBA()
		dc.SetRGBA(float64(r)/max, float64(g)/max, float64(b)/max, float64(a)/max)
	} else {
		dc.SetRGB(0, 0, 0)
	}
}

func maximize(mbb *MBB, width, height int) (float64, int, int) {
	scale := float64(width) / (mbb.Max[0] - mbb.Min[0])
	if yScale := float64(height) / (mbb.Max[1] - mbb.Min[1]); yScale < scale {
		scale = yScale
		width = int(0.5 + scale*(mbb.Max[0]-mbb.Min[0]))
	} else {
		height = int(0.5 + scale*(mbb.Max[1]-mbb.Min[1]))
	}
	return scale, width, height
}

func (r *Render) ToImage(width, height int) image.Image {
	dc := r.ToContext(width, height)
	return dc.Image()
}

func (r *Render) SavePNG(filename string, width, height int) error {
	dc := r.ToContext(width, height)
	return dc.SavePNG(filename)
}

func (r *Render) ToContext(width, height int) *gg.Context {
	var scale float64
	scale, width, height = maximize(&r.MBB, width, height)
	log.Printf("MBB: %v, scale=%.2f, size=(%v,%v)", r.MBB, scale, width, height)

	dc := gg.NewContext(width, height)
	background(dc, r)
	dc.Clear()
	r.RenderToDC(dc, 0, 0, scale, height)
	return dc
}

func (r *Render) RenderToDC(dc *gg.Context, dx, dy, scale float64, height int) {
	for _, poly := range r.Polygons {
		if poly.Dark {
			foreground(dc, r)
		} else {
			background(dc, r)
		}
		for i, pt := range poly.Pts {
			if i == 0 {
				dc.MoveTo(dx+scale*(pt[0]-r.MBB.Min[0]), dy+float64(height)-scale*(pt[1]-r.MBB.Min[1]))
			} else {
				dc.LineTo(dx+scale*(pt[0]-r.MBB.Min[0]), dy+float64(height)-scale*(pt[1]-r.MBB.Min[1]))
			}
		}
		dc.Fill()
	}
}
