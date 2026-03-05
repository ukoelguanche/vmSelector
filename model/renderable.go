package model

import "math"

type Renderable interface {
	GetBitmap() *Bitmap
	GetSprite() *Sprite
	ProcessColor(color []byte) []byte
	NextFrame()
	GetPosition() Point
	GetTargetPosition() Point
	GetSpeed() Size
	SetPosition(Point)
	EndMovement()
	IsMoving() bool
	SetTargetPosition(Point, Size)
}

func UpdatePosition(r Renderable) {
	if !r.IsMoving() {
		return
	}
	currentPosition := r.GetPosition()
	targetPosition := r.GetTargetPosition()
	nextPosition := currentPosition

	speed := r.GetSpeed()

	dx := targetPosition.X - currentPosition.X
	dy := targetPosition.Y - currentPosition.Y

	if dx > 0 {
		nextPosition.X += speed.W
	} else if dx < 0 {
		nextPosition.X -= speed.W
	}

	if dy > 0 {
		nextPosition.Y += speed.H
	} else if dy < 0 {
		nextPosition.Y -= speed.H
	}

	if currentPosition.Equals(nextPosition) || (math.Abs(float64(dx)) < float64(speed.W) && math.Abs(float64(dy)) < float64(speed.H)) {
		r.SetPosition(targetPosition)
		r.EndMovement()
	} else {
		r.SetPosition(nextPosition)
	}

}
