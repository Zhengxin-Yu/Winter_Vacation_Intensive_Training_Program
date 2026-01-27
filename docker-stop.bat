@echo off
chcp 65001 >nul
echo ========================================
echo Stopping Docker Services
echo ========================================
echo.

echo Stopping all services...
docker-compose stop
echo.

echo Service status:
docker-compose ps
echo.

echo ========================================
echo All services stopped successfully!
echo ========================================
echo.
echo To remove containers: docker-compose down
echo To remove containers and data: docker-compose down -v
echo.

pause
