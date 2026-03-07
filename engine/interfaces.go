package engine

import (
	"time"

	"github.com/ukoelguanche/graphicsengine/core"
)

type Drawer interface {
	DrawSpriteRect(sprite *core.Sprite, rect core.Rect, position core.Point)
}

type Renderable interface {
	Draw(d Drawer)

	NextFrame()
	GetMovementFrameCount() float64
	GetMovementFrame() float64

	GetSprite() *core.Sprite
	GetPosition() core.Point
	GetTargetPosition() core.Point
	GetSpeed() core.Size
	SetPosition(core.Point)
	EndMovement()
	IsMoving() bool
	SetTargetPosition(core.Point)
	SetSpeed(float64)
	GetTotalDistance() float64
	SetOnMovementComplete(func(Renderable))
	// SetOnAnimationComplete(func(Renderable))
	GetStartPosition() core.Point
	GetStartTime() time.Time
	GetDuration() time.Duration

	GetEaseFunction() func(float64) float64
	SetEaseFunction(func(float64) float64)
}

func UpdatePosition(r Renderable) {
	if !r.IsMoving() {
		return
	}

	elapsed := time.Since(r.GetStartTime())
	duration := r.GetDuration()

	t := elapsed.Seconds() / duration.Seconds()
	if t > 1.0 {
		t = 1.0
	}

	easedT := t
	if easeFunc := r.GetEaseFunction(); easeFunc != nil {
		easedT = easeFunc(t)
	}

	start := r.GetStartPosition()
	target := r.GetTargetPosition()

	nextX := start.X + (target.X-start.X)*easedT
	nextY := start.Y + (target.Y-start.Y)*easedT

	r.SetPosition(core.Point{X: nextX, Y: nextY})

	if t >= 1.0 {
		r.SetPosition(target)
		r.EndMovement()
	}
}
