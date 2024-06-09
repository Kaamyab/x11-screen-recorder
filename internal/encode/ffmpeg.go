package encode

import (
	"fmt"
	"os/exec"
)

func EncodeImagesToVideo(imagesDir string, outputFilePath string, frameRate int) error {
	cmd := exec.Command("ffmpeg", "-framerate", fmt.Sprintf("%d", frameRate), "-i", imagesDir+"/frame_%04d.png", "-c:v", "libx264", "-pix_fmt", "yuv420p", outputFilePath)
	return cmd.Run()
}
