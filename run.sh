#!/bin/bash

MIGRATE=false
RUN=false

for arg in "$@"; do
    case $arg in
        --migrate)
        MIGRATE=true
        ;;
        --run)
        RUN=true
        ;;
    esac
done

echo "Starting services with docker-compose..."
docker-compose up -d

echo "Waiting for the containers to start..."
sleep 7

cleanup() {
    echo "Shutting down containers..."
    docker-compose down
}

if [ "$MIGRATE" = true ]; then
    echo "Running migration..."
    MIGRATE=true go run app/cmd/main.go 
fi

if [ "$RUN" = true ]; then
    trap cleanup SIGINT
    echo "Running main app..."
    go run app/cmd/main.go
fi