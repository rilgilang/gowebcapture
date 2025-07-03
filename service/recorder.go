package service

import (
	"fmt"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func StartFFmpeg(currentDateTime *time.Time) (*exec.Cmd, error) {
	fmt.Println("start ffmpeg")
	dir, _ := os.Getwd()
	outputPath := filepath.Join(dir, fmt.Sprintf(`/output/%s.mp4`, currentDateTime.Format("2006-01-02-15-04-05")))

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
		"framerate":  "120",
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

func StopFFmpeg(cmd *exec.Cmd) error {
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

func deleteUnusedOutput(time *time.Time) error {
	dir, _ := os.Getwd()
	outputPath := filepath.Join(dir, fmt.Sprintf(`/output/%s.mp4`, time.Format("2006-01-02-15-04-05")))

	err := os.Remove(outputPath)
	if err != nil {
		fmt.Println("error deleting file: ", err)
	}
	return nil
}
