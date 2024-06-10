//go:build linux

package capture

import (
	"fmt"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
	"image"
	"image/color"
)

type X11 struct {
	conn *xgb.Conn
}

func NewScreenCapturer() (ScreenCapturer, error) {
	return NewX11()
}

func NewX11() (*X11, error) {
	conn, err := xgb.NewConn()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to X server: %w", err)
	}
	return &X11{conn: conn}, nil
}

func (x *X11) CaptureScreenShot() (*image.RGBA, error) {
	setup := xproto.Setup(x.conn)
	screen := setup.DefaultScreen(x.conn)
	return captureWindow(x.conn, xproto.Window(screen.Root))
}

func (x *X11) CaptureWindowShot(windowID uint32) (*image.RGBA, error) {
	return captureWindow(x.conn, xproto.Window(windowID))
}

func captureWindow(conn *xgb.Conn, window xproto.Window) (*image.RGBA, error) {
	geo, err := xproto.GetGeometry(conn, xproto.Drawable(window)).Reply()
	if err != nil {
		return nil, fmt.Errorf("failed to get window geometry: %w", err)
	}

	xImg, err := xproto.GetImage(conn, xproto.ImageFormatZPixmap, xproto.Drawable(window), 0, 0, geo.Width, geo.Height, 0xffffffff).Reply()
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, int(geo.Width), int(geo.Height)))
	for y := 0; y < int(geo.Height); y++ {
		for x := 0; x < int(geo.Width); x++ {
			offset := (y*int(geo.Width) + x) * 4
			r := xImg.Data[offset+2]
			g := xImg.Data[offset+1]
			b := xImg.Data[offset]
			a := uint8(255)
			img.SetRGBA(x, y, color.RGBA{r, g, b, a})
		}
	}

	return img, nil
}
