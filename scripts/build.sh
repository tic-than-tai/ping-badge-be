#!/bin/bash

# Build the application
echo "Building PingBadge backend..."
go build -o bin/ping-badge-be main.go

# Check if build was successful
if [ $? -eq 0 ]; then
    echo "Build successful! Binary created at bin/ping-badge-be"
else
    echo "Build failed!"
    exit 1
fi
