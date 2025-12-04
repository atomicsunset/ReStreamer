// DOM elements
const videoPathInput = document.getElementById('videoPath');
const browseBtn = document.getElementById('browseBtn');
const serviceSelect = document.getElementById('serviceSelect');
const qualitySelect = document.getElementById('qualitySelect');
const rtmpUrlGroup = document.getElementById('rtmpUrlGroup');
const streamKeyGroup = document.getElementById('streamKeyGroup');
const rtmpUrlInput = document.getElementById('rtmpUrl');
const streamKeyInput = document.getElementById('streamKey');
const startBtn = document.getElementById('startBtn');
const stopBtn = document.getElementById('stopBtn');
const statusIndicator = document.getElementById('statusIndicator');
const statusText = document.getElementById('statusText');
const errorMessage = document.getElementById('errorMessage');
const timerDisplay = document.getElementById('timerDisplay');
const elapsedTime = document.getElementById('elapsedTime');
const remainingTime = document.getElementById('remainingTime');
const totalTime = document.getElementById('totalTime');
const autoStopCheckbox = document.getElementById('autoStopCheckbox');
const compactModeBtn = document.getElementById('compactModeBtn');
const streamEndedPopup = document.getElementById('streamEndedPopup');
const dismissPopupBtn = document.getElementById('dismissPopupBtn');

// Castr channel keys mapping
// NOTE: This is for demo purposes only. For production use, store your keys securely
// and do NOT commit them to version control!
const castrChannelKeys = {
    // Add your own Castr channel keys here
    // Example format:
    // '000': 'live_xxxxxxxxxxxxxxxxxxxxxxxx?password=xxxxxxxx',
};

// State
let isStreaming = false;
let wasStreaming = false;  // Track previous streaming state to detect stream end
let isCompactMode = false;

// Window sizes
const FULL_WIDTH = 800;
const FULL_HEIGHT = 600;
const COMPACT_WIDTH = 400;
const COMPACT_HEIGHT = 280;

// Helper function to format seconds as HH:MM:SS
function formatTime(seconds) {
    if (!seconds || seconds < 0) return '00:00:00';

    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = Math.floor(seconds % 60);

    return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
}

// Service dropdown handler
serviceSelect.addEventListener('change', () => {
    // Clear stream key when changing services
    streamKeyInput.value = '';

    // Set RTMP URL based on service selection
    if (serviceSelect.value) {
        rtmpUrlInput.value = serviceSelect.value;
    }
});

// Quality dropdown handler
qualitySelect.addEventListener('change', async () => {
    try {
        await window.go.main.App.SetQuality(qualitySelect.value);
    } catch (err) {
        console.error('Failed to set quality preference:', err);
        showError('Failed to set quality: ' + err);
    }
});

// Auto-stop checkbox handler
autoStopCheckbox.addEventListener('change', async () => {
    try {
        await window.go.main.App.SetAutoStopOnEnd(autoStopCheckbox.checked);
    } catch (err) {
        console.error('Failed to set auto-stop preference:', err);
    }
});

// File selection using Wails runtime
browseBtn.addEventListener('click', async () => {
    try {
        const result = await window.go.main.App.OpenFileDialog();

        if (result) {
            videoPathInput.value = result;
        }
    } catch (err) {
        showError('Failed to open file dialog: ' + err);
    }
});

// Start streaming
startBtn.addEventListener('click', async () => {
    const videoPath = videoPathInput.value.trim();
    const rtmpUrl = rtmpUrlInput.value.trim();
    const streamKey = streamKeyInput.value.trim();

    // Validate inputs
    if (!videoPath) {
        showError('Please select a video file');
        return;
    }

    if (!rtmpUrl) {
        showError('Please enter an RTMP URL');
        return;
    }

    if (!streamKey) {
        showError('Please enter a stream key');
        return;
    }

    try {
        hideError();
        setStreaming(true);

        await window.go.main.App.StartStream(videoPath, rtmpUrl, streamKey);

        // Start status polling
        startStatusPolling();
    } catch (err) {
        setStreaming(false);
        showError('Failed to start stream: ' + err);
    }
});

// Stop streaming
stopBtn.addEventListener('click', async () => {
    try {
        // Get current status to check if video is still playing
        const status = await window.go.main.App.GetStreamStatus();

        // If video hasn't finished, ask for confirmation
        if (status.remainingSeconds > 5) {
            const remaining = formatTime(status.remainingSeconds);
            const confirmed = confirm(`Stream still has ${remaining} remaining. Are you sure you want to stop?`);
            if (!confirmed) {
                return; // User cancelled
            }
        }

        await window.go.main.App.StopStream();
        setStreaming(false);
        hideError();
    } catch (err) {
        showError('Failed to stop stream: ' + err);
    }
});

// Status polling
let statusInterval = null;

