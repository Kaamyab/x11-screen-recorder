package record

import (
	"fmt"
	"image"
	"sync"
	"time"
	"x11-screen-recorder/internal/capture"
	"x11-screen-recorder/internal/util"
)

type RecorderModel struct {
	Capturer capture.ScreenCapturer
	Channels Channels
	Config   Config
	Ticker   *time.Ticker
	WG       sync.WaitGroup
	Statics  RecorderStatics
}

type RecorderStatics struct {
	averageFrameTime float64
	bufferLength     int
}

type Channels struct {
	FrameChannel chan *image.RGBA
	Done         chan struct{}
}

func NewRecorder() (*RecorderModel, error) {
	screenCapturer, err := capture.NewScreenCapturer()
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize X11 capture: %w", err)
	}

	frameChan := make(chan *image.RGBA, 100)

	recorder := &RecorderModel{
		Capturer: screenCapturer,
		Channels: Channels{
			FrameChannel: frameChan,
			Done:         make(chan struct{}),
		},
		Ticker: time.NewTicker(1 * time.Second),
	}

	return recorder, nil
}

func (recorder *RecorderModel) Start() {
	frameTimes := []time.Duration{}
	recorder.WG.Add(1)
	go recorder.monitorFrameChannelLength()

	for {
		select {
		case <-recorder.Channels.Done:
			return
		default:
			start := time.Now()
			frame, err := recorder.Capturer.CaptureScreenShot()
			if err != nil {
				util.LogError(fmt.Errorf("Failed to capture screen: %w", err))
				continue
			}
			recorder.Channels.FrameChannel <- frame
			frameTime := time.Since(start)
			frameTimes = append(frameTimes, frameTime)
			if len(frameTimes) > recorder.Config.RecordFPS {
				frameTimes = frameTimes[1:]
			}
			recorder.Statics.averageFrameTime = float64(avgDuration(frameTimes).Milliseconds())
		}
	}
}

func (recorder *RecorderModel) Stop() {
	close(recorder.Channels.Done)
	recorder.Ticker.Stop()
	recorder.Wait()
}

func (recorder *RecorderModel) Wait() {
	recorder.WG.Wait()
}

func avgDuration(durations []time.Duration) time.Duration {
	var sum time.Duration
	for _, d := range durations {
		sum += d
	}
	return sum / time.Duration(len(durations))
}

func (recorder *RecorderModel) monitorFrameChannelLength() {
	defer recorder.WG.Done()
	for {
		select {
		case <-recorder.Channels.Done:
			return
		case <-recorder.Ticker.C:
			recorder.Statics.bufferLength = len(recorder.Channels.FrameChannel)
			// Log the current length of frameChan for debugging purposes
			util.LogError(fmt.Errorf("Current frameChan length: %d", recorder.Statics.bufferLength))
		}
	}
}
