#!/bin/bash -ex
for i in $(ls fonts | grep -v .go) ; do echo $i ; sed -e "s/latoregular/${i}/ig" < cmd/render-fonts/main.go > main-${i}.go ; go run main-${i}.go -all -out fonts/${i}/${i}.png ; rm main-${i}.go ; done
