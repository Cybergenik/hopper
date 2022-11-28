#!/bin/bash
#go build -race -buildmode=plugin ./mut.go
go build .
./hopper -I ./test/in/ -H=5
