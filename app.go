package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx        context.Context
	streamer   *Streamer
	binaryPath string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		streamer: NewStreamer(),
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Extract embedded FFmpeg binaries
	binaryPath, err := ExtractEmbeddedBinaries()
	if err != nil {
		runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Startup Error",
			Message: fmt.Sprintf("Failed to extract FFmpeg binaries: %v", err),
		})
		return
	}
	a.binaryPath = binaryPath
	a.streamer.SetBinaryPath(binaryPath)
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	if a.streamer.IsStreaming() {
		a.streamer.Stop()
	}

	// Cleanup extracted binaries
	CleanupEmbeddedBinaries()
}

// OpenFileDialog opens a file selection dialog
func (a *App) OpenFileDialog() (string, error) {
	result, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Video File or M3U Playlist",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Video Files & Playlists",
				Pattern:     "*.mp4;*.avi;*.mov;*.mkv;*.flv;*.wmv;*.m3u;*.m3u8",
			},
			{
				DisplayName: "Video Files",
				Pattern:     "*.mp4;*.avi;*.mov;*.mkv;*.flv;*.wmv",
			},
			{
				DisplayName: "Playlists",
				Pattern:     "*.m3u;*.m3u8",
			},
			{
				DisplayName: "All Files",
				Pattern:     "*.*",
			},
		},
	})
	return result, err
}

// StartStream starts streaming a video file to Castr
func (a *App) StartStream(videoPath, rtmpURL, streamKey string) error {
	if a.streamer.IsStreaming() {
		return fmt.Errorf("stream is already running")
	}

	// Clean the RTMP URL - remove trailing slash if present
	rtmpURL = strings.TrimRight(rtmpURL, "/")

	// Validate RTMP URL format
	if !strings.HasPrefix(rtmpURL, "rtmp://") && !strings.HasPrefix(rtmpURL, "rtmps://") {
		return fmt.Errorf("invalid RTMP URL: must start with rtmp:// or rtmps://")
	}

	// Construct the full RTMP URL (including any authentication parameters)
	fullRTMPURL := fmt.Sprintf("%s/%s", rtmpURL, streamKey)

	return a.streamer.Start(videoPath, fullRTMPURL)
}

// StopStream stops the current stream
func (a *App) StopStream() error {
	return a.streamer.Stop()
}

// GetStreamStatus returns the current streaming status
func (a *App) GetStreamStatus() StreamStatus {
	return a.streamer.GetStatus()
}

// SetAutoStopOnEnd sets whether to automatically stop streaming when the video ends
func (a *App) SetAutoStopOnEnd(enabled bool) {
	a.streamer.SetAutoStopOnEnd(enabled)
}

// SetQuality sets the streaming quality preset
func (a *App) SetQuality(quality string) error {
	return a.streamer.SetQuality(quality)
}

// GetQuality returns the current streaming quality preset
func (a *App) GetQuality() string {
	return a.streamer.GetQuality()
}

// SetWindowSize sets the window to the specified dimensions
func (a *App) SetWindowSize(width, height int) {
	runtime.WindowSetSize(a.ctx, width, height)
}

// SetWindowAlwaysOnTop sets whether the window should stay on top
func (a *App) SetWindowAlwaysOnTop(onTop bool) {
	runtime.WindowSetAlwaysOnTop(a.ctx, onTop)
}

// StreamStatus represents the current status of the stream
type StreamStatus struct {
	IsStreaming      bool    `json:"isStreaming"`
	VideoPath        string  `json:"videoPath"`
	RTMPURL          string  `json:"rtmpUrl"`
	Error            string  `json:"error"`
	ElapsedSeconds   float64 `json:"elapsedSeconds"`
	DurationSeconds  float64 `json:"durationSeconds"`
	RemainingSeconds float64 `json:"remainingSeconds"`
	Quality          string  `json:"quality"`
	ConnectionHealth string  `json:"connectionHealth"`
	RetryCount       int     `json:"retryCount"`
	MaxRetries       int     `json:"maxRetries"`
}
