# Atomation ReStreamer v1.1.0

## Major Features

- **Dark Neomorphic UI** - Modern, clean interface with soft shadows and gradient effects
- **Stream Health Monitoring** - Real-time connection tracking with 4 states:
  - Healthy (green) - Stream running smoothly
  - Degraded (yellow) - Connection issues detected
  - Reconnecting (orange) - Attempting to reconnect
  - Disconnected (red) - Stream stopped
- **Auto-Reconnect** - Automatic reconnection with up to 3 retry attempts and 5-second delays between attempts
- **Quality Presets** - Four streaming quality options optimized for different use cases:
  - **Low** (540p, 1000 kbps) - For slow connections
  - **Medium** (720p, 2000 kbps) - Balanced quality
  - **High** (1080p, 3000 kbps) - Default, great quality
  - **Ultra** (1080p, 5000 kbps) - Maximum quality
- **Smart Stop Confirmation** - Warns before stopping streams with significant time remaining (>5 seconds)
- **Stream Timer** - Real-time display of elapsed, remaining, and total time

## Improvements

- **Console Windows Hidden** - FFmpeg and FFprobe no longer show console windows for a cleaner user experience
- **Error Handling** - Suppressed false error messages when manually stopping streams
- **UI Layout** - Two-column responsive layout with improved spacing and visual hierarchy
- **Color-Coded Status** - Health indicators use distinct colors with animated glowing effects
- **Batch File Reliability** - Improved Wails executable path detection in build scripts

## Bug Fixes

- Fixed timer display for remaining and total time
- Fixed false error messages appearing on manual stream stop
- Improved auto-stop behavior when video completes

## Download

Download **Atomation-ReStreamer.exe** below and run it. FFmpeg binaries are included in the executable.

## Installation

1. Download `Atomation-ReStreamer.exe` from the assets below
2. Run the executable - no installation required!
3. Select your video file, select your streaming service, enter your stream key, and start streaming

## Requirements

- Windows 10/11
- Active internet connection for streaming

## Support

- **Issues**: Report bugs or request features in the Issues tab
- **Documentation**: See [README.md](https://github.com/atomicsunset/ReStreamer/blob/main/README.md)
- **Changelog**: Full version history in [CHANGELOG.md](https://github.com/atomicsunset/ReStreamer/blob/main/CHANGELOG.md)

---

---

**Created by Atomation**
