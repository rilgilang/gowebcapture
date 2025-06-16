package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	dir, _ := os.Getwd()
	outputPath := filepath.Join(dir, "output.mp4")

	// Ensure DISPLAY is set
	display := os.Getenv("DISPLAY")
	if display == "" {
		display = ":99"
	}
	os.Setenv("DISPLAY", display) // Ensure Chromium and ffmpeg both use the same

	// Start ffmpeg recording
	ffmpeg := exec.Command("ffmpeg",
		"-y",
		"-video_size", "1366x768",
		"-framerate", "30",
		"-f", "x11grab",
		"-i", display,
		"-c:v", "libx264",
		outputPath,
	)

	ffmpeg.Stdout = os.Stdout
	ffmpeg.Stderr = os.Stderr

	if err := ffmpeg.Start(); err != nil {
		log.Fatalf("Failed to start ffmpeg: %v", err)
	}
	log.Println("Recording started...")

	// Launch browser with Rod
	url := launcher.New().
		Headless(false).
		MustLaunch()

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://example.com")
	page.MustWaitLoad()

	// Wait for it to render
	time.Sleep(5 * time.Second)

	// Stop ffmpeg
	log.Println("Stopping recording...")
	err := ffmpeg.Process.Signal(os.Interrupt)
	if err != nil {
		log.Printf("Error stopping ffmpeg: %v", err)
	} else {
		log.Println("Video should be saved as", outputPath)
	}
}
