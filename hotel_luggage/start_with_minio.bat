@echo off

REM Change to the directory where this bat file is located
cd /d %~dp0

REM Database configuration
set DB_DSN=root:123456@tcp(127.0.0.1:3306)/hotel_luggage?charset=utf8mb4^&parseTime=True^&loc=Local

REM Redis configuration
set REDIS_ADDR=127.0.0.1:6379
set REDIS_PASSWORD=
set REDIS_DB=0

REM MinIO configuration
set MINIO_ENDPOINT=minio.2huo.tech
set MINIO_ACCESS_KEY=i8IuD8lJYxE5kAL1HOwS
set MINIO_SECRET_KEY=lAfdJNMqAQDmNrK8peuIwu5un6PFI0EtgWlae7jv
set MINIO_USE_SSL=true
set MINIO_BUCKET_NAME=traning-hotel

REM JWT secret
set JWT_SECRET=your-secret-key-change-in-production

echo Starting Hotel Luggage System...
echo Database: hotel_luggage@127.0.0.1:3306
echo Redis: %REDIS_ADDR%
echo MinIO: %MINIO_ENDPOINT%
echo Bucket: %MINIO_BUCKET_NAME%
echo.

go run cmd\main.go

pause