function startStatusPolling() {
    if (statusInterval) {
        clearInterval(statusInterval);
    }

    statusInterval = setInterval(async () => {
        try {
            const status = await window.go.main.App.GetStreamStatus();

            if (status.isStreaming) {
                setStreaming(true, status.connectionHealth, status.retryCount, status.maxRetries);

                // Update timer display
                if (status.durationSeconds > 0) {
                    elapsedTime.textContent = formatTime(status.elapsedSeconds);
                    remainingTime.textContent = formatTime(status.remainingSeconds);
                    totalTime.textContent = formatTime(status.durationSeconds);
                    timerDisplay.style.display = 'flex';
                } else {
                    // No duration info, just show elapsed
                    elapsedTime.textContent = formatTime(status.elapsedSeconds);
                    remainingTime.textContent = '--:--:--';
                    totalTime.textContent = '--:--:--';
                    timerDisplay.style.display = 'flex';
                }

                if (status.error) {
                    showError(status.error);
                }
            } else {
                // Detect stream completion (was streaming, now stopped)
                if (wasStreaming && !status.error) {
                    showStreamEndedPopup();
                }

                setStreaming(false);
                timerDisplay.style.display = 'none';

                if (status.error) {
                    showError(status.error);
                }
                if (statusInterval) {
                    clearInterval(statusInterval);
                    statusInterval = null;
                }
            }
        } catch (err) {
            console.error('Status check failed:', err);
        }
    }, 1000);
}

// UI helpers
function setStreaming(streaming, connectionHealth = 'healthy', retryCount = 0, maxRetries = 3) {
    wasStreaming = isStreaming;  // Track previous state before updating
    isStreaming = streaming;

    if (streaming) {
        // Set status class based on connection health
        let statusClass = 'status-indicator streaming';
        let statusMessage = 'Streaming...';

        switch(connectionHealth) {
            case 'healthy':
                statusClass = 'status-indicator streaming healthy';
                statusMessage = 'Streaming';
                break;
            case 'degraded':
                statusClass = 'status-indicator streaming degraded';
                statusMessage = 'Connection Degraded';
                break;
            case 'reconnecting':
                statusClass = 'status-indicator streaming reconnecting';
                statusMessage = `Reconnecting (${retryCount}/${maxRetries})...`;
                break;
            case 'disconnected':
                statusClass = 'status-indicator streaming disconnected';
                statusMessage = 'Disconnected';
                break;
        }

        statusIndicator.className = statusClass;
        statusText.textContent = statusMessage;
        startBtn.disabled = true;
        stopBtn.disabled = false;

        // Disable inputs while streaming
        videoPathInput.disabled = true;
        browseBtn.disabled = true;
        rtmpUrlInput.disabled = true;
        streamKeyInput.disabled = true;
        qualitySelect.disabled = true;
        serviceSelect.disabled = true;
    } else {
        statusIndicator.className = 'status-indicator stopped';
        statusText.textContent = 'Ready';
        startBtn.disabled = false;
        stopBtn.disabled = true;

        // Enable inputs
        videoPathInput.disabled = false;
        browseBtn.disabled = false;
        rtmpUrlInput.disabled = false;
        streamKeyInput.disabled = false;
        qualitySelect.disabled = false;
        serviceSelect.disabled = false;
    }
}

function showError(message) {
    errorMessage.textContent = message;
    errorMessage.style.display = 'block';
}

function hideError() {
    errorMessage.style.display = 'none';
    errorMessage.textContent = '';
}

// Initialize quality setting on page load
async function initializeQuality() {
    try {
        const currentQuality = await window.go.main.App.GetQuality();
        if (currentQuality) {
            qualitySelect.value = currentQuality;
        }
    } catch (err) {
        console.error('Failed to get initial quality:', err);
    }
}

// Initialize on page load
window.addEventListener('DOMContentLoaded', () => {
    initializeQuality();
});

// Cleanup on window close
window.addEventListener('beforeunload', () => {
    if (statusInterval) {
        clearInterval(statusInterval);
    }
});

// Stream Ended Popup functions
function showStreamEndedPopup() {
    streamEndedPopup.style.display = 'flex';
    // Make window always on top when showing popup
    window.go.main.App.SetWindowAlwaysOnTop(true);
}

function hideStreamEndedPopup() {
    streamEndedPopup.style.display = 'none';
    // Restore normal window behavior
    window.go.main.App.SetWindowAlwaysOnTop(false);
}

// Dismiss popup button handler
dismissPopupBtn.addEventListener('click', () => {
    hideStreamEndedPopup();
});

// Also dismiss popup when clicking outside
streamEndedPopup.addEventListener('click', (e) => {
    if (e.target === streamEndedPopup) {
        hideStreamEndedPopup();
    }
});

// Compact Mode functions
async function toggleCompactMode() {
    isCompactMode = !isCompactMode;

    if (isCompactMode) {
        // Switch to compact mode
        document.body.classList.add('compact-mode');
        compactModeBtn.querySelector('.compact-icon').textContent = '⊞';
        compactModeBtn.title = 'Toggle Full Mode';
        await window.go.main.App.SetWindowSize(COMPACT_WIDTH, COMPACT_HEIGHT);
    } else {
        // Switch to full mode
        document.body.classList.remove('compact-mode');
        compactModeBtn.querySelector('.compact-icon').textContent = '⊟';
        compactModeBtn.title = 'Toggle Compact Mode';
        await window.go.main.App.SetWindowSize(FULL_WIDTH, FULL_HEIGHT);
    }
}

// Compact mode button handler
compactModeBtn.addEventListener('click', toggleCompactMode);
