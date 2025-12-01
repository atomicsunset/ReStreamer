# OH.Report ReStreamer - Developer Documentation

## Project Overview

OH.Report ReStreamer is a standalone Windows application built with Wails v2 that streams pre-recorded videos to RTMP platforms without requiring vMix or OBS. The application embeds FFmpeg for video processing and uses WebView2 for the frontend.

## Technology Stack

- **Framework**: Wails v2.11.0
- **Backend**: Go 1.25.4
- **Frontend**: Vanilla HTML/CSS/JavaScript
- **Video Processing**: FFmpeg (embedded)
- **UI Runtime**: WebView2 (built into Windows 10/11)

## Project Structure

```
ReStreamer/
├── main.go                 # Application entry point
├── app.go                  # Main application logic and Wails bindings
├── streamer.go            # FFmpeg streaming functionality
├── embedded.go            # Embedded binary extraction logic
├── wails.json             # Wails configuration
├── go.mod / go.sum        # Go dependencies
├── frontend/
│   ├── dist/
│   │   ├── index.html     # Main UI (two-column layout)
│   │   ├── style.css      # Styling
│   │   ├── app.js         # Frontend logic
│   │   └── wailsjs/       # Auto-generated Wails bindings
│   └── wailsjs/           # Go-to-JS bindings
├── build/                 # Build output (gitignored)
│   └── bin/
│       └── OH.Report-ReStreamer.exe
├── dist/                  # Distribution folder
│   └── OH.Report-ReStreamer.exe
├── build.bat             # Build script
├── deploy.bat            # Build + copy FFmpeg
├── run-dev.bat           # Development mode
└── download-ffmpeg.bat   # FFmpeg download helper
```

## Architecture

### Backend (Go)

#### `main.go`
- Entry point for the application
- Configures Wails runtime with window size (800x600)
- Embeds frontend assets

#### `app.go`
- `App` struct: Main application state
- `startup()`: Initializes app, extracts embedded FFmpeg binaries
- `shutdown()`: Cleanup, stops active streams
- `OpenFileDialog()`: File picker for video selection
- `StartStream()`: Initiates FFmpeg streaming process
- `StopStream()`: Terminates active stream
- `GetStreamStatus()`: Returns current stream state
- `SetAutoStopOnEnd()`: Configure auto-stop behavior

#### `streamer.go`
- `Streamer` struct: Manages FFmpeg process
- `Start()`: Launches FFmpeg with RTMP streaming parameters
- `Stop()`: Gracefully terminates FFmpeg
- `GetStatus()`: Returns streaming metrics (elapsed, remaining, total time)
- `getVideoDuration()`: Uses ffprobe to get video length
- `getFFmpegPath()`: Resolves FFmpeg binary location

#### `embedded.go`
- Embeds ffmpeg.exe and ffprobe.exe into the compiled binary
- `ExtractEmbeddedBinaries()`: Extracts to temp directory on startup
- `CleanupEmbeddedBinaries()`: Removes temp files on shutdown

### Frontend (JavaScript)

#### Key Components

1. **Service Selection Logic**
   - Castr: Shows channel selector, hides RTMP/key fields
   - Facebook/YouTube/Custom: Shows RTMP URL and stream key fields

2. **Castr Channel Keys**
   - Predefined mapping of 9 channels (000-008) to stream keys
   - Auto-populates when channel is selected

3. **Stream Status Polling**
   - Polls backend every 1 second while streaming
   - Updates timer display and status indicator

4. **UI State Management**
   - Disables inputs during streaming
   - Shows/hides timer based on stream state
   - Error message display

## Building from Source

### Prerequisites

1. **Go 1.25+**
   ```bash
   go version
   ```

2. **Wails CLI**
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

3. **FFmpeg Binaries**
   - Place `ffmpeg.exe` and `ffprobe.exe` in the root directory
   - Or run `download-ffmpeg.bat`

### Build Commands

**Development Mode:**
```bash
wails dev
# Or: run-dev.bat
```

**Production Build:**
```bash
wails build
# Or: build.bat
```

**Build + Package:**
```bash
deploy.bat
```
- Builds the application
- Copies FFmpeg binaries to build/bin/
- Creates standalone distribution

## Development Workflow

### Adding New Features

