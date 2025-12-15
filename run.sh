#!/bin/bash

# GLB Model Manager - Quick Start Script for Linux/Mac

echo "╔════════════════════════════════════════╗"
echo "║  GLB Model Manager - Quick Start      ║"
echo "╚════════════════════════════════════════╝"
echo ""
echo "Pilih opsi:"
echo "1. Run Backend"
echo "2. Run Frontend"
echo "3. Run Both (requires 2 terminals)"
echo "4. Setup Backend"
echo "5. Setup Frontend"
echo "6. Seed Database"
echo ""

read -p "Pilih (1-6): " choice

case $choice in
    1)
        echo "Starting Backend..."
        cd backend
        go run main.go
        ;;
    2)
        echo "Starting Frontend..."
        cd frontend
        npm run dev
        ;;
    3)
        echo "Opening 2 terminals..."
        echo "Pastikan sudah install Go dan Node.js!"
        open -a Terminal backend/
        sleep 2
        open -a Terminal frontend/
        echo "Running: cd backend && go run main.go"
        echo "Running: cd frontend && npm run dev"
        echo "Both servers started!"
        ;;
    4)
        echo "Setting up Backend..."
        cd backend
        echo "Installing Go dependencies..."
        go mod download
        echo "Backend setup complete!"
        ;;
    5)
        echo "Setting up Frontend..."
        cd frontend
        echo "Installing Node dependencies..."
        npm install
        echo "Frontend setup complete!"
        ;;
    6)
        echo "Seeding Database..."
        cd backend
        go run main.go seed.go
        ;;
    *)
        echo "Invalid choice!"
        ;;
esac
