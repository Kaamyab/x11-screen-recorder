package main

import (
	"x11-screen-recorder/internal/record"
	"x11-screen-recorder/internal/terminal"
	"x11-screen-recorder/internal/util"
)

func main() {
	util.InitLogger()
	defer util.HandlePanic()

	recorder, err := record.NewRecorder()
	if err != nil {
		util.LogError(err)
		return
	}

	if err := terminal.SetupAndRun(recorder); err != nil {
		util.LogError(err)
		return
	}

	recorder.Wait()
}
