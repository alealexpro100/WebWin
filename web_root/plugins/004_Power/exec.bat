@echo off
IF "%1"=="poweroff" (
    shutdown /s
)
IF "%1"=="force" (
    shutdown /s /f
)
if "%1"=="reboot" (
    shutdown /r
)
if errorlevel 1 exit /b %errorlevel%