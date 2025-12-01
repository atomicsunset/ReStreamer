package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// QualityPreset defines a streaming quality preset
type QualityPreset struct {
	Name        string
	MaxBitrate  string // e.g., "3000k"
	BufferSize  string // e.g., "6000k"
	Resolution  string // e.g., "1920x1080", "1280x720", "960x540"
	Description string
}

// Quality presets
var QualityPresets = map[string]QualityPreset{
	"low": {
		Name:        "Low",
		MaxBitrate:  "1000k",
		BufferSize:  "2000k",
		Resolution:  "960x540",
		Description: "540p, for slow connections",
	},
	"medium": {
		Name:        "Medium",
		MaxBitrate:  "2000k",
		BufferSize:  "4000k",
		Resolution:  "1280x720",
		Description: "720p, balanced",
	},
	"high": {
		Name:        "High",
		MaxBitrate:  "3000k",
		BufferSize:  "6000k",
		Resolution:  "1920x1080",
		Description: "1080p, current default",
	},
	"ultra": {
		Name:        "Ultra",
		MaxBitrate:  "5000k",
		BufferSize:  "10000k",
		Resolution:  "1920x1080",
		Description: "1080p high quality",
	},
}

// Streamer handles video streaming to RTMP
type Streamer struct {
	mu               sync.Mutex
	cmd              *exec.Cmd
	isStreaming      bool
	videoPath        string
	rtmpURL          string
	lastError        string
	errorOutput      string
	startTime        time.Time
	duration         float64 // Video duration in seconds
	autoStopOnEnd    bool
	binaryPath       string  // Path to extracted binaries
	quality          string  // Quality preset: "low", "medium", "high", "ultra"
	connectionHealth string  // "healthy", "degraded", "disconnected", "reconnecting"
	retryCount       int     // Current retry attempt
	maxRetries       int     // Maximum reconnection attempts
	manualStop       bool    // True if user manually stopped
	lastHealthCheck  time.Time
}

// NewStreamer creates a new Streamer instance
func NewStreamer() *Streamer {
	return &Streamer{
		quality:          "high", // Default to high quality
		maxRetries:       3,      // Auto-reconnect up to 3 times
		connectionHealth: "disconnected",
	}
}

// SetQuality sets the quality preset for streaming
func (s *Streamer) SetQuality(quality string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate quality preset
	if _, ok := QualityPresets[quality]; !ok {
		return fmt.Errorf("invalid quality preset: %s", quality)
	}

	s.quality = quality
	return nil
}

// GetQuality returns the current quality preset
func (s *Streamer) GetQuality() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.quality
}

// SetBinaryPath sets the path to the extracted FFmpeg binaries
func (s *Streamer) SetBinaryPath(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.binaryPath = path
}

// getVideoDuration uses ffprobe to get the duration of a video file
func (s *Streamer) getVideoDuration(videoPath string) (float64, error) {
	// Use extracted binary path if available
	var ffprobePath string
	if s.binaryPath != "" {
		ffprobePath = filepath.Join(s.binaryPath, "ffprobe.exe")
	} else {
		// Fallback to executable directory
		exePath, err := os.Executable()
		if err != nil {
			return 0, err
		}
		exeDir := filepath.Dir(exePath)
		ffprobePath = filepath.Join(exeDir, "ffprobe")
		if runtime.GOOS == "windows" {
			ffprobePath += ".exe"
		}

		// Check if ffprobe exists
		if _, err := os.Stat(ffprobePath); os.IsNotExist(err) {
			// Try to find ffprobe in PATH
			ffprobePath, err = exec.LookPath("ffprobe")
			if err != nil {
				return 0, fmt.Errorf("ffprobe not found")
			}
		}
	}

	// Run ffprobe to get duration
	cmd := exec.Command(ffprobePath,
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoPath,
	)

	// Hide console window on Windows
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x08000000, // CREATE_NO_WINDOW
		}
	}

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get video duration: %v", err)
	}

	durationStr := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %v", err)
	}

	return duration, nil
}

// Start begins streaming a video file to the specified RTMP URL
func (s *Streamer) Start(videoPath, rtmpURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isStreaming {
		return fmt.Errorf("stream is already running")
	}

	// Check if video file exists
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		return fmt.Errorf("video file does not exist: %s", videoPath)
	}

	// Get video duration
	duration, err := s.getVideoDuration(videoPath)
	if err != nil {
		// Log warning but don't fail - duration is optional
		fmt.Printf("Warning: Could not get video duration: %v\n", err)
		duration = 0
	}

	// Reset retry state for new stream
	s.retryCount = 0
	s.manualStop = false

	return s.startStream(videoPath, rtmpURL, duration)
}

