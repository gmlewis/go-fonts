#!/bin/bash -ex
echo "## Font samples" >> README.md
for i in $(ls fonts | grep -v .go) ; do
    echo $i
    sed -e "s/latoregular/${i}/ig" < cmd/render-fonts/main.go > main-${i}.go
    time go run main-${i}.go -msg "Sample from ${i}: 0123456789" -out "images/sample_${i}.png"
    echo "[![${i}](images/sample_${i}.png)](fonts/${i})" >> README.md
    rm main-${i}.go
done
