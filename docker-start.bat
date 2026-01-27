@echo off
chcp 65001 >nul
echo ========================================
echo Docker Deployment - Hotel Luggage System
echo ========================================
echo.

echo [1/4] Checking Docker installation...
docker --version >nul 2>&1
if errorlevel 1 (
    echo Error: Docker is not installed or not running!
    echo Please install Docker Desktop first.
    pause
    exit /b 1
)
echo OK: Docker is ready
echo.

echo [2/4] Checking docker-compose...
docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo Error: docker-compose is not available!
    pause
    exit /b 1
)
echo OK: docker-compose is ready
echo.

echo [3/4] Starting services (MySQL + Redis + App)...
docker-compose up -d
echo.

echo [4/4] Waiting for services to be ready...
timeout /t 5 /nobreak >nul
echo.

echo ========================================
echo Checking service status...
echo ========================================
docker-compose ps
echo.

echo ========================================
echo Viewing application logs...
echo ========================================
echo Press Ctrl+C to stop viewing logs
echo.
timeout /t 2 /nobreak >nul

docker-compose logs -f app
