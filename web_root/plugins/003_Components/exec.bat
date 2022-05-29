@echo off
powershell -NoLogo -NoProfile -ExecutionPolicy Bypass -File %~dp0\exec.ps1 %*
if errorlevel 1 exit /b %errorlevel%