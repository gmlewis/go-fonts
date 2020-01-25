#!/bin/bash -ex
# -*- compile-command: "./gen-latoregular.sh"; -*-
go run cmd/font2go/*.go \
   fonts/latoregular/Lato-Reg-webfont.svg \
   $@
