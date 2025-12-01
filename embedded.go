package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed ffmpeg.exe
var ffmpegBinary []byte

//go:embed ffprobe.exe
var ffprobeBinary []byte

var extractedDir string

// ExtractEmbeddedBinaries extracts ffmpeg and ffprobe to a temporary directory
func ExtractEmbeddedBinaries() (string, error) {
	// Create a temp directory in the user's temp folder
	tempDir := filepath.Join(os.TempDir(), "OH.Report-ReStreamer")

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}

	// Extract ffmpeg.exe
	ffmpegPath := filepath.Join(tempDir, "ffmpeg.exe")
	if err := os.WriteFile(ffmpegPath, ffmpegBinary, 0755); err != nil {
		return "", fmt.Errorf("failed to extract ffmpeg: %v", err)
	}

	// Extract ffprobe.exe
	ffprobePath := filepath.Join(tempDir, "ffprobe.exe")
	if err := os.WriteFile(ffprobePath, ffprobeBinary, 0755); err != nil {
		return "", fmt.Errorf("failed to extract ffprobe: %v", err)
	}

	extractedDir = tempDir
	return tempDir, nil
}

// CleanupEmbeddedBinaries removes the temporary directory
func CleanupEmbeddedBinaries() {
	if extractedDir != "" {
		os.RemoveAll(extractedDir)
	}
}
