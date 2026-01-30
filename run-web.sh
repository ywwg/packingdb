#!/bin/bash
# Run the packingweb server

cd "$(dirname "$0")/cmd/packingweb"

echo "Building packingweb..."
go build

if [ $? -eq 0 ]; then
    echo "Starting packingweb server..."
    echo ""
    echo "🌐 Open http://localhost:8080 in your browser"
    echo ""
    ./packingweb
else
    echo "Build failed!"
    exit 1
fi
