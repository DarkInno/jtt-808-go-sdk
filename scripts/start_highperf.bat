@echo off
echo Starting JT808 High Performance Server...
cd /d "%~dp0\.."
go run cmd/highperf/main.go
pause
