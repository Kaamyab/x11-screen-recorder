package terminal

import (
	"bufio"
	"fmt"
	"github.com/charmbracelet/huh"
	"os"
	"x11-screen-recorder/internal/api"
	"x11-screen-recorder/internal/encode"
	"x11-screen-recorder/internal/record"
	"x11-screen-recorder/internal/util"
)

func listenForStop(recorder *record.RecorderModel) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadByte()
		if err != nil {
			util.LogError(fmt.Errorf("Failed to read input: %w", err))
			continue
		}
		if input == 'S' || input == 's' {
			recorder.Stop()
			return
		}
	}
}

func SetupAndRun(recorder *record.RecorderModel) error {

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your method").
				Options(
					huh.NewOption("Capture Screen", "screen"),
					huh.NewOption("Capture a Window", "window"),
				).
				Value(&recorder.Config.CaptureMethod),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	switch recorder.Config.CaptureMethod {
	case "screen":
		settingsForm := createSettingsForm(recorder)
		if err := settingsForm.Run(); err != nil {
			return err
		}
	case "window":
		return fmt.Errorf("Unsupported Capture Method: %s", recorder.Config.CaptureMethod)
	default:
		return fmt.Errorf("Unsupported Capture Method: %s", recorder.Config.CaptureMethod)
	}

	recorder.WG.Add(1)
	go func() {
		defer recorder.WG.Done()
		encode.HandleEncoding(recorder.Channels.FrameChannel, recorder.Channels.Done, recorder.Config)
	}()

	go recorder.Start()
	go listenForStop(recorder)
	return nil
}

func createSettingsForm(recorder *record.RecorderModel) *huh.Form {
	microphones := api.GetMicrophoneDevices()
	microphoneOptions := make([]huh.Option[string], 0, len(microphones))

	for _, mic := range microphones {
		microphoneOptions = append(microphoneOptions, huh.NewOption(
			fmt.Sprintf("%s - %s - %s", mic.State, mic.Description, mic.VendorName),
			mic.Description,
		))
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select recording FPS").
				Options(
					huh.NewOption("60fps", 60),
					huh.NewOption("30fps", 30),
					huh.NewOption("23fps", 23),
				).
				Value(&recorder.Config.RecordFPS),
			huh.NewSelect[record.Resolution]().
				Title("Select recording resolution").
				Options(
					huh.NewOption("1920x1080", record.Resolution{X: 1920, Y: 1080}),
					huh.NewOption("1280x720", record.Resolution{X: 1280, Y: 720}),
					huh.NewOption("640x480", record.Resolution{X: 640, Y: 480}),
				).
				Value(&recorder.Config.RecordResolution),
			huh.NewSelect[string]().
				Title("Select your microphone").
				Options(microphoneOptions...),
		),
	)
}
