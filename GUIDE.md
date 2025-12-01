# Atomation ReStreamer - User Guide

## What is ReStreamer?

Atomation ReStreamer is a simple tool for streaming pre-recorded videos to platforms like Facebook and YouTube without needing heavy software like vMix or OBS. Perfect for delayed broadcasts and scheduled content.

## System Requirements

- **Operating System**: Windows 10 or Windows 11
- **Storage**: ~200MB for the application
- **Network**: Stable internet connection for streaming

## Installation

**There is no installation!**

1. Download `Atomation-ReStreamer.exe` from the GitHub release
2. Copy it anywhere on your computer
3. Double-click to run

That's it. The application is completely standalone.

## Getting Started

### Interface Overview

When you open ReStreamer, you'll see a two-column layout:

**Left Side (Configuration)**
- Video file selection
- Streaming service selection
- Channel or stream key settings

**Right Side (Status & Controls)**
- Stream status indicator
- Timer (elapsed, remaining, total)
- Start/Stop buttons
- Auto-stop option

## Streaming to Castr

### Step 1: Select Service
1. Click the **Service** dropdown
2. Select **Castr.io**

### Step 2: Select Channel
1. A **Channel** dropdown will appear
2. Select your channel (Castr 000 through Castr 008)
3. The stream key is automatically filled in

### Step 3: Select Video
1. Click **Browse** next to Video File
2. Choose your video file (MP4, AVI, MOV, MKV, etc.)
3. Or select an M3U playlist for continuous streaming

### Step 4: Start Streaming
1. Click the **Start** button
2. Watch the status change to "Streaming"
3. Monitor progress with the timer

### Step 5: Stop When Done
- Click **Stop** to end the stream manually
- Or enable **Auto-stop when ended** to stop automatically

## Streaming to Facebook Live

