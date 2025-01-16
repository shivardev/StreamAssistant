@echo off
:: Change directory to where chrome.exe is located
cd /d "C:\Program Files\Google\Chrome\Application"

:: Start Chrome with a specific profile and remote debugging port
start "" "chrome.exe" --remote-debugging-port=8989 --user-data-dir="C:\streaming"
