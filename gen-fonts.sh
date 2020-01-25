#!/bin/bash -ex
# -*- compile-command: "./gen-fonts.sh"; -*-
go run cmd/font2go/*.go fonts/*/*.svg
