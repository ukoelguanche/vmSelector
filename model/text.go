package model

import (
	"math"
	"time"

	"apodeiktikos.com/fbtest/util"
)

type Text struct {
	Sprite             *Sprite
	Text               string
	Speed              Size
	AbsSpeed           float64
	movementFrameCount float64
	movementFrame      float64
	Moving             bool
	OnMovementComplete func(t *Text)
	easeFunc           func(float64) float64
	totalDistance      float64

	// Movement
	Position       Point
	StartPosition  Point
	TargetPosition Point
	StartTime      time.Time
	Duration       time.Duration
	TotalDistance  float64
}

func (t *Text) GetSprite() *Sprite {
	return t.Sprite
}
func (t *Text) GetBitmap() *Bitmap {
	return t.Sprite.Bitmap
}
func (t *Text) NextFrame()                 { UpdatePosition(t) }
func (t *Text) SetPosition(position Point) { t.Position = position }
func (t *Text) GetStartTime() time.Time    { return t.StartTime }
func (t *Text) GetDuration() time.Duration { return t.Duration }
func (t *Text) GetStartPosition() Point    { return t.StartPosition }
func (si *Text) MoveTo(target Point, duration time.Duration) {
	si.StartPosition = si.Position
	si.TargetPosition = target
	si.StartTime = time.Now()
	si.Duration = duration
	si.Moving = true
}

func (t *Text) SetEaseFunction(f func(float64) float64) { t.easeFunc = f }
func (t *Text) GetEaseFunction() func(float64) float64  { return t.easeFunc }
func (t *Text) GetTotalDistance() float64               { return t.totalDistance }
func (t *Text) GetMovementFrameCount() float64          { return t.movementFrameCount }
func (t *Text) GetMovementFrame() float64               { return t.movementFrame }
func (t *Text) GetPosition() Point                      { return t.Position }
func (t *Text) GetTargetPosition() Point                { return t.TargetPosition }
func (t *Text) GetSpeed() Size                          { return t.Speed }
func (t *Text) IsMoving() bool                          { return t.Moving }

func (t *Text) SetTargetPosition(targetPosition Point) {
	t.TargetPosition = targetPosition
	t.Moving = true
	t.totalDistance = math.Sqrt(math.Pow(targetPosition.X-t.Position.X, 2) + math.Pow(targetPosition.Y-t.Position.Y, 2))
	t.easeFunc = util.EeaseLinear
}

func (t *Text) SetSpeed(absSpeed float64) {
	dx := t.TargetPosition.X - t.Position.X
	dy := t.TargetPosition.Y - t.Position.Y
	angle := math.Atan2(dy, dx)

	t.movementFrameCount = t.totalDistance / absSpeed
	t.movementFrame = 0

	t.AbsSpeed = absSpeed
	t.Speed = Size{W: absSpeed * math.Cos(angle), H: absSpeed * math.Sin(angle)}
}

func (t *Text) EndMovement() {
	if !t.Moving {

		return
	}
	t.Moving = false
	if t.OnMovementComplete != nil {
		t.OnMovementComplete(t)
	}
}
func (t *Text) ProcessColor(color []byte) []byte { return color }

func BuildTextInstance(sprite *Sprite, text string, position Point) *Text {
	return &Text{
		Sprite:         sprite,
		Text:           text,
		Position:       position,
		TargetPosition: position,
		Moving:         false,
	}

}
