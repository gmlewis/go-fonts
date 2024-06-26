#!/bin/bash -e
for i in $(find . -name go.mod); do
    echo $i
    pushd ${i%go.mod} && go mod tidy && popd
done