// startStream is an internal method that handles the actual FFmpeg execution
func (s *Streamer) startStream(videoPath, rtmpURL string, duration float64) error {
	// Get FFmpeg path
	ffmpegPath, err := s.getFFmpegPath()
	if err != nil {
		return err
	}

	// Get quality preset
	preset, ok := QualityPresets[s.quality]
	if !ok {
		// Fallback to high if quality is invalid
		preset = QualityPresets["high"]
	}

	// Build FFmpeg command for streaming to RTMP
	// Note: FFmpeg's internal RTMP handler treats everything after the app name as the playpath
	// So rtmp://server/app/key?password=xxx sends "key?password=xxx" as the stream key to the server
	args := []string{
		"-re", // Read input at native frame rate
		"-i", videoPath, // Input file
		"-c:v", "libx264", // Video codec
		"-preset", "veryfast", // Encoding preset (faster encoding, larger file)
		"-maxrate", preset.MaxBitrate, // Maximum bitrate from quality preset
		"-bufsize", preset.BufferSize, // Buffer size from quality preset
		"-s", preset.Resolution, // Output resolution from quality preset
		"-pix_fmt", "yuv420p", // Pixel format
		"-g", "50", // GOP size
		"-c:a", "aac", // Audio codec
		"-b:a", "128k", // Audio bitrate
		"-ar", "44100", // Audio sample rate
		"-f", "flv", // Output format (Flash Video for RTMP)
		"-flvflags", "no_duration_filesize", // Improve compatibility
		// Important: Quote the URL to prevent shell interpretation of special characters
		rtmpURL,
	}

	// Log the command for debugging
	fmt.Printf("Starting FFmpeg with command: %s %v\n", ffmpegPath, args)

	s.cmd = exec.Command(ffmpegPath, args...)

	// Hide console window on Windows
	if runtime.GOOS == "windows" {
		s.cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x08000000, // CREATE_NO_WINDOW
		}
	}

	// Capture stderr to get error messages
	stderrPipe, err := s.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Start the command
	if err := s.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start FFmpeg: %v", err)
	}

	s.isStreaming = true
	s.videoPath = videoPath
	s.rtmpURL = rtmpURL
	s.lastError = ""
	s.errorOutput = ""
	s.startTime = time.Now()
	s.duration = duration
	s.connectionHealth = "healthy"
	s.lastHealthCheck = time.Now()

	// Read stderr in a goroutine
	go s.monitorFFmpegOutput(stderrPipe)

	// Monitor the process in a goroutine
	go s.monitorStreamProcess()

	// Monitor stream health in a goroutine
	go s.monitorStreamHealth()

	return nil
}

// monitorFFmpegOutput reads FFmpeg stderr for errors and health indicators
func (s *Streamer) monitorFFmpegOutput(stderrPipe io.ReadCloser) {
	scanner := bufio.NewScanner(stderrPipe)
	var errorLines []string
	lastFrameTime := time.Now()

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line) // Print to console for debugging
		errorLines = append(errorLines, line)

		// Keep only the last 20 lines to avoid memory issues
		if len(errorLines) > 20 {
			errorLines = errorLines[len(errorLines)-20:]
		}

		// Check for frame output (indicates healthy streaming)
		if strings.Contains(line, "frame=") {
			lastFrameTime = time.Now()
			s.mu.Lock()
			if s.connectionHealth != "healthy" {
				s.connectionHealth = "healthy"
				s.retryCount = 0 // Reset retry count on recovery
			}
			s.lastHealthCheck = lastFrameTime
			s.mu.Unlock()
		}

		// Check for connection errors
		if strings.Contains(line, "Connection refused") ||
		   strings.Contains(line, "Connection timed out") ||
		   strings.Contains(line, "Failed to update") ||
		   strings.Contains(line, "I/O error") {
			s.mu.Lock()
			s.connectionHealth = "degraded"
			s.mu.Unlock()
		}
	}

	s.mu.Lock()
	// Store last few lines of output
	if len(errorLines) > 0 {
		s.errorOutput = ""
		for _, line := range errorLines {
			s.errorOutput += line + "\n"
		}
	}
	s.mu.Unlock()
}

// monitorStreamHealth checks connection health periodically
func (s *Streamer) monitorStreamHealth() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		if !s.isStreaming {
			s.mu.Unlock()
			return
		}

		// Check if we've received frames recently
		timeSinceLastHealth := time.Since(s.lastHealthCheck)
		if timeSinceLastHealth > 30*time.Second {
			if s.connectionHealth == "healthy" {
				s.connectionHealth = "degraded"
				fmt.Println("Warning: No frame updates in 30 seconds, connection degraded")
			}
		}
		s.mu.Unlock()
	}
}

