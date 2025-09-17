#!/bin/bash

rm -Rf build
mkdir -p build 
go build -ldflags="-s -w" -o build/markovian-word-gen .
