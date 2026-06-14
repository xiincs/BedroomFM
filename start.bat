@echo off
echo Starting BedroomFM...

echo [1/3] Starting backend on :8080
start "BedroomFM Backend" cmd /c "cd /d %~dp0backend && go run main.go"

echo [2/3] Starting frontend on :5173
start "BedroomFM Frontend" cmd /c "cd /d %~dp0frontend && npm run dev"

echo [3/3] Starting NetEase Music API on :3000
echo    (requires: npm install -g NeteaseCloudMusicApi)
start "NetEase API" cmd /c "NeteaseCloudMusicApi --port 3000 || npx NeteaseCloudMusicApi --port 3000"

timeout /t 3 /nobreak > nul
start http://localhost:5173
echo Done! App running at http://localhost:5173
