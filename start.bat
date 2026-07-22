@echo off
start "Backend" cmd /c "cd /d %~dp0backend && go run ./cmd/server/"
start "Frontend" cmd /c "cd /d %~dp0frontend && npm run dev"
start "ML Service" cmd /c "cd /d %~dp0ml && .venv\Scripts\activate && uvicorn main:app --port 8000"
echo Backend, Frontend, and ML Service starting...
timeout /t 5
start http://localhost:5173