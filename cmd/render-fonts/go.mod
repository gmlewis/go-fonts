module github.com/gmlewis/go-fonts/cmd/make-puzzle

go 1.22.4

require (
	github.com/gmlewis/go-fonts v0.19.0
	github.com/gmlewis/go-fonts/fonts/latoregular v0.0.0-20240626233958-3409c190883f
)

require (
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/gmlewis/go3d v0.0.4 // indirect
	github.com/gmlewis/ponoko2d v0.0.0-20190404133045-d77d370bec9a // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/yofu/dxf v0.0.0-20190710012328-5a6d1e83f16c // indirect
	golang.org/x/image v0.18.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/gmlewis/go-fonts v0.19.0 => ../../

replace github.com/gmlewis/go-fonts-l/fonts/latoregular v0.1.0 => ../../../go-fonts-l/fonts/latoregular
