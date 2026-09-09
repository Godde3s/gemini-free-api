@echo off
REM ============================================================================
REM  Gemini-Free-API - one-command launcher (Windows)
REM  First run:  double-click, or:  start.bat
REM ============================================================================
setlocal enabledelayedexpansion
cd /d "%~dp0"

if not exist .env (
    copy .env.example .env >nul
    echo.
    echo  ! -- First-run setup ------------------------------------------
    echo   A starter .env was created from .env.example.
    echo   Open it and paste your GEMINI_COOKIES (optional — guest works) into GEMINI_COOKIES:
    echo       notepad .env
    echo       (token: gemini.google.com -^> F12 -^> Application -^> Cookies -^> token)
    echo  ----------------------------------------------------------------
    echo.
)

if exist gemini-api.exe (
    set BIN=gemini-api.exe
    goto :run
)

where go >nul 2>nul
if %errorlevel%==0 (
    echo Building gemini-api with Go...
    go build -trimpath -ldflags="-s -w" -o gemini-api.exe .
    set BIN=gemini-api.exe
    goto :run
)

echo No prebuilt binary and no Go toolchain found.
echo Install Go: https://go.dev/dl/   then re-run start.bat
echo Or grab a release binary: https://github.com/Godde3s/gemini-free-api/releases
pause
exit /b 1

:run
for /f "tokens=1,2 delims==" %%a in (.env) do (
    if not "%%b"=="" set "%%a=%%b"
)
echo Starting Gemini-Free-API...
"%BIN%"
