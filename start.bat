@echo off
start "Backend" cmd /c "cd /d c:\aplikasi-analisis-keuangan\backend && go run ./cmd/server/"
start "Frontend" cmd /c "cd /d c:\aplikasi-analisis-keuangan\frontend && npm run dev"
echo Backend & Frontend starting...
timeout /t 5
start http://localhost:5173