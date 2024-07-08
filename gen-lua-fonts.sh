#!/bin/bash -ex
# -*- compile-command: "./gen-fonts.sh"; -*-
go run cmd/font2lua/main.go fonts/*/*.svg ../go-fonts-*/fonts/*/*.svg
