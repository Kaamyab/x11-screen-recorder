//go:build linux

package api

import (
	"fmt"
	"os/exec"
	"strings"
	"x11-screen-recorder/internal/util"
)

// Placeholder function for getting resolutions
func getResolutions() []string {
	// Replace this with actual code to get the supported resolutions
	return []string{"1920x1080", "1280x720", "1024x768"}
}

// Placeholder function for getting audio devices
func getAudioDevices() []string {
	out, err := exec.Command("pactl", "list", "short", "sinks").Output()
	if err != nil {
		util.LogError(fmt.Errorf("Failed to list audio devices: %w", err))
		return []string{}
	}
	var devices []string
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if fields := strings.Fields(line); len(fields) > 1 {
			devices = append(devices, fields[1])
		}
	}
	return devices
}

type Microphones struct {
	State       string
	Name        string
	Description string
	VendorName  string
}

// Placeholder function for getting microphone devices with user-friendly names
func GetMicrophoneDevices() []Microphones {
	out, err := exec.Command("pactl", "list", "sources").Output()
	if err != nil {
		util.LogError(fmt.Errorf("Failed to list microphones devices: %w", err))
		return []Microphones{}
	}

	var devices []Microphones
	var currentState, currentDescription, currentVendorName string
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "State:") {
			currentState = strings.TrimPrefix(line, "State: ")
		} else if strings.HasPrefix(line, "Description:") {
			currentDescription = strings.TrimPrefix(line, "Description: ")
		} else if strings.HasPrefix(line, "device.vendor.name =") {
			currentVendorName = strings.Trim(line[len("device.vendor.name ="):], `"`)
		}

		// When all required fields are collected, format and store the result
		if currentState != "" && currentDescription != "" && currentVendorName != "" {

			devices = append(devices, Microphones{
				State:       currentState,
				Name:        "",
				Description: currentDescription,
				VendorName:  currentVendorName,
			})
			currentState, currentDescription, currentVendorName = "", "", "" // Reset for the next source
		}
	}
	return devices
}
