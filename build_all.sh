#!/bin/sh

rm -rf build/*
mkdir -p build/mac_amd64 build/win_amd64 build/linux_amd64
GOOS=darwin GOARCH=amd64 go build -o build/mac_amd64/platform_cc
GOOS=windows GOARCH=amd64 go build -o build/win_amd64/platform_cc.exe
GOOS=linux GOARCH=amd64 go build -o build/linux_amd64/platform_cc