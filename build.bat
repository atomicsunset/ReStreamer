@echo off
echo ========================================
echo Building ReStreamer
echo ========================================
echo.

REM Set the Wails executable path
set WAILS_EXE=C:\Users\ohrep\go\bin\wails.exe

echo Checking if Wails is installed...
if not exist "%WAILS_EXE%" (
    echo ERROR: Wails not found at %WAILS_EXE%
    echo Trying to find wails in PATH...
    where wails >nul 2>&1
    if %errorlevel% neq 0 (
        echo ERROR: Wails is not installed!
        echo Please install Wails first:
        echo   go install github.com/wailsapp/wails/v2/cmd/wails@latest
        echo.
        pause
        exit /b 1
    )
    set WAILS_EXE=wails
)

"%WAILS_EXE%" version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Wails executable found but failed to run!
    pause
    exit /b 1
)

echo Wails found and working.
echo.

echo Downloading Go dependencies...
go mod download
if %errorlevel% neq 0 (
    echo ERROR: Failed to download dependencies
    pause
    exit /b 1
)

echo.
echo Building application...
"%WAILS_EXE%" build

if %errorlevel% neq 0 (
    echo.
    echo ERROR: Build failed!
    pause
    exit /b 1
)

echo.
echo ========================================
echo Build completed successfully!
echo ========================================
echo.
echo Your application is ready in: build\bin\Atomation-ReStreamer.exe
echo.
pause
