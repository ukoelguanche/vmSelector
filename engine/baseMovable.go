package engine

import (
	"math"
	"time"

	"apodeiktikos.com/fbtest/interfaces"
	"github.com/ukoelguanche/graphicsengine/core"
)

type BaseMovable struct {
	interfaces.Movable

	moving             bool
	onMovementComplete func(interfaces.Movable)
	movementFrameCount float64
	movementFrame      float64

	position       core.Point
	startPosition  core.Point
	targetPosition core.Point
	startTime      time.Time
	duration       time.Duration
	totalDistance  float64

	Acceleration float64
	MaxSpeed     float64
	Speed        core.Size
}

func (b *BaseMovable) IsMoving() bool { return b.moving }

func (b *BaseMovable) GetSpeed() core.Size { return b.Speed }

func (b *BaseMovable) GetStartTime() time.Time    { return b.startTime }
func (b *BaseMovable) GetDuration() time.Duration { return b.duration }

func (b *BaseMovable) GetPosition() core.Point         { return b.position }
func (b *BaseMovable) SetPosition(position core.Point) { b.position = position }

func (b *BaseMovable) GetStartPosition() core.Point { return b.startPosition }

func (b *BaseMovable) GetTargetPosition() core.Point { return b.targetPosition }
func (b *BaseMovable) SetTargetPosition(targetPosition core.Point) {

	b.targetPosition = targetPosition

	b.moving = true
	b.totalDistance = math.Sqrt(math.Pow(targetPosition.X-b.position.X, 2) + math.Pow(targetPosition.Y-b.position.Y, 2))
	return
}

func (b *BaseMovable) EndMovement() {
	if !b.moving {
		return
	}
	b.moving = false
	if b.onMovementComplete != nil {
		b.onMovementComplete(b)
	}
}

func (b *BaseMovable) SetOnMovementComplete(f func(movable interfaces.Movable)) {
	b.onMovementComplete = f
}

func (b *BaseMovable) UpdatePosition(r interfaces.Movable) {
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
