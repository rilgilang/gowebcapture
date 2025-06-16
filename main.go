package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func main() {

	dir, _ := os.Getwd()
	outputPath := filepath.Join(dir, "output.mp4")

	var input string
	var format string

	if runtime.GOOS == "darwin" {
		// Check with: ffmpeg -f avfoundation -list_devices true -i ""
		input = "1:none" // or "0:none" depending on screen index
		format = "avfoundation"
	} else {
		display := os.Getenv("DISPLAY")
		if display == "" {
			display = ":99.0"
		}
		os.Setenv("DISPLAY", display)
		input = display
		format = "x11grab"
	}

	cmd, err := startFFmpeg(input, format, outputPath)
	if err != nil {
		fmt.Println("Failed to start ffmpeg:", err)
		return
	}

	// Launch browser and interact
	runBrowser()

	// Wait a little after browser interaction
	time.Sleep(1 * time.Second)

	// Stop ffmpeg cleanly
	fmt.Println("Stopping ffmpeg...")
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		fmt.Println("Failed to stop ffmpeg gracefully:", err)
	}
	_ = cmd.Wait()

	fmt.Println("Recording complete. Saved to:", outputPath)
}

func startFFmpeg(input, format, outputPath string) (*exec.Cmd, error) {
	stream := ffmpeg_go.Input(input, ffmpeg_go.KwArgs{
		"f":          format,
		"framerate":  "30",
		"video_size": "375x812",
	}).Output(outputPath, ffmpeg_go.KwArgs{
		"c:v": "libx264",
		"y":   "",
	}).OverWriteOutput()

	cmd := stream.Compile()
	proc := exec.Command(cmd.Path, cmd.Args[1:]...)
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr

	if err := proc.Start(); err != nil {
		return nil, err
	}

	return proc, nil
}

func runBrowser() {
	url := launcher.New().
		Headless(false). // show browser
		MustLaunch()

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("")                        // open blank first
	page.MustEmulate(devices.IPhoneX)                   // emulate full mobile device
	page.MustNavigate("https://ourmoment.my.id/art-1/") // load page
	page.MustWaitLoad()

	time.Sleep(5 * time.Second)
}

func testBrowser() {
	url := launcher.New().
		Headless(false). // show browser
		MustLaunch()

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("")                        // open blank first
	page.MustEmulate(devices.IPhone5orSE)               // emulate full mobile device
	page.MustNavigate("https://ourmoment.my.id/art-1/") // load page
	page.MustWaitLoad()
}
