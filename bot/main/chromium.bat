@echo off
:: Change directory to where chrome.exe is located
cd /d "C:\Program Files\Chromium\Application"

:: Start Chromium with the remote debugging port
start "" "chrome.exe" --remote-debugging-port=8989
