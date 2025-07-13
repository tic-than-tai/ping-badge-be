#!/bin/bash

# Run database migrations
echo "Running database migrations..."

# Check if migrate tool is installed
if ! command -v migrate &> /dev/null; then
    echo "golang-migrate tool not found. Installing..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
fi

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | xargs)
fi

# Run migrations
migrate -path migrations -database "$DATABASE_URL" up

echo "Migrations completed!"
