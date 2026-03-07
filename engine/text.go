package engine

import (
	"math"
	"time"

	"github.com/ukoelguanche/graphicsengine/core"
)

type Text struct {
	Sprite              *core.Sprite
	Text                string
	Speed               core.Size
	AbsSpeed            float64
	movementFrameCount  float64
	movementFrame       float64
	Moving              bool
	OnMovementComplete  func(Renderable)
	OnAnimationComplete func(Renderable)
	easeFunc            func(float64) float64
	totalDistance       float64

	// Movement
	Position       core.Point
	StartPosition  core.Point
	TargetPosition core.Point
	StartTime      time.Time
	Duration       time.Duration
	TotalDistance  float64
}

func (t *Text) GetSprite() *core.Sprite { return t.Sprite }
func (t *Text) Draw(d Drawer) {
	cursorX := t.Position.X

	letters := t.Sprite.Frames
	characters := t.Sprite.Characters

	for _, char := range t.Text {
		sChar := string(char)
		rect := letters[characters[sChar]]

		d.DrawSpriteRect(t.Sprite, rect, core.Point{X: cursorX, Y: t.Position.Y})
		cursorX += rect.Size.W + 1
	}

	t.NextFrame()
}

func (t *Text) NextFrame()                      { UpdatePosition(t) }
func (t *Text) SetPosition(position core.Point) { t.Position = position }
func (t *Text) GetStartTime() time.Time         { return t.StartTime }
func (t *Text) GetDuration() time.Duration      { return t.Duration }
func (t *Text) GetStartPosition() core.Point    { return t.StartPosition }
func (si *Text) MoveTo(target core.Point, duration time.Duration) {
	si.StartPosition = si.Position
	si.TargetPosition = target
	si.StartTime = time.Now()
	si.Duration = duration
	si.Moving = true
}
func (t *Text) SetOnMovementComplete(f func(Renderable)) { t.OnMovementComplete = f }
func (t *Text) SetEaseFunction(f func(float64) float64)  { t.easeFunc = f }
func (t *Text) GetEaseFunction() func(float64) float64   { return t.easeFunc }
func (t *Text) GetTotalDistance() float64                { return t.totalDistance }
func (t *Text) GetMovementFrameCount() float64           { return t.movementFrameCount }
func (t *Text) GetMovementFrame() float64                { return t.movementFrame }
func (t *Text) GetPosition() core.Point                  { return t.Position }
func (t *Text) GetTargetPosition() core.Point            { return t.TargetPosition }
func (t *Text) GetSpeed() core.Size                      { return t.Speed }
func (t *Text) IsMoving() bool                           { return t.Moving }

func (t *Text) SetTargetPosition(targetPosition core.Point) {
	t.TargetPosition = targetPosition
	t.Moving = true
	t.totalDistance = math.Sqrt(math.Pow(targetPosition.X-t.Position.X, 2) + math.Pow(targetPosition.Y-t.Position.Y, 2))
	t.easeFunc = EeaseLinear
}

func (t *Text) SetSpeed(absSpeed float64) {
	dx := t.TargetPosition.X - t.Position.X
	dy := t.TargetPosition.Y - t.Position.Y
	angle := math.Atan2(dy, dx)

	t.movementFrameCount = t.totalDistance / absSpeed
	t.movementFrame = 0

	t.AbsSpeed = absSpeed
	t.Speed = core.Size{W: absSpeed * math.Cos(angle), H: absSpeed * math.Sin(angle)}
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

func BuildTextInstance(sprite *core.Sprite, text string, position core.Point) *Text {
	return &Text{
		Sprite:         sprite,
		Text:           text,
		Position:       position,
		TargetPosition: position,
		Moving:         false,
	}

}
