#!/bin/bash -ex
# -*- compile-command: "./render-one-font.sh"; -*-

echo $@
sed -e "s/latoregular/$@/g" < cmd/render-fonts/main.go > main-$@.go
time go run main-$@.go -all -out fonts/$@/$@.png
# rm main-$@.go
