package encode

import (
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
	"x11-screen-recorder/internal/record"
)

func HandleEncoding(frameChan <-chan *image.RGBA, stopChan <-chan struct{}, config record.Config) {
	cmd := exec.Command("ffmpeg",
		"-y",
		"-f", "rawvideo",
		"-pixel_format", "rgba",
		"-video_size", fmt.Sprintf("%vx%v", config.RecordResolution.X, config.RecordResolution.Y), // Change this to your screen resolution
		"-framerate", fmt.Sprintf("%v", config.RecordFPS), // Initial frame rate
		"-i", "pipe:0",
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-pix_fmt", "yuv420p",
		"output.mp4")

	ffmpegStdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to get FFmpeg stdin pipe: %v", err)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start FFmpeg: %v", err)
	}

	defer func() {
		ffmpegStdin.Close()
		cmd.Wait()
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		frameInterval := time.Second / 60 // Target frame interval for 60fps
		lastFrameTime := time.Now()

		for {
			select {
			case frame := <-frameChan:
				currentTime := time.Now()
				timeSinceLastFrame := currentTime.Sub(lastFrameTime)
				if timeSinceLastFrame < frameInterval {
					time.Sleep(frameInterval - timeSinceLastFrame)
				}
				lastFrameTime = time.Now()

				if err := writeRawFrame(ffmpegStdin, frame); err != nil {
					log.Fatalf("Failed to write frame to FFmpeg: %v", err)
				}
			case <-stopChan:
				log.Println("Received stop signal, closing FFmpeg stdin...")
				ffmpegStdin.Close()
				cmd.Wait()
				log.Println("FFmpeg process completed.")
				return
			}
		}
	}()

	wg.Wait()
}

func writeRawFrame(w io.Writer, img *image.RGBA) error {
	if _, err := w.Write(img.Pix); err != nil {
		return fmt.Errorf("failed to write raw frame data: %w", err)
	}
	return nil
}
