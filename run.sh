#!/bin/bash


echo "Starting services with docker-compose..."
docker-compose up -d

echo "Waiting for the containers to start..."
sleep 4

cleanup() {
    echo "Shutting down containers..."
    docker-compose down
}

echo "Running migration..."
MIGRATE=true go run app/cmd/main.go 

trap cleanup SIGINT
echo "Running main app..."
go run app/cmd/main.go