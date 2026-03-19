package colorProcessors

import (
	"time"

	"github.com/ukoelguanche/graphicsengine/drivers"
)

type PixelFade struct {
	drivers.PixelTransformer
	StartTime time.Time
	Duration  time.Duration
	GridSize  int
	Reverse   bool

	OnComplete func()
	completed  bool
}

func BuildPixelFade(gridSize int, duration time.Duration) *PixelFade {
	return BuildPixelFadeWithDirection(gridSize, duration, false)
}

func BuildReversePixelFade(gridSize int, duration time.Duration) *PixelFade {
	return BuildPixelFadeWithDirection(gridSize, duration, true)
}

func BuildPixelFadeWithDirection(gridSize int, duration time.Duration, reverse bool) *PixelFade {
	return &PixelFade{
		StartTime: time.Now(),
		GridSize:  gridSize,
		Duration:  duration,
		Reverse:   reverse,
	}
}

func (pf *PixelFade) Transform(pixels []byte) {
	step := pf.getTransformStep()

	for i := 0; i < drivers.VH; i++ {
		for j := 0; j < drivers.VW; j++ {
			ii := i % pf.GridSize
			jj := j % pf.GridSize

			if jj <= step-ii {
				pixelPos := (i*drivers.VW + j) * 4
				pixels[pixelPos] = 0
				pixels[pixelPos+1] = 0
				pixels[pixelPos+2] = 0
			}
		}
	}

	if pf.completed {
		pf.Complete()
	}
}

func (pf *PixelFade) getTransformStep() int {
	elapsed := time.Since(pf.StartTime)
	duration := pf.Duration

	t := elapsed.Seconds() / duration.Seconds()

	if t > 1.0 {
		t = 1.0
	}

	start := 0.0
	target := float64(pf.GridSize) * 2

	if t >= 1.0 {
		pf.completed = true
	}

	return int(float32(start + (target-start)*t))
}

func (pf *PixelFade) IsFinished() bool {
	return pf.completed
}

func (pf *PixelFade) Complete() {
	pf.completed = true
	if pf.OnComplete != nil {
		pf.OnComplete()
		pf.OnComplete = nil
	}
}
