#!/bin/bash -ex

echo $@
sed -e "s/latoregular/$@/ig" < cmd/render-fonts/main.go > main-$@.go
go run main-$@.go -all -out fonts/$@/$@.png
rm main-$@.go
