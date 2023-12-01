#!/bin/bash -ex
# -*- compile-command: "./gen-fonts.sh"; -*-
go run cmd/font2lua/*.go fonts/*/*.svg
