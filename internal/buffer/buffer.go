package buffer

import (
	"image"
	"sync"
)

type Buffer struct {
	images []*image.RGBA
	lock   sync.Mutex
}

func NewBuffer() *Buffer {
	return &Buffer{
		images: make([]*image.RGBA, 0),
	}
}

func (b *Buffer) AddImage(img *image.RGBA) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.images = append(b.images, img)
}

func (b *Buffer) GetImages() []*image.RGBA {
	b.lock.Lock()
	defer b.lock.Unlock()
	imgs := b.images
	b.images = make([]*image.RGBA, 0) // Clear buffer
	return imgs
}
