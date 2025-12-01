# Changelog

All notable changes to OH.Report ReStreamer will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2025-01-30

### Added
- **Dark Neomorphic UI** - Complete visual overhaul with modern dark theme and soft neumorphic shadows
- **Stream Health Monitoring** - Real-time connection health tracking with 4 states (healthy, degraded, reconnecting, disconnected)
- **Auto-Reconnect** - Automatic reconnection with up to 3 retry attempts and 5-second delays between attempts
- **Smart Stop Confirmation** - Confirmation dialog when manually stopping a stream with significant time remaining (>5 seconds)
- **Quality Presets** - Four streaming quality options: Low (540p), Medium (720p), High (1080p), Ultra (1080p high bitrate)

### Changed
- **Error Handling** - Suppressed error messages when manually stopping streams (expected behavior)
- **Console Windows** - Hidden FFmpeg and FFprobe console windows for cleaner user experience
- **UI Layout** - Two-column responsive layout with improved spacing and visual hierarchy
- **Color-coded Status** - Health indicators use distinct colors (green/yellow/orange/red) with glowing effects

### Fixed
- **Timer Display** - Fixed remaining and total time display
- **Manual Stop Errors** - Eliminated false error messages when user manually stops streaming
- **Batch File Reliability** - Improved Wails executable path detection in build scripts

## [1.0.0] - 2025-01-29

### Added
- Initial release
- RTMP streaming to Castr.io, YouTube, Facebook
- Quality preset selection
- Video file browser
- Stream timer with elapsed/remaining/total time
- Auto-stop when video ends option
- Multi-service support (Castr, YouTube, Facebook, Custom RTMP)

[Unreleased]: https://github.com/ohreport/restreamer/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/ohreport/restreamer/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/ohreport/restreamer/releases/tag/v1.0.0
