# OH.Report ReStreamer - Deployment Guide

## Is it 100% Standalone?

**ABSOLUTELY YES!** The application is a **SINGLE .EXE FILE** - truly standalone!

### What's Included:
- **OH.Report-ReStreamer.exe** (209 MB) - Everything in one file!
  - Application UI and logic
  - FFmpeg video processing (embedded)
  - FFprobe video analysis (embedded)

**Single file deployment:** Just one .exe file!

### What Users Need:

#### ✅ Already Have (Windows 10/11):
- **WebView2 Runtime** - Pre-installed on Windows 10 (version 2004+) and all Windows 11

#### ⚠️ May Need (Older Windows):
- **WebView2 Runtime** for Windows 7, 8, 8.1, or old Windows 10
  - Download: https://developer.microsoft.com/microsoft-edge/webview2/
  - Installer: ~100KB download that fetches the runtime

### Distribution Options:

#### Single File Distribution (Recommended) ⭐
1. Copy `build/bin/OH.Report-ReStreamer.exe`
2. Send the single .exe file to users
3. Users double-click to run - that's it!

**No ZIP needed! No folder structure! Just ONE .exe file!**

#### How It Works:
- FFmpeg and FFprobe are embedded inside the .exe
- On first run, they're extracted to: `%TEMP%\OH.Report-ReStreamer\`
- Automatic cleanup when app closes
- No clutter in user's directories!

### Zero Configuration Required

Users literally just:
1. Download `OH.Report-ReStreamer.exe`
2. Double-click it
3. Done!

**That's it! No:**
- ❌ Installation
- ❌ Extraction
- ❌ PATH setup
- ❌ Admin rights
- ❌ Other files needed
- ❌ Configuration

**Just download and run!**

## Building from Source

```bash
# Simple build - everything is embedded automatically!
wails build

# The standalone exe is ready in:
# build/bin/OH.Report-ReStreamer.exe
```

**Note:** Make sure `ffmpeg.exe` and `ffprobe.exe` are in the project root before building. They will be automatically embedded into the final executable.

## Professional Features

- ✅ Professional dark theme (no more Lisa Frank!)
- ✅ Multi-platform support (Castr/Facebook/YouTube)
- ✅ M3U playlist support
- ✅ Real-time runtime timer
- ✅ Auto-stop when video ends
- ✅ Live streaming progress tracking
- ✅ Standalone executable

## System Requirements

- **OS:** Windows 7+ (WebView2 required for Win7/8)
- **RAM:** 100MB minimum
- **Disk:** 250MB for application
- **Network:** Internet connection for streaming

## Troubleshooting

### "Application won't start"
- **Windows 7/8 users:** Install WebView2 Runtime
- **All users:** Make sure all 3 exe files are together

### "FFmpeg not found"
- Ensure `ffmpeg.exe` is in the same folder as the main application

### "Stream won't connect"
- Verify RTMP URL and stream key from your platform
- Check your internet connection
- Test with a short video first

## License & Credits

- Built with [Wails](https://wails.io/)
- Powered by [FFmpeg](https://ffmpeg.org/)
- © 2025 OH.Report
