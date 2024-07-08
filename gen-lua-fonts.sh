#!/bin/bash -ex
# -*- compile-command: "./gen-fonts.sh"; -*-
go run cmd/font2lua/main.go ../go-fonts-*/fonts/*/*.svg
