#!/bin/bash

# Run tests
echo "Running tests..."
go test ./... -v

# Run tests with coverage
echo "Running tests with coverage..."
go test ./... -v -coverprofile=coverage.out

# Display coverage report
echo "Coverage report:"
go tool cover -html=coverage.out -o coverage.html
echo "Coverage report generated: coverage.html"
