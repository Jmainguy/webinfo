#!/bin/bash

# Build the WebAssembly binary
cd wasm
GOOS=js GOARCH=wasm go build -o main.wasm
cd ..

# Build the server binary statically
cd server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server
cd ..

# Build the Docker image
docker build -t web-info .

# Run the Docker container
docker run -p 8080:8080 web-info
