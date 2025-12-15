@echo off
REM ========================================
REM GLB Model Manager - Quick Start Script
REM ========================================

REM Check if folder exists
if not exist "backend" (
    echo Error: backend folder not found!
    exit /b 1
)

if not exist "frontend" (
    echo Error: frontend folder not found!
    exit /b 1
)

echo.
echo ╔════════════════════════════════════════╗
echo ║  GLB Model Manager - Quick Start      ║
echo ╚════════════════════════════════════════╝
echo.
echo Pilih opsi:
echo 1. Run Backend
echo 2. Run Frontend
echo 3. Run Both (requires 2 terminals)
echo 4. Setup Backend
echo 5. Setup Frontend
echo 6. Seed Database
echo.

set /p choice="Pilih (1-6): "

if "%choice%"=="1" (
    echo Starting Backend...
    cd backend
    go run main.go
    pause
) else if "%choice%"=="2" (
    echo Starting Frontend...
    cd frontend
    npm run dev
    pause
) else if "%choice%"=="3" (
    echo Opening 2 terminals...
    echo Pastikan sudah install Go dan Node.js!
    start cmd /k "cd /d %cd%\backend && go run main.go"
    timeout /t 2
    start cmd /k "cd /d %cd%\frontend && npm run dev"
    echo Both servers started!
    pause
) else if "%choice%"=="4" (
    echo Setting up Backend...
    cd backend
    echo Installing Go dependencies...
    go mod download
    echo Backend setup complete!
    pause
) else if "%choice%"=="5" (
    echo Setting up Frontend...
    cd frontend
    echo Installing Node dependencies...
    call npm install
    echo Frontend setup complete!
    pause
) else if "%choice%"=="6" (
    echo Seeding Database...
    cd backend
    go run main.go seed.go
    pause
) else (
    echo Invalid choice!
    pause
)