1. **Backend Changes (Go)**
   - Add methods to `app.go` or `streamer.go`
   - Methods must be exported (capitalized) to be callable from frontend
   - Update return types in `StreamStatus` struct if needed

2. **Frontend Changes**
   - Modify `frontend/dist/` files
   - Wails automatically embeds these on build
   - Call Go methods via `window.go.main.App.MethodName()`

3. **Rebuild**
   - Run `wails build` to embed frontend changes
   - Frontend is compiled into the binary

### FFmpeg Parameters

Current encoding settings (in `streamer.go:124-140`):
```go
"-re",              // Read at native frame rate
"-i", videoPath,    // Input file
"-c:v", "libx264",  // H.264 video codec
"-preset", "veryfast", // Encoding speed
"-maxrate", "3000k",   // Max bitrate
"-bufsize", "6000k",   // Buffer size
"-pix_fmt", "yuv420p", // Pixel format
"-g", "50",            // GOP size
"-c:a", "aac",         // AAC audio codec
"-b:a", "128k",        // Audio bitrate
"-ar", "44100",        // Sample rate
"-f", "flv",           // FLV output format
"-flvflags", "no_duration_filesize",
rtmpURL,               // Destination
```

To adjust quality/performance, modify these parameters.

## Configuration

### Wails Configuration (`wails.json`)

```json
{
  "name": "OH.Report ReStreamer",
  "outputfilename": "OH.Report-ReStreamer",
  "frontend:install": "echo 'No frontend install needed'",
  "frontend:build": "echo 'Using pre-built frontend'",
  "info": {
    "productVersion": "1.0.0"
  }
}
```

### Window Size
Defined in `main.go:23-24`:
```go
Width:  800,
Height: 600,
```

## Castr Channel Configuration

To add/modify Castr channels, edit `frontend/dist/app.js`:

```javascript
const castrChannelKeys = {
    '000': 'live_xxxxx?password=xxxxx',
    '001': 'live_xxxxx?password=xxxxx',
    // Add more channels here
};
```

Also update the HTML dropdown in `frontend/dist/index.html`.

## Debugging

### Enable Wails Dev Tools

In `main.go`, change build mode to include dev tools:
```go
err := wails.Run(&options.App{
    // ... other options
    Debug: options.Debug{
        OpenInspectorOnStartup: true,
    },
})
```

### Check FFmpeg Output

FFmpeg stderr is captured in `streamer.go:166-189`. Check logs for encoding issues.

### Common Issues

1. **Stream not starting**: Check RTMP URL format and stream key
2. **FFmpeg not found**: Ensure binaries are embedded (check `embedded.go`)
3. **Video not playing**: Verify video codec compatibility (use H.264/AAC)
4. **Build fails**: Ensure Wails CLI and Go versions match requirements

## Deployment

### Distribution

The final executable is completely standalone:
- Single file: `dist/OH.Report-ReStreamer.exe` (~199MB)
- Includes embedded FFmpeg and ffprobe
- No installation required
- Works on Windows 10/11 (WebView2 pre-installed)

### Deployment Checklist

- [ ] Build production binary (`wails build`)
- [ ] Test on clean Windows 10/11 system
- [ ] Verify FFmpeg extraction works (check temp directory)
- [ ] Test all streaming services (Castr, Facebook, YouTube)
- [ ] Confirm M3U playlist support
- [ ] Validate timer accuracy

## Performance Considerations

- FFmpeg process runs as separate child process
- CPU usage depends on video encoding settings
- Memory: ~50-100MB for app + FFmpeg overhead
- Temp files cleaned up on exit (FFmpeg binaries extracted to `%TEMP%\OH.Report-ReStreamer`)

## Security Notes

- Stream keys stored in JavaScript (visible in source)
- For production, consider encrypting keys or using backend storage
- RTMPS (secure RTMP) supported for Facebook Live

## Future Enhancement Ideas

- YouTube API integration for metadata setup
- Facebook API integration
- Multiple simultaneous streams
- Stream health monitoring
- Recording while streaming
- Scheduled streaming
- Stream presets/profiles

## License & Copyright

Copyright © 2025 OH.Report
Product Version: 1.0.0

## Support

For issues or questions, refer to:
- GUIDE.md - User documentation
- README.md - Project overview
- DEPLOYMENT.md - Deployment instructions

---

**Last Updated**: 2025-01-30
**Version**: 1.0.0
