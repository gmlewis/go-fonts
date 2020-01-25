#!/bin/bash -ex
# -*- compile-command: "./gen-one-font-sample.sh"; -*-
sed -e "s/latoregular/$@/g" < cmd/render-fonts/main.go > main-$@.go
time go run main-$@.go -msg "Sample from $@: 0123456789" -out "images/sample_$@.png"
# rm main-$@.go
