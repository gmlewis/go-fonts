module github.com/gmlewis/go-fonts

go 1.22.4

require (
	github.com/fogleman/gg v1.3.0
	github.com/gmlewis/go-fonts-b/fonts/baloo v0.1.0
	github.com/gmlewis/go-fonts/fonts/latoregular v0.0.0-20240626233958-3409c190883f
	github.com/gmlewis/go3d v0.0.4
	github.com/gmlewis/ponoko2d v0.0.0-20190404133045-d77d370bec9a
	github.com/yofu/dxf v0.0.0-20190710012328-5a6d1e83f16c
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	golang.org/x/image v0.18.0 // indirect
)

replace github.com/gmlewis/go-fonts-b/fonts/baloo v0.1.0 => ../go-fonts-b/fonts/baloo

replace github.com/gmlewis/go-fonts-l/fonts/latoregular v0.1.0 => ../go-fonts-l/fonts/latoregular
