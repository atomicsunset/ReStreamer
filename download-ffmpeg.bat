@echo off
echo ========================================
echo FFmpeg Downloader for ReStreamer
echo ========================================
echo.
echo This script will download FFmpeg for Windows
echo.

set FFMPEG_URL=https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip
set DOWNLOAD_DIR=%~dp0
set ZIP_FILE=%DOWNLOAD_DIR%ffmpeg.zip
set EXTRACT_DIR=%DOWNLOAD_DIR%ffmpeg-temp

echo Downloading FFmpeg...
echo This may take a few minutes depending on your connection.
echo.

:: Download using PowerShell
powershell -Command "& {Invoke-WebRequest -Uri '%FFMPEG_URL%' -OutFile '%ZIP_FILE%'}"

if not exist "%ZIP_FILE%" (
    echo ERROR: Download failed!
    pause
    exit /b 1
)

echo.
echo Download complete. Extracting...

:: Extract using PowerShell
powershell -Command "& {Expand-Archive -Path '%ZIP_FILE%' -DestinationPath '%EXTRACT_DIR%' -Force}"

:: Find and copy ffmpeg.exe and ffprobe.exe
for /r "%EXTRACT_DIR%" %%i in (ffmpeg.exe) do (
    copy "%%i" "%DOWNLOAD_DIR%"
    echo Copied ffmpeg.exe
)

for /r "%EXTRACT_DIR%" %%i in (ffprobe.exe) do (
    copy "%%i" "%DOWNLOAD_DIR%"
    echo Copied ffprobe.exe
)

:: Cleanup
del "%ZIP_FILE%"
rmdir /s /q "%EXTRACT_DIR%"

echo.
echo ========================================
echo FFmpeg installation complete!
echo ========================================
echo.
echo FFmpeg has been installed to the application directory.
echo You can now build and run ReStreamer.
echo.
pause