### Step 1: Get Your Stream Key
1. Go to [facebook.com/live/producer](https://www.facebook.com/live/producer)
2. Click "Go Live"
3. Copy your **Stream Key**

### Step 2: Configure ReStreamer
1. Select **Facebook** from the Service dropdown
2. The RTMP URL is automatically set
3. Paste your **Stream Key**

### Step 3: Select Video & Start
1. Browse and select your video file
2. Click **Start** to begin streaming

## Streaming to YouTube Live

### Step 1: Get Your Stream Key
1. Go to [YouTube Studio](https://studio.youtube.com)
2. Click "Create" → "Go Live"
3. Under "Stream settings", copy your **Stream Key**
4. Set your title, description, and privacy settings on YouTube

### Step 2: Configure ReStreamer
1. Select **YouTube** from the Service dropdown
2. The RTMP URL is automatically set
3. Paste your **Stream Key**

### Step 3: Select Video & Start
1. Browse and select your video file
2. Click **Start** to begin streaming

**Important**: You still need to configure stream details (title, description, thumbnail) on YouTube's website before or during streaming.

## Streaming to Custom RTMP Server

### Step 1: Select Custom
1. Select **Custom RTMP** from the Service dropdown

### Step 2: Enter Details
1. Enter your **RTMP URL** (e.g., `rtmp://server.com/live`)
2. Enter your **Stream Key**

### Step 3: Select Video & Start
1. Browse and select your video file
2. Click **Start**

## Using M3U Playlists

### What is M3U?

M3U playlists let you stream multiple videos in sequence. Perfect for automated, multi-hour broadcasts.

### Creating an M3U Playlist

Create a text file with `.m3u` extension:

```
#EXTM3U
C:\Videos\intro.mp4
C:\Videos\main-content.mp4
C:\Videos\outro.mp4
```

Or with titles:

```
#EXTM3U
#EXTINF:-1,Introduction
C:\Videos\intro.mp4
#EXTINF:-1,Main Content
C:\Videos\main-content.mp4
#EXTINF:-1,Closing
C:\Videos\outro.mp4
```

### Using the Playlist

1. Create your M3U file
2. In ReStreamer, click **Browse**
3. Select your `.m3u` file
4. Start streaming - videos will play in order

## Understanding the Timer

During streaming, you'll see three time displays:

- **Elapsed**: How long you've been streaming
- **Remaining**: Time left in the current video
- **Total**: Total length of the video

If the timer shows `--:--:--`, the video duration couldn't be detected (rare with standard formats).

## Auto-Stop Feature

Enable **Auto-stop when ended** to automatically stop streaming when the video finishes.

**When to use it:**
- Single video streams that should end automatically
- Timed content with specific end times

**When NOT to use it:**
- M3U playlists (will stop after first video)
- Looping content
- When you want to manually control the end time

## Tips & Best Practices

### Video Format Recommendations
- **Best**: MP4 with H.264 video and AAC audio
- **Also works**: AVI, MOV, MKV, FLV, WMV
- **Resolution**: Any (1080p, 720p, 4K all supported)
- **Avoid**: Exotic codecs or broken files

### Network Considerations
- Use wired Ethernet for best stability
- Minimum upload speed: 5 Mbps for HD content
- Close bandwidth-heavy applications during streaming
- Test with a short video first

### Stream Quality
The app streams at:
- **Video**: Up to 3 Mbps (adjustable in source code)
- **Audio**: 128 kbps AAC
- **Format**: H.264 video in FLV container

This provides excellent quality for most online platforms.

## Troubleshooting

### Stream Won't Start

**Check:**
- Is the video file path correct?
- Is the RTMP URL valid?
- Is the stream key correct?
- Do you have internet connection?

**For Castr:**
- Try a different channel
- Verify channel is active in Castr dashboard

### Video Plays But Stream Doesn't Work

**Possible causes:**
- Incorrect stream key
- RTMP server is down
- Platform rejected the connection (check platform dashboard)
- Firewall blocking RTMP traffic (port 1935)

### Stream Stops Unexpectedly

**Check:**
- Internet connection stability
- Platform didn't disconnect (check platform dashboard)
- Video file isn't corrupted
- Sufficient disk space in temp folder

### Timer Shows Wrong Duration

This is usually harmless. The app uses FFprobe to detect duration, but some files don't report it correctly. Streaming still works normally.

### Poor Stream Quality

**Solutions:**
- Check your upload bandwidth
- Reduce video resolution before streaming
- Use wired connection instead of WiFi
- Close other applications using bandwidth

## Advanced Usage

### Running Multiple Instances

You can run multiple ReStreamer instances simultaneously to stream to different platforms at once:

1. Launch ReStreamer
2. Configure for Platform A and start streaming
3. Launch ReStreamer again (new window)
4. Configure for Platform B and start streaming

Each instance operates independently.

### Keyboard Shortcuts

Currently, there are no keyboard shortcuts. Use mouse/touch for all interactions.

## Known Limitations

- No built-in recording while streaming
- No stream preview (you see only the interface)
- No automatic reconnection if stream drops
- Facebook/YouTube require manual setup of title/description on platform
- One video at a time per instance (use M3U for playlists)

## Getting Help

### Before Asking for Help

1. Check this guide's Troubleshooting section
2. Verify video file plays in VLC or Windows Media Player
3. Test stream key on platform's website
4. Try with a different video file

### Technical Support

For technical issues:
- Review DEVELOPERS.md for technical details
- Check DEPLOYMENT.md for deployment information
- Contact OH.Report support at info@oh.report

## Privacy & Security

- Stream keys are stored in the app's memory only during your session
- No data is sent anywhere except your chosen streaming platform
- Video files are read directly (no copying to temp folders)
- FFmpeg binaries extracted to system temp folder on launch

## Updates

To update ReStreamer:
1. Download the new version
2. Replace the old .exe file
3. Launch the new version

Your settings don't persist between sessions (by design for security).

## Frequently Asked Questions

**Q: Can I stream to multiple platforms at once?**
A: Yes, run multiple instances of ReStreamer.

**Q: Does this work on Mac or Linux?**
A: No, Windows only. The app uses Windows-specific features.

**Q: Can I edit videos in ReStreamer?**
A: No, ReStreamer only streams existing videos. Edit videos in other software first.

**Q: Why can't I see the video preview?**
A: ReStreamer is a streaming tool, not a player. Use VLC to preview videos beforehand.

**Q: Can I schedule streams to start later?**
A: Not currently. You must manually start the stream.

**Q: Does this replace vMix or OBS?**
A: Only for simple file streaming. vMix and OBS offer much more (live mixing, graphics, etc.).

**Q: Is this free?**
A: Check with OH.Report for licensing information.

**Q: Can I add my own Castr channels?**
A: Not in the UI, but developers can edit the source code (see DEVELOPERS.md).

---

**Version**: 1.0.0
**Last Updated**: 2025-01-30

For technical documentation, see **DEVELOPERS.md**
