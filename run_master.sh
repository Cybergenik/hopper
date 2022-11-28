#!/bin/bash
#go build -race -buildmode=plugin ./mut.go
go run -race main.go -I ./test/in/ -H=5
