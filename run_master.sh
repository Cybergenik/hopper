#!/bin/bash
#go build -race -buildmode=plugin ./mut.go
go run main.go -I ./test/in/ -H=5
