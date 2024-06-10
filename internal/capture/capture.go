package capture

import "image"

type ScreenCapturer interface {
	CaptureScreenShot() (*image.RGBA, error)
	CaptureWindowShot(windowID uint32) (*image.RGBA, error)
}
