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
	"strings"
	"time"
)

func main() {

	path := ""

	if runtime.GOOS == "darwin" {
		path = "/Applications/Brave Browser.app/Contents/MacOS/Brave Browser"
	}

	//runBrowser(path)

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
	runBrowser(path)

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

func runBrowser(browserPath string) {
	now := time.Now()
	fmt.Println("start --> ", now)
	url := launcher.New().
		Bin(browserPath). // use Chrome instead of default Chromium
		Headless(false).  // show browser
		MustLaunch()

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("")       // open blank first
	page.MustEmulate(devices.GalaxyS5) // emulate full mobile device

	//page.MustNavigate("https://satumomen.com/preview/peresmian-rs")
	page.MustNavigate("https://joinedwithshan.viding.co/")
	//page.MustNavigate("https://app.sangmempelai.id/pilihan-tema/sunda-01")
	//page.MustNavigate("https://adirara.webnikah.com/?templatecoba=156/kepada:Budi%20dan%20Ani-Bandung")
	//page.MustNavigate("https://ourmoment.my.id/art-6/")
	page.MustWaitLoad()

	// Find all possible clickable elements
	elements := page.MustElements("a, button, div, span")

	var target *rod.Element
	for _, el := range elements {

		text, err := el.Text()
		if err != nil {
			continue
		}

		if strings.EqualFold(strings.TrimSpace(text), "Buka Undangan") {
			target = el
			break
		}
		if strings.EqualFold(strings.TrimSpace(text), "Open The Invitation") {
			target = el
			break
		}
	}

	if target == nil {
		fmt.Println("Element with text 'Buka Undangan' not found")
		time.Sleep(1 * time.Second) // wait before exit
		return
	}

	time.Sleep(2 * time.Second) // First page wait

	target.MustClick()
	//target.MustScrollIntoView()

	videoExists := page.MustEval(`() => !!document.querySelector("video.elementor-background-video-hosted")`).Bool()

	if videoExists {
		time.Sleep(13 * time.Second) // After Click button delay
	} else {
		time.Sleep(5 * time.Second) // After Click button delay
	}

	scrollInterval := 2 * time.Second // configurable scroll interval
	scrollToBottom(page, scrollInterval, 450)

	time.Sleep(1 * time.Second) // wait before exit
	fmt.Println("end --> ", time.Now().Sub(now).Minutes())
}

func scrollToBottom(page *rod.Page, interval time.Duration, step int) {
	for {
		// Get current scroll position and total height
		pos := page.MustEval(`() => window.scrollY`).Int()
		height := page.MustEval(`() => document.body.scrollHeight`).Int()

		// Stop if we are at or near the bottom
		if pos+step >= height {
			page.MustEval(`() => window.scrollTo(0, document.body.scrollHeight)`)
			fmt.Println("Reached the bottom.")
			break
		}

		// Scroll by step
		script := fmt.Sprintf("() => window.scrollTo(0, %d)", pos+step)
		page.MustEval(script)

		// Wait for new content to possibly load
		time.Sleep(interval)

		if isAtBottom(page) {
			fmt.Println("Reached the bottom.")
			break
		}
	}
}

func isAtBottom(page *rod.Page) bool {
	result := page.MustEval(`() => {
		return (window.innerHeight + window.scrollY) >= (document.body.scrollHeight - 100);
	}`)
	return result.Bool()
}
