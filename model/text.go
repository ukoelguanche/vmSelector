package model

import (
	"math"

	"apodeiktikos.com/fbtest/util"
)

type Text struct {
	Sprite             *Sprite
	Position           Point
	TargetPosition     Point
	Text               string
	Speed              Size
	moving             bool
	OnMovementComplete func(t *Text)
	easeFunc           func(float64) float64
	totalDistance      float64
}

func (t *Text) GetSprite() *Sprite {
	return t.Sprite
}
func (t *Text) GetBitmap() *Bitmap {
	return t.Sprite.Bitmap
}
func (t *Text) NextFrame()                 { UpdatePosition(t) }
func (t *Text) SetPosition(position Point) { t.Position = position }
func (t *Text) SetTargetPosition(targetPosition Point, size Size) {
	t.TargetPosition = targetPosition
	t.Speed = size
	t.moving = true
	t.totalDistance = math.Sqrt(math.Pow(targetPosition.X-t.Position.X, 2) + math.Pow(targetPosition.Y-t.Position.Y, 2))
	t.easeFunc = util.EeaseLinear
}
func (t *Text) SetEaseFunction(f func(float64) float64) { t.easeFunc = f }
func (t *Text) GetEaseFunction() func(float64) float64  { return t.easeFunc }
func (t *Text) GetTotalDistance() float64               { return t.totalDistance }
func (t *Text) GetPosition() Point                      { return t.Position }
func (t *Text) GetTargetPosition() Point                { return t.TargetPosition }
func (t *Text) GetSpeed() Size                          { return t.Speed }
func (t *Text) IsMoving() bool                          { return t.moving }

func (t *Text) EndMovement() {
	if !t.moving {

		return
	}
	t.moving = false
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
		moving:         false,
	}

}
