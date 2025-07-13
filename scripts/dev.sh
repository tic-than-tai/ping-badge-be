#!/bin/bash

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | xargs)
fi

# Run the application
echo "Starting PingBadge backend server..."
go run main.go
