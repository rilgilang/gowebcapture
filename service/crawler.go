package service

import (
	"context"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"
	"go/src/github.com/rilgilang/gowebcapture/bootstrap"
	"go/src/github.com/rilgilang/gowebcapture/pkg"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Crawler interface {
	RunBrowserAndInteract(ctx context.Context, urlLink string) error
}

type crawler struct {
	storage pkg.Storage
	config  *bootstrap.Config
}

func NewCrawler(storage pkg.Storage, config *bootstrap.Config) Crawler {
	return &crawler{
		storage: storage,
		config:  config,
	}
}

func (c *crawler) RunBrowserAndInteract(ctx context.Context, urlLink string) error {

	path := ""
	dir, _ := os.Getwd() // Used for output file

	switch runtime.GOOS {
	case "darwin":
		// macOS with Brave Browser
		path = c.config.DarwinBrowserPath
	case "linux":
		// Linux inside Docker with Google Chrome
		path = c.config.LinuxBrowserPath
	default:
		fmt.Println("Unsupported OS:", runtime.GOOS)
		os.Exit(1)
	}

	now := time.Now()

	url := launcher.New().
		Headless(true).
		Delete("disable-gpu").
		Set("hide-scrollbars"). // Hide scrollbars
		Leakless(true)

	if path != "" {
		url.Bin(path)
		url.RemoteDebuggingPort(3000)
		// url.Set("display", ":99")
	}

	controlUrl := url.Headless(false).MustLaunch()

	browser := rod.New().ControlURL(controlUrl).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("")      // open blank first
	page.MustEmulate(devices.Nexus6P) // emulate full mobile device
	page.MustSetWindow(0, 0, 450, 850)

	page.MustNavigate(urlLink)

	page.MustWaitLoad()

	time.Sleep(2 * time.Second) // wait before exit

	// Start ffmpeg
	cmd, err := StartFFmpeg(&now, c.config, dir)
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

		err = StopFFmpeg(cmd)

		e := os.Remove(fmt.Sprintf(`%s.mp4`, now.Format("2006-01-02-15-04-05")))
		if e != nil {
			fmt.Println("err --> ", err)
		}

		return err
	}

	time.Sleep(2 * time.Second) // First page wait

	target.MustClick()

	videoExists := page.MustEval(`() => !!document.querySelector("video.elementor-background-video-hosted")`).Bool()

	if videoExists {
		time.Sleep(13 * time.Second) // After Click button delay
	} else {
		time.Sleep(1 * time.Second) // After Click button delay
	}

	//TODO make throttle is configurable
	scrollToBottomSmoothWithThrottle(page, 3, 15) // 30px max step, ~60fps (16ms delay)

	//Stop ffmpeg
	err = StopFFmpeg(cmd)
	if err != nil {
		return err
	}

	// Put output video to storage
	if err = c.putOutputToStorage(ctx, &now, dir); err != nil {
		return err
	}

	// Remove output when done
	if err = deleteUnusedOutput(&now); err != nil {
		return err
	}

	return nil
}

func scrollToBottomSmoothWithThrottle(page *rod.Page, maxStep int, delayMs int) {
	script := fmt.Sprintf(`() => {
		if (window.__scrolling) return;
		window.__scrolling = true;

		const delay = %d;
		const maxStep = %d;

		let lastScrollTime = Date.now();

		function step() {
			const now = Date.now();
			const deltaTime = now - lastScrollTime;

			const currentY = window.scrollY;
			const viewHeight = window.innerHeight;
			const scrollHeight = document.body.scrollHeight;
			const distanceToBottom = scrollHeight - (currentY + viewHeight);

			if (distanceToBottom <= 10) {
				window.__scrolling = false;
				return;
			}

			// Easing: slow down as we approach bottom
			let speedFactor = distanceToBottom / scrollHeight; // near bottom => 0
			if (speedFactor < 0.1) speedFactor = 0.1;

			const scrollAmount = Math.min(maxStep, distanceToBottom * speedFactor);
			window.scrollBy({ top: scrollAmount, behavior: 'smooth' });

			lastScrollTime = now;
			setTimeout(step, delay);
		}

		step();
	}`, delayMs, maxStep)

	page.MustEval(script)
	waitUntilScrollStops(page)
}

func waitUntilScrollStops(page *rod.Page) {
	startTime := time.Now()
	lastPosition := 0
	stableCount := 0

	for {
		currentPos, err := page.Eval(`() => window.scrollY`)
		if err != nil {
			fmt.Println("Scroll position check failed:", err)
			break
		}

		pos := currentPos.Value.Int()

		// Check if we've reached bottom
		if err == nil && isAtBottom(page) {
			time.Sleep(2 * time.Second) // Final wait at bottom
			break
		}

		// Check if scrolling has stopped
		if pos == lastPosition {
			stableCount++
			if stableCount > 3 { // 3 consecutive checks with no movement
				fmt.Println("Scroll appears to have stopped")
				break
			}
		} else {
			stableCount = 0
			lastPosition = pos
		}

		// Timeout after 5 minutes
		if time.Since(startTime) > 5*time.Minute {
			fmt.Println("Scroll timeout reached")
			break
		}

		time.Sleep(300 * time.Millisecond)
	}

	// Clean up
	page.MustEval(`() => { window.__scrolling = false; }`)
}

func isAtBottom(page *rod.Page) bool {
	result := page.MustEval(`() => {
		return (window.innerHeight + window.scrollY) >= (document.body.scrollHeight - 100);
	}`)
	return result.Bool()
}

// TODO move this function to recorder.go
func (c *crawler) putOutputToStorage(ctx context.Context, currentDateTime *time.Time, dir string) error {

	filePath := filepath.Join(dir, fmt.Sprintf(`/output/%s.mp4`, currentDateTime.Format("2006-01-02-15-04-05")))

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("err open file: ", err)
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("Error reading file stats: ", err)
		return err
	}
	stat.Size()

	defer file.Close() // Ensure the file is closed after function exits

	// Read all content from the file into a byte slice
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file: ", err)
		return err
	}

	if err := c.storage.Put(ctx, c.config.StorageBucket, filePath, data, stat.Size(), true, "video/mp4"); err != nil {
		fmt.Println("Error put file to storage: ", err)
		return err
	}

	return nil
}
