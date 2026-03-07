package interfaces

import (
	"time"

	"github.com/ukoelguanche/graphicsengine/core"
)

type Drawer interface {
	DrawSpriteRect(sprite *core.Sprite, rect core.Rect, position core.Point)
}

type Animatable interface {
	GetFrame(index int32) core.Rect
	GetSequences(sequenceName string) []int
	// SetOnAnimationComplete(func(Renderable))
	GetCurrentSequencePosition() float32
	SetCurrentSequencePosition(float32)
	//IncrementCurrentSequencePosition(float32)
	GetSequenceOffset() float32
	GetOnAnimationComplete() func(Animatable)
	ExecOnAnimationComplete()
}

type Drawable interface {
	Draw(d Drawer)
}

type Easable interface {
	GetEaseFunction() func(float64) float64
	SetEaseFunction(func(float64) float64)
}

type Movable interface {
	Easable
	IsMoving() bool

	GetStartTime() time.Time
	GetDuration() time.Duration

	GetPosition() core.Point
	SetPosition(core.Point)

	GetTargetPosition() core.Point
	SetTargetPosition(core.Point)

	GetStartPosition() core.Point

	EndMovement()
	SetOnMovementComplete(func(Renderable))
}

type Renderable interface {
	Drawable
	Movable

	NextFrame()
	GetSprite() *core.Sprite
	GetSpeed() core.Size
}

type BaseMovable struct {
	Movable
}

func (b *BaseMovable) UpdatePosition(r Movable) {
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

type BaseAnimatable struct {
	Animatable
}

func (b *BaseAnimatable) UpdateFrame(a Animatable) {
	currentSequencePosition := a.GetCurrentSequencePosition()
	currentSequencePosition += a.GetSequenceOffset()

	if currentSequencePosition >= 1 {
		currentSequencePosition = 0
		a.ExecOnAnimationComplete()
	}

	a.SetCurrentSequencePosition(currentSequencePosition)

}