// monitorStreamProcess waits for FFmpeg to exit and handles reconnection
func (s *Streamer) monitorStreamProcess() {
	err := s.cmd.Wait()
	s.mu.Lock()
	defer s.mu.Unlock()

	s.isStreaming = false
	s.connectionHealth = "disconnected"

	if err != nil {
		// Only show error if not manually stopped
		if !s.manualStop {
			if s.errorOutput != "" {
				s.lastError = fmt.Sprintf("FFmpeg error: %v\n\nLast output:\n%s", err, s.errorOutput)
			} else {
				s.lastError = fmt.Sprintf("FFmpeg error: %v", err)
			}
		}

		// Attempt auto-reconnect if not manually stopped
		if !s.manualStop && s.retryCount < s.maxRetries {
			s.retryCount++
			fmt.Printf("Stream disconnected. Attempting reconnect %d/%d in 5 seconds...\n", s.retryCount, s.maxRetries)
			s.connectionHealth = "reconnecting"

			// Unlock before sleeping to avoid deadlock
			s.mu.Unlock()
			time.Sleep(5 * time.Second)
			s.mu.Lock()

			// Attempt to restart stream
			videoPath := s.videoPath
			rtmpURL := s.rtmpURL
			duration := s.duration

			s.mu.Unlock()
			reconnectErr := s.startStream(videoPath, rtmpURL, duration)
			s.mu.Lock()

			if reconnectErr != nil {
				s.lastError = fmt.Sprintf("Reconnect attempt %d failed: %v", s.retryCount, reconnectErr)
				fmt.Printf("Reconnect failed: %v\n", reconnectErr)
			} else {
				fmt.Printf("Reconnect attempt %d successful!\n", s.retryCount)
			}
		} else if s.retryCount >= s.maxRetries {
			fmt.Printf("Maximum reconnection attempts (%d) reached. Stream stopped.\n", s.maxRetries)
			s.lastError = fmt.Sprintf("Stream failed after %d reconnection attempts", s.maxRetries)
		}
	}
}

// Stop stops the current stream
func (s *Streamer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isStreaming || s.cmd == nil {
		return fmt.Errorf("no stream is running")
	}

	// Mark as manual stop to prevent auto-reconnect
	s.manualStop = true

	// Send interrupt signal to FFmpeg for graceful shutdown
	if s.cmd.Process != nil {
		if err := s.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to stop stream: %v", err)
		}
	}

	s.isStreaming = false
	return nil
}

// IsStreaming returns whether a stream is currently active
func (s *Streamer) IsStreaming() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isStreaming
}

// SetAutoStopOnEnd sets whether to automatically stop when video ends
func (s *Streamer) SetAutoStopOnEnd(enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.autoStopOnEnd = enabled
}

// GetStatus returns the current stream status
func (s *Streamer) GetStatus() StreamStatus {
	s.mu.Lock()
	defer s.mu.Unlock()

	var elapsed, remaining float64
	if s.isStreaming && !s.startTime.IsZero() {
		elapsed = time.Since(s.startTime).Seconds()
		if s.duration > 0 {
			remaining = s.duration - elapsed
			if remaining < 0 {
				remaining = 0
			}
		}
	}

	return StreamStatus{
		IsStreaming:      s.isStreaming,
		VideoPath:        s.videoPath,
		RTMPURL:          s.rtmpURL,
		Error:            s.lastError,
		ElapsedSeconds:   elapsed,
		DurationSeconds:  s.duration,
		RemainingSeconds: remaining,
		Quality:          s.quality,
		ConnectionHealth: s.connectionHealth,
		RetryCount:       s.retryCount,
		MaxRetries:       s.maxRetries,
	}
}

// getFFmpegPath returns the path to the FFmpeg executable
func (s *Streamer) getFFmpegPath() (string, error) {
	// First, check if we have extracted binaries
	if s.binaryPath != "" {
		ffmpegPath := filepath.Join(s.binaryPath, "ffmpeg.exe")
		if _, err := os.Stat(ffmpegPath); err == nil {
			return ffmpegPath, nil
		}
	}

	// Fallback: check if ffmpeg is in the same directory as the executable
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		localFFmpeg := filepath.Join(exeDir, "ffmpeg")
		if runtime.GOOS == "windows" {
			localFFmpeg += ".exe"
		}
		if _, err := os.Stat(localFFmpeg); err == nil {
			return localFFmpeg, nil
		}
	}

	// Check if ffmpeg is in PATH
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", fmt.Errorf("FFmpeg not found. Please install FFmpeg or place ffmpeg.exe in the application directory")
	}

	return ffmpegPath, nil
}
