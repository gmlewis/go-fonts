#!/bin/bash -ex
protoc --go_out=. ./glyphs.proto
mv github.com/gmlewis/go-fonts/pb/glyphs/glyphs.pb.go . && rm -rf github.com
