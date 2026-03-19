package colorProcessors

import (
	"time"

	"apodeiktikos.com/fbtest/engine"
	"github.com/ukoelguanche/graphicsengine/drivers"
)

type FadeToBlack struct {
	drivers.PixelTransformer
	StartTime time.Time
	Duration  time.Duration
}

func (ftb *FadeToBlack) Transform(pixels []byte) {
	elapsed := time.Since(ftb.StartTime)
	duration := ftb.Duration

	t := elapsed.Seconds() / duration.Seconds()

	if t > 1.0 {
		t = 1.0
	}

	easedT := engine.EaseInOutCubic(t)

	start := 1.0
	target := 0.0

	nextY := float32(start + (target-start)*easedT)

	for i := 0; i < len(pixels); i += 4 {
		pixels[i] = uint8(float32(pixels[i]) * nextY)
		pixels[i+1] = uint8(float32(pixels[i+1]) * nextY)
		pixels[i+2] = uint8(float32(pixels[i+2]) * nextY)
	}
}
