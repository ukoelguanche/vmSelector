package colorProcessors

import (
	"time"

	"github.com/ukoelguanche/graphicsengine/drivers"
)

type BlackScreen struct {
	drivers.PixelTransformer

	StartTime time.Time
	Duration  time.Duration

	OnComplete func()
	completed  bool
}

func BuildBlackScreen(duration time.Duration) *BlackScreen {
	return &BlackScreen{
		StartTime: time.Now(),
		Duration:  duration,
	}
}

func (bs *BlackScreen) Transform(pixels []byte) {
	elapsed := time.Since(bs.StartTime)

	if elapsed >= bs.Duration {
		bs.Complete()
	}

	for i := 0; i < len(pixels); i += 4 {
		pixels[i] = 0
		pixels[i+1] = 0
		pixels[i+2] = 0
	}
}

func (bs *BlackScreen) Complete() {
	bs.completed = true
	if bs.OnComplete != nil {
		bs.OnComplete()
		bs.OnComplete = nil
	}
}
