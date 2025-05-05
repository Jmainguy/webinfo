#!/bin/bash

# Build the WebAssembly binary
cd wasm
GOOS=js GOARCH=wasm go build -o main.wasm
cd ..

# Build the server binary for ARM architecture
cd server
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o server
cd ..

# Build the Docker image for ARM
docker build --platform linux/arm64 -t zot.soh.re/info:latest .
docker push zot.soh.re/info:latest
