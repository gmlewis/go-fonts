#!/bin/bash -ex
protoc --go_out=plugins=grpc:. *.proto
