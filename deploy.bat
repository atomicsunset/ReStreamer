@echo off
echo ========================================
echo Atomation ReStreamer - Deployment
echo ========================================
echo.

REM Set the Wails executable path
set WAILS_EXE=C:\Users\ohrep\go\bin\wails.exe

echo Checking if Wails is available...
if not exist "%WAILS_EXE%" (
    echo ERROR: Wails not found at %WAILS_EXE%
    echo Trying to find wails in PATH...
    where wails >nul 2>&1
    if %errorlevel% neq 0 (
        echo ERROR: Wails is not installed!
        pause
        exit /b 1
    )
    set WAILS_EXE=wails
)

:: Build the application
echo Building application...
"%WAILS_EXE%" build
if errorlevel 1 (
    echo ERROR: Build failed!
    pause
    exit /b 1
)

:: Copy FFmpeg dependencies
echo.
echo Copying FFmpeg dependencies...
copy /Y ffmpeg.exe build\bin\ffmpeg.exe
copy /Y ffprobe.exe build\bin\ffprobe.exe

echo.
echo ========================================
echo Deployment Complete!
echo ========================================
echo.
echo The standalone application is ready in: build\bin\
echo.
echo Contents:
echo - Atomation-ReStreamer.exe  (Main application)
echo - ffmpeg.exe                (Video processing)
echo - ffprobe.exe               (Video analysis)
echo.
echo NOTE: Windows 10/11 users have WebView2 pre-installed.
echo For Windows 7/8, users need to install WebView2 Runtime from:
echo https://developer.microsoft.com/en-us/microsoft-edge/webview2/
echo.
pause
