# Use Chainguard's distroless image as the base
FROM cgr.dev/chainguard/static:latest

# Set the working directory
WORKDIR /app

# Copy the pre-built server binary and static files
COPY server/server /app/server
COPY wasm/main.wasm /app/main.wasm
COPY index.html /app/index.html
COPY wasm_exec.js /app/wasm_exec.js

# Expose port 8080
EXPOSE 8080

# Run the Go server
CMD ["/app/server"]
