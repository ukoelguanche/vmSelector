package model

type Text struct {
	Sprite             *Sprite
	Position           Point
	TargetPosition     Point
	Text               string
	Speed              Size
	moving             bool
	OnMovementComplete func(t *Text)
}

func (t *Text) GetSprite() *Sprite {
	return t.Sprite
}
func (t *Text) GetBitmap() *Bitmap {
	return t.Sprite.Bitmap
}
func (t *Text) NextFrame()                 { UpdatePosition(t) }
func (t *Text) SetPosition(position Point) { t.Position = position }
func (t *Text) SetTargetPosition(position Point, size Size) {
	t.TargetPosition = position
	t.Speed = size
	t.moving = true
}
func (t *Text) GetPosition() Point       { return t.Position }
func (t *Text) GetTargetPosition() Point { return t.TargetPosition }
func (t *Text) GetSpeed() Size           { return t.Speed }
func (t *Text) IsMoving() bool           { return t.moving }

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
