#!/bin/bash
export GOOS=windows 
export GOARCH=amd64 
export CC=x86_64-w64-mingw32-gcc
export CXX=x86_64-w64-mingw32-gcc 
export CGO_ENABLED=1
export CGO_CXXFLAGS="-static-libgcc -static-libstdc++ -Wl,-Bstatic -lstdc++ -lpthread -Wl,-Bdynamic"
go build -ldflags "-s -w -X main.version=$2" -o $3 -v -x $1