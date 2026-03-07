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
