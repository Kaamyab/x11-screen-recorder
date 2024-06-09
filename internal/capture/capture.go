package capture

import (
	"image"
)

// CaptureInterface defines the methods for screen capturing
type CaptureInterface interface {
	CaptureScreen() (*image.RGBA, error)
	CaptureWindow(windowID uint32) (*image.RGBA, error)
}
