#!/bin/bash -ex
# -*- compile-command: "./gen-latoregular-irmf.sh"; -*-
go run ./cmd/font2irmf/... \
    $@ \
    fonts/latoregular/Lato-Reg-webfont.svg
