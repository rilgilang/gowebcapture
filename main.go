package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"log"
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

	now := time.Now()
	// Launch browser and interact
	err := runBrowser(path, now)

	if err != nil {
		fmt.Println("err --> ", err)
	}

	fmt.Println("Recording complete.")
}

func startFFmpeg(currentDateTime *time.Time) (*exec.Cmd, error) {

	dir, _ := os.Getwd()
	outputPath := filepath.Join(dir, fmt.Sprintf(`%s.mp4`, currentDateTime.Format("2006-01-02-15-04-05")))

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

	stream := ffmpeg_go.Input(input, ffmpeg_go.KwArgs{
		"f":          format,
		"framerate":  "30",
		"video_size": "390x844",
	}).Filter("crop", ffmpeg_go.Args{"720:1280:0:210"}).Output(outputPath, ffmpeg_go.KwArgs{
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

func stopFFmpeg(cmd *exec.Cmd) error {
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		fmt.Println("Failed to interrupt ffmpeg:", err)
		if killErr := cmd.Process.Kill(); killErr != nil {
			fmt.Println("Failed to kill ffmpeg:", killErr)
		}
		return err
	}
	_ = cmd.Wait()

	return nil
}

func runBrowser(browserPath string, now time.Time) error {
	url := launcher.New()

	if browserPath != "" {
		url.Bin(browserPath) // use Chrome instead of default Chromium
	}

	controlUrl := url.Headless(false).MustLaunch()

	browser := rod.New().ControlURL(controlUrl).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("")       // open blank first
	page.MustEmulate(devices.GalaxyS5) // emulate full mobile device

	page.MustSetWindow(0, 0, 360, 1024)

	//page.MustNavigate("https://satumomen.com/preview/peresmian-rs")
	//page.MustNavigate("https://joinedwithshan.viding.co/")
	//page.MustNavigate("https://app.sangmempelai.id/pilihan-tema/sunda-01")
	//page.MustNavigate("https://adirara.webnikah.com/?templatecoba=156/kepada:Budi%20dan%20Ani-Bandung")
	//page.MustNavigate("https://ourmoment.my.id/art-6/")
	page.MustNavigate("http://127.0.0.1:8000")

	page.MustWaitLoad()

	// Start ffmpeg

	cmd, err := startFFmpeg(&now)
	if err != nil {
		fmt.Println("Failed to start ffmpeg:", err)
		return err
	}

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

		err = stopFFmpeg(cmd)

		e := os.Remove(fmt.Sprintf(`%s.mp4`, now.Format("2006-01-02-15-04-05")))
		if e != nil {
			log.Fatal(e)
		}

		return err
	}

	time.Sleep(5 * time.Second) // First page wait

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

	// Wait a little after browser interaction
	//time.Sleep(5 * time.Second)

	//Stop ffmpeg
	err = stopFFmpeg(cmd)
	if err != nil {
		return err
	}

	//time.Sleep(1 * time.Second) // wait before exit
	fmt.Println("end --> ", time.Now().Sub(now).Minutes())
	return nil
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
