package service

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"
	"runtime"
	"strings"
	"time"
)

func RunBrowser(urlLink string) error {

	path := ""

	if runtime.GOOS == "darwin" {
		path = "/Applications/Brave Browser.app/Contents/MacOS/Brave Browser"
	}

	now := time.Now()

	url := launcher.New()

	if path != "" {
		url.Bin(path) // use Chrome instead of default Chromium
	}

	controlUrl := url.Headless(false).MustLaunch()

	browser := rod.New().ControlURL(controlUrl).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("")       // open blank first
	page.MustEmulate(devices.GalaxyS5) // emulate full mobile device

	page.MustSetWindow(0, 0, 360, 1024)

	page.MustNavigate(urlLink)

	page.MustWaitLoad()

	// Start ffmpeg

	//cmd, err := StartFFmpeg(&now)
	//if err != nil {
	//	fmt.Println("Failed to start ffmpeg:", err)
	//	return err
	//}

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

		//err = StopFFmpeg(cmd)
		//
		//e := os.Remove(fmt.Sprintf(`%s.mp4`, now.Format("2006-01-02-15-04-05")))
		//if e != nil {
		//	fmt.Println("err --> ", err)
		//}
		//
		//return err
	}

	time.Sleep(2 * time.Second) // First page wait

	target.MustClick()
	//target.MustScrollIntoView()

	//videoExists := page.MustEval(`() => !!document.querySelector("video.elementor-background-video-hosted")`).Bool()

	//if videoExists {
	//	time.Sleep(13 * time.Second) // After Click button delay
	//} else {
	time.Sleep(1 * time.Second) // After Click button delay
	//}

	//scrollInterval := 2 * time.Second // configurable scroll interval
	//scrollToBottom(page, scrollInterval, 450)
	scrollToBottomSmoothWithThrottle(page, 8, 30) // 30px max step, ~60fps (16ms delay)

	////Stop ffmpeg
	//err = StopFFmpeg(cmd)
	//if err != nil {
	//	return err
	//}

	//time.Sleep(1 * time.Second) // wait before exit
	fmt.Println("end --> ", time.Now().Sub(now).Minutes())
	return nil
}

//func scrollToBottom(page *rod.Page, interval time.Duration, step int) {
//	for {
//		// Get current scroll position and total height
//		pos := page.MustEval(`() => window.scrollY`).Int()
//		height := page.MustEval(`() => document.body.scrollHeight`).Int()
//
//		// Stop if we are at or near the bottom
//		if pos+step >= height {
//			page.MustEval(`() => window.scrollTo(0, document.body.scrollHeight)`)
//			fmt.Println("Reached the bottom.")
//			break
//		}
//
//		// Scroll by step
//		script := fmt.Sprintf("() => window.scrollTo(0, %d)", pos+step)
//		page.MustEval(script)
//
//		// Wait for new content to possibly load
//		time.Sleep(interval)
//
//		if isAtBottom(page) {
//			fmt.Println("Reached the bottom.")
//			break
//		}
//	}
//}

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
		atBottom, err := page.Eval(`() => {
            return window.innerHeight + window.scrollY >= document.body.scrollHeight - 10;
        }`)

		if err == nil && atBottom.Value.Bool() {
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

		// Timeout after 2 minutes
		if time.Since(startTime) > 2*time.Minute {
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
