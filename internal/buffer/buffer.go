package buffer

import (
	"image"
	"sync"
)

type Buffer struct {
	frames []*image.RGBA
	mu     sync.Mutex
}

func NewBuffer() *Buffer {
	return &Buffer{
		frames: make([]*image.RGBA, 0),
	}
}

func (b *Buffer) AddFrame(frame *image.RGBA) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.frames = append(b.frames, frame)
}

func (b *Buffer) GetFrames() []*image.RGBA {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.frames
}

func (b *Buffer) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.frames = make([]*image.RGBA, 0)
}
