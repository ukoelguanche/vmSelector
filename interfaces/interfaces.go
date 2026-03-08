package interfaces

import (
	"time"

	"github.com/ukoelguanche/graphicsengine/core"
)

type Drawer interface {
	DrawSpriteRect(sprite *core.Sprite, rect core.Frame, position core.Point)
}

type Animatable interface {
	GetSprite() *core.Sprite
	GetSequences(sequenceName string) []int

	// GetFrame(index int32) core.Frame
	GetCurrentFrame() core.Frame
	GetCurrentSequence() []int

	// SetOnAnimationComplete(func(Renderable))
	GetCurrentSequencePosition() float32
	SetCurrentSequencePosition(float32)
	//IncrementCurrentSequencePosition(float32)
	GetSequenceOffset() float32
	GetOnAnimationComplete() func(Animatable)
	SetOnAnimationComplete(f func(Animatable))
	ExecOnAnimationComplete()
	SetCurrentSequence([]int)
	//GetSprite() *core.Sprite
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

	GetSpeed() core.Size

	GetStartTime() time.Time
	GetDuration() time.Duration

	GetPosition() core.Point
	SetPosition(core.Point)

	GetTargetPosition() core.Point
	SetTargetPosition(core.Point)

	GetStartPosition() core.Point

	EndMovement()
	SetOnMovementComplete(func(Movable))
}

type Renderable interface {
	Drawable
	Movable

	NextFrame()
	GetSprite() *core.Sprite
}
